package authorization_test

import (
	"auth/config"
	"auth/pkg/authorization"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOPA_Validate(t *testing.T) {
	input := authorization.OPAInput{
		Action: "GET",
		Path:   "/api/users",
	}
	input.User.Role = "admin"

	t.Run("success allowed", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"result": {"allow": true}}`))
		}))
		defer server.Close()

		env := &config.Env{OPA_ADDRESS: server.URL}
		opa := authorization.NewOPA(env)

		err := opa.Validate(input)

		assert.NoError(t, err)
	})

	t.Run("error forbidden", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"result": {"allow": false}}`))
		}))
		defer server.Close()

		env := &config.Env{OPA_ADDRESS: server.URL}
		opa := authorization.NewOPA(env)

		err := opa.Validate(input)

		assert.ErrorIs(t, err, authorization.ErrForbidden)
	})

	t.Run("error unexpected status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		env := &config.Env{OPA_ADDRESS: server.URL}
		opa := authorization.NewOPA(env)

		err := opa.Validate(input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "opa returned unexpected status: 500")
	})

	t.Run("error malformed json response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{ invalid_json:`))
		}))
		defer server.Close()

		env := &config.Env{OPA_ADDRESS: server.URL}
		opa := authorization.NewOPA(env)

		err := opa.Validate(input)

		assert.Error(t, err)
	})

	t.Run("error connection failure", func(t *testing.T) {
		env := &config.Env{OPA_ADDRESS: "http://invalid-url-that-does-not-exist"}
		opa := authorization.NewOPA(env)

		err := opa.Validate(input)

		assert.Error(t, err)
	})
}
