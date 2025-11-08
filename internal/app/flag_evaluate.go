package app

import (
	"context"
	"crypto/sha1"
)

// Инпут — кто и какой флаг проверяет
type EvaluateInput struct {
	Tenant TenantID
	Key    FlagKey
	UserID string // простой актор; потом расширим сегментами/атрибутами
}

type EvaluateOutput struct {
	Enabled bool
}

type FlagEvaluateUsecase struct {
	Repo FlagRepository
}

func NewFlagEvaluateUsecase(repo FlagRepository) *FlagEvaluateUsecase {
	return &FlagEvaluateUsecase{Repo: repo}
}

func (u *FlagEvaluateUsecase) Do(ctx context.Context, in EvaluateInput) (EvaluateOutput, error) {
	f, err := u.Repo.Get(ctx, in.Tenant, in.Key)
	if err != nil {
		return EvaluateOutput{}, err
	}
	// простейший детерминированный bucketing: hash(userID+key) % 100 < rollout
	h := sha1.Sum([]byte(in.UserID + string(in.Key)))
	bucket := int(h[0]) % 100
	return EvaluateOutput{Enabled: bucket < int(f.Rule.Rollout)}, nil
}
