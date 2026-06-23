package authorization

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnforce_AllowsWildcardActionForModuleRoutes(t *testing.T) {
	enforcer, err := casbin.NewSyncedEnforcer("../../model.conf")
	require.NoError(t, err)

	_, err = enforcer.AddPolicy("role:admin", "/company/*", "*")
	require.NoError(t, err)
	_, err = enforcer.AddRoleForUser("user:123", "role:admin")
	require.NoError(t, err)

	casbinEnforcer := &CasbinEnforcer{cas: enforcer}

	tests := []struct {
		name   string
		route  string
		method string
		want   bool
	}{
		{name: "module get", route: "/company/123", method: "GET", want: true},
		{name: "module patch", route: "/company/123", method: "PATCH", want: true},
		{name: "module delete", route: "/company/123/logo", method: "DELETE", want: true},
		{name: "different module denied", route: "/role/123", method: "GET", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := casbinEnforcer.Enforce("user:123", tt.route, tt.method)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
