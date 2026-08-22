package auth

import (
	"log/slog"
	"net/http"

	"gateway/internal/domain/contract"
	authv1 "proto/gen/auth/v1"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService contract.AuthService
	log         *slog.Logger
}

func NewAuthController(authService contract.AuthService, log *slog.Logger) *AuthController {
	return &AuthController{
		authService: authService,
		log:         log,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	reqID := ctx.Writer.Header().Get("X-Request-ID")
	c.log.Info("handling login request", slog.String("request_id", reqID))

	var input LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.log.Warn("invalid login payload", slog.String("request_id", reqID), slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authService.Login(ctx.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.log.Error("auth service login failed", slog.String("request_id", reqID), slog.String("email", input.Email), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.log.Info("login successful", slog.String("request_id", reqID), slog.String("email", input.Email))
	ctx.JSON(http.StatusOK, LoginOutput{Token: token})
}

func (c *AuthController) CreateUser(ctx *gin.Context) {
	reqID := ctx.Writer.Header().Get("X-Request-ID")
	c.log.Info("handling create user request", slog.String("request_id", reqID))

	var input CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.log.Warn("invalid create user payload", slog.String("request_id", reqID), slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.authService.CreateUser(ctx.Request.Context(), input.Name, input.Email, input.Password, input.RoleID)
	if err != nil {
		c.log.Error("auth service create user failed", slog.String("request_id", reqID), slog.String("email", input.Email), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.log.Info("user created successfully", slog.String("request_id", reqID), slog.String("user_id", resp.Id))
	ctx.JSON(http.StatusCreated, UserOutput{
		ID:    resp.Id,
		Name:  resp.Name,
		Email: resp.Email,
	})
}

func (c *AuthController) DeleteUser(ctx *gin.Context) {
	reqID := ctx.Writer.Header().Get("X-Request-ID")
	id := ctx.Param("id")
	c.log.Info("handling delete user request", slog.String("request_id", reqID), slog.String("user_id", id))

	resp, err := c.authService.DeleteUser(ctx.Request.Context(), id)
	if err != nil {
		c.log.Error("auth service delete user failed", slog.String("request_id", reqID), slog.String("user_id", id), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.log.Info("user deleted successfully", slog.String("request_id", reqID), slog.String("user_id", id))
	ctx.JSON(http.StatusOK, DeleteUserOutput{Msg: resp.Msg})
}

func (c *AuthController) CreateRole(ctx *gin.Context) {
	reqID := ctx.Writer.Header().Get("X-Request-ID")
	c.log.Info("handling create role request", slog.String("request_id", reqID))

	var input CreateRoleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.log.Warn("invalid create role payload", slog.String("request_id", reqID), slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var permissions []*authv1.Permission
	for _, p := range input.Permissions {
		permissions = append(permissions, &authv1.Permission{
			Action: p.Action,
			Path:   p.Path,
		})
	}

	resp, err := c.authService.CreateRole(ctx.Request.Context(), input.Name, input.Description, permissions)
	if err != nil {
		c.log.Error("auth service create role failed", slog.String("request_id", reqID), slog.String("role_name", input.Name), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var createdPermissions []Permission
	for _, p := range resp.Permissions {
		createdPermissions = append(createdPermissions, Permission{
			Action: p.Action,
			Path:   p.Path,
		})
	}

	c.log.Info("role created successfully", slog.String("request_id", reqID), slog.String("role_id", resp.Id))
	ctx.JSON(http.StatusCreated, RoleOutput{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		Permissions: createdPermissions,
	})
}

func (c *AuthController) DeleteRole(ctx *gin.Context) {
	reqID := ctx.Writer.Header().Get("X-Request-ID")
	id := ctx.Param("id")
	c.log.Info("handling delete role request", slog.String("request_id", reqID), slog.String("role_id", id))

	resp, err := c.authService.DeleteRole(ctx.Request.Context(), id)
	if err != nil {
		c.log.Error("auth service delete role failed", slog.String("request_id", reqID), slog.String("role_id", id), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.log.Info("role deleted successfully", slog.String("request_id", reqID), slog.String("role_id", id))
	ctx.JSON(http.StatusOK, DeleteRoleOutput{Msg: resp.Msg})
}

func (c *AuthController) FindUserByID(ctx *gin.Context) {
	reqID := ctx.Writer.Header().Get("X-Request-ID")
	id := ctx.Param("id")
	c.log.Info("handling find user by id request", slog.String("request_id", reqID), slog.String("user_id", id))

	resp, err := c.authService.FindUserByID(ctx.Request.Context(), id)
	if err != nil {
		c.log.Error("auth service find user by id failed", slog.String("request_id", reqID), slog.String("user_id", id), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, UserOutput{
		ID:     resp.Id,
		Email:  resp.Email,
		Name:   resp.Name,
		RoleID: resp.RoleId,
	})
}

func (c *AuthController) FindRoleByID(ctx *gin.Context) {
	reqID := ctx.Writer.Header().Get("X-Request-ID")
	id := ctx.Param("id")
	c.log.Info("handling find role by id request", slog.String("request_id", reqID), slog.String("role_id", id))

	resp, err := c.authService.FindRoleByID(ctx.Request.Context(), id)
	if err != nil {
		c.log.Error("auth service find role by id failed", slog.String("request_id", reqID), slog.String("role_id", id), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	perm := make([]Permission, len(resp.Permissions))
	for i, v := range resp.Permissions {
		perm[i] = Permission{
			Action: v.Action,
			Path:   v.Path,
		}
	}

	ctx.JSON(http.StatusOK, RoleOutput{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		Permissions: perm,
	})
}

func (c *AuthController) FindAllRoles(ctx *gin.Context) {
	reqID := ctx.Writer.Header().Get("X-Request-ID")
	c.log.Info("handling find all roles request", slog.String("request_id", reqID))

	resp, err := c.authService.FindAllRoles(ctx.Request.Context())
	if err != nil {
		c.log.Error("auth service find all roles failed", slog.String("request_id", reqID), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roles := make([]RoleOutput, len(resp.Roles))
	for i, v := range resp.Roles {
		roles[i] = RoleOutput{
			ID:          v.Id,
			Description: v.Description,
			Name:        v.Name,
		}
	}

	ctx.JSON(http.StatusOK, FindAllRolesOutput{
		Roles: roles,
	})
}
