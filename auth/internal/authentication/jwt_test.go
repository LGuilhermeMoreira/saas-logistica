package authentication_test

import (
	"testing"
	"time"

	"auth/config"
	"auth/internal/authentication"

	jwt_go "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// helper pra criar JWT sem depender de .env
func newTestJWT(ttl time.Duration, secret string) *authentication.JWT {
	env := &config.Env{}
	env.JWT_TTL = ttl
	env.JWT_SECRET = secret

	return authentication.NewJWT(env)
}

func TestGenerate_Success(t *testing.T) {
	j := newTestJWT(time.Minute, "secret")

	data := map[string]any{"id": "123"}

	token, err := j.Generate(data)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidate_Success(t *testing.T) {
	j := newTestJWT(time.Minute, "secret")

	token, err := j.Generate(map[string]any{"id": "123"})
	assert.NoError(t, err)

	err = j.Validate(token)

	assert.NoError(t, err)
}

func TestValidate_InvalidToken(t *testing.T) {
	j := newTestJWT(time.Minute, "secret")

	err := j.Validate("token_invalido")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token inválido")
}

func TestValidate_ExpiredToken(t *testing.T) {
	j := newTestJWT(-time.Minute, "secret") // já nasce expirado

	token, err := j.Generate(map[string]any{"id": "123"})
	assert.NoError(t, err)

	err = j.Validate(token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expirado")
}

func TestValidate_InvalidSignature(t *testing.T) {
	j1 := newTestJWT(time.Minute, "secret1")
	j2 := newTestJWT(time.Minute, "secret2")

	token, err := j1.Generate(map[string]any{"id": "123"})
	assert.NoError(t, err)

	err = j2.Validate(token)

	assert.Error(t, err)
}

func TestExtractClaims_Success(t *testing.T) {
	j := newTestJWT(time.Minute, "secret")

	data := map[string]any{"id": "123"}

	token, err := j.Generate(data)
	assert.NoError(t, err)

	claims, err := j.ExtractClaims(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)

	assert.Equal(t, "123", claims["data"].(map[string]any)["id"])
}

func TestExtractClaims_InvalidToken(t *testing.T) {
	j := newTestJWT(time.Minute, "secret")

	claims, err := j.ExtractClaims("token_invalido")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestExtractClaims_InvalidSignature(t *testing.T) {
	j1 := newTestJWT(time.Minute, "secret1")
	j2 := newTestJWT(time.Minute, "secret2")

	token, err := j1.Generate(map[string]any{"id": "123"})
	assert.NoError(t, err)

	claims, err := j2.ExtractClaims(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGenerate_ContainsClaims(t *testing.T) {
	j := newTestJWT(time.Minute, "secret")

	data := map[string]any{"id": "123"}

	tokenStr, err := j.Generate(data)
	assert.NoError(t, err)

	token, err := jwt_go.Parse(tokenStr, func(t *jwt_go.Token) (any, error) {
		return []byte("secret"), nil
	})
	assert.NoError(t, err)

	claims := token.Claims.(jwt_go.MapClaims)

	assert.Contains(t, claims, "exp")
	assert.Contains(t, claims, "iat")
	assert.Contains(t, claims, "data")
}
