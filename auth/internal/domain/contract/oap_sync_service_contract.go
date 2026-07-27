package contract

import "context"

type OPASyncServiceInterface interface {
	SyncPolicies(ctx context.Context) error
}
