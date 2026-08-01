package authentication

import (
	"auth/config"
	"errors"
	"fmt"
	"time"

	jwt_go "github.com/golang-jwt/jwt/v5"
)

type JWTToken string

type JWTInterface interface {
	TokenGenerator
	TokenValidator
}

type TokenGenerator interface {
	GenerateToken(data any) (JWTToken, error)
	GenerateShortToken(data any) (JWTToken, error)
}

type TokenValidator interface {
	Validate(tokenStr string) error
	ExtractClaims(tokenStr string) (jwt_go.MapClaims, error)
}
type JWT struct {
	timeToLive time.Duration
	secret     string
}

func NewJWT(env *config.Env) *JWT {
	return &JWT{
		timeToLive: env.JWT_TTL,
		secret:     env.JWT_SECRET,
	}
}

func (j *JWT) Generate(data any) (string, error) {

	claims := jwt_go.MapClaims{
		"exp":  time.Now().Add(j.timeToLive).Unix(),
		"iat":  time.Now().Unix(),
		"data": data,
	}

	token := jwt_go.NewWithClaims(jwt_go.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("jwt: falha ao assinar token: %w", err)
	}

	return signed, nil
}

func (j *JWT) GenerateShortToken(data any) (string, error) {
	claims := jwt_go.MapClaims{
		"exp":  time.Now().Add(time.Minute * 30).Unix(),
		"iat":  time.Now().Unix(),
		"data": data,
	}

	token := jwt_go.NewWithClaims(jwt_go.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("jwt: falha ao assinar token: %w", err)
	}

	return signed, nil

}

func (j *JWT) Validate(tokenStr string) error {
	token, err := jwt_go.Parse(tokenStr, func(t *jwt_go.Token) (any, error) {
		if _, ok := t.Method.(*jwt_go.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: método de assinatura inválido: %v", t.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt_go.ErrTokenExpired) {
			return fmt.Errorf("jwt: token expirado")
		}
		return fmt.Errorf("jwt: token inválido: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("jwt: token é inválido")
	}

	return nil
}

func (j *JWT) ExtractClaims(tokenStr string) (jwt_go.MapClaims, error) {
	token, err := jwt_go.Parse(tokenStr, func(t *jwt_go.Token) (any, error) {
		if _, ok := t.Method.(*jwt_go.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt: %w", err)
	}

	claims, ok := token.Claims.(jwt_go.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("jwt: could not extract claims")
	}

	return claims, nil
}
