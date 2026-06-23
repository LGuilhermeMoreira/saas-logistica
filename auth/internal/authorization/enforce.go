package authorization

import (
	"fmt"

	"auth/config"

	"github.com/casbin/casbin/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxadapter "github.com/pckhoi/casbin-pgx-adapter/v3"
)

type Enforcer interface {
	Enforce(subject, route, method string) (bool, error)
	AddPermissionsToRole(roleID string, permissions []Permission) error
	RemovePermissionsToRole(roleID string, permissions []Permission) error
	RemoveRole(roleID string) error
	AssignRoleToUser(userID, roleID string) error
}

type CasbinEnforcer struct {
	cas *casbin.SyncedEnforcer
}

func NewEnforcer(env *config.Env, pool *pgxpool.Pool) (Enforcer, error) {
	adapter, err := pgxadapter.NewAdapter(
		env.PostgresURI(),
		pgxadapter.WithConnectionPool(pool),
		pgxadapter.WithDatabase(env.DATABASE_NAME),
	)
	if err != nil {
		return nil, fmt.Errorf("casbin adapter: %w", err)
	}

	enforcer, err := casbin.NewSyncedEnforcer("model.conf", adapter)
	if err != nil {
		return nil, fmt.Errorf("casbin enforcer: %w", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("load policy: %w", err)
	}

	enforcer.EnableAutoSave(true)

	return &CasbinEnforcer{cas: enforcer}, nil
}

func (e *CasbinEnforcer) AddPermissionsToRole(roleID string, permissions []Permission) error {
	var rules [][]string
	for _, perm := range permissions {
		rules = append(rules, []string{
			"role:" + roleID,
			perm.Route,
			perm.Action,
		})
	}
	_, err := e.cas.AddPolicies(rules)
	return err
}

func (e *CasbinEnforcer) RemovePermissionsToRole(roleID string, permissions []Permission) error {
	var rules [][]string
	for _, perm := range permissions {
		rules = append(rules, []string{
			"role:" + roleID,
			perm.Route,
			perm.Action,
		})
	}
	_, err := e.cas.RemovePolicies(rules)
	return err
}

func (e *CasbinEnforcer) RemoveRole(roleID string) error {
	_, err := e.cas.DeleteRole("role:" + roleID)
	return err
}

func (e *CasbinEnforcer) AssignRoleToUser(userID, roleID string) error {
	_, err := e.cas.AddRoleForUser("user:"+userID, "role:"+roleID)
	return err
}

func (e *CasbinEnforcer) Enforce(subject, route, method string) (bool, error) {
	return e.cas.Enforce(subject, route, method)
}
