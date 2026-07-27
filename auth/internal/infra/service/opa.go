package service

import (
	"auth/config"
	"auth/internal/domain/contract"
	"context"
)

type OPASyncService struct {
	repo contract.RoleRepositoryInterface
}

func NewOPASyncService(env *config.Env) contract.OPASyncServiceInterface {
	return &OPASyncService{}
}

func (o *OPASyncService) SyncPolicies(ctx context.Context) error
