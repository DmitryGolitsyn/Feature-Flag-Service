package app

import "context"

type FlagRepository interface {
	Upsert(ctx context.Context, f Flag) error
	Get(ctx context.Context, tenant TenantID, flagKey FlagKey) (Flag, error)
}
