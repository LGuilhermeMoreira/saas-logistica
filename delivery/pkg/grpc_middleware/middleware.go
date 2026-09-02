package grpcmiddleware

import (
	"context"
	"delivery/pkg/authentication"
	"delivery/pkg/authorization"
	"delivery/pkg/logger"
	"log/slog"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCMiddleware struct {
	jwt authentication.TokenValidator
	opa authorization.OPAInterface
}

func NewGRPCMiddleware(jwt authentication.TokenValidator, opa authorization.OPAInterface) *GRPCMiddleware {
	return &GRPCMiddleware{
		jwt: jwt,
		opa: opa,
	}
}

func (g *GRPCMiddleware) ValidateCredentials() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/auth.v1.AuthService/Login" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			slog.Error("missing metadata", slog.String("grpc_method", info.FullMethod))
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		var reqID string
		if ids := md.Get("x-request-id"); len(ids) > 0 {
			reqID = ids[0]
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 || authHeader[0] == "" {
			slog.Warn("missing authorization token", slog.String("request_id", reqID))
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")

		claims, err := g.jwt.ExtractClaims(token)
		if err != nil {
			slog.Warn("invalid authorization token", slog.String("request_id", reqID), slog.Any("error", err))
			return nil, status.Error(codes.Unauthenticated, "invalid authorization token")
		}

		data, ok := claims["data"].(map[string]any)
		if !ok {
			slog.Warn("invalid token claims structure", slog.String("request_id", reqID))
			return nil, status.Error(codes.Unauthenticated, "invalid token claims structure")
		}

		role, ok := data["role_name"].(string)
		if !ok {
			slog.Warn("missing role in token", slog.String("request_id", reqID))
			return nil, status.Error(codes.Unauthenticated, "missing role in token")
		}

		httpPath := md.Get("http-path")
		httpMethod := md.Get("http-method")
		if len(httpPath) == 0 || len(httpMethod) == 0 {
			slog.Error("missing http path or method for validation", slog.String("request_id", reqID))
			return nil, status.Error(codes.InvalidArgument, "missing http path or method for validation")
		}

		opaInput := authorization.OPAInput{
			Action: httpMethod[0],
			Path:   httpPath[0],
		}
		opaInput.User.Role = role

		if err := g.opa.Validate(opaInput); err != nil {
			slog.Warn("access denied by OPA",
				slog.String("request_id", reqID),
				slog.String("role", role),
				slog.String("http_path", httpPath[0]),
				slog.String("http_method", httpMethod[0]),
			)
			return nil, status.Error(codes.PermissionDenied, "access denied by OPA")
		}
		ctx = context.WithValue(ctx, logger.RequestIDKey, reqID)

		return handler(ctx, req)
	}
}
