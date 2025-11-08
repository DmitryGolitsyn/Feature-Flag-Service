package mem

import (
	"context"
	"sync"

	"ffs-tutorial/internal/app"
)

type FlagRepo struct {
	mu   sync.RWMutex
	data map[string]app.Flag // key = tenant+"|"+flagKey
}

func NewFlagRepo() *FlagRepo {
	return &FlagRepo{
		data: make(map[string]app.Flag),
	}
}

func key(t app.TenantID, k app.FlagKey) string {
	return string(t) + "|" + string(k)
}

func (r *FlagRepo) Upsert(ctx context.Context, f app.Flag) error {
	if err := f.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key(f.Tenant, f.Key)] = f
	return nil
}

func (r *FlagRepo) Get(ctx context.Context, tenant app.TenantID, flagKey app.FlagKey) (app.Flag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if f, ok := r.data[key(tenant, flagKey)]; ok {
		return f, nil
	}
	return app.Flag{}, app.ErrNotFound
}
