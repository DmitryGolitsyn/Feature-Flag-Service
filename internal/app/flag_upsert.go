package app

import "context"

type FlagUpsertUsecase struct {
	Repo FlagRepository
}

func NewFlagUpsertUsecase(repo FlagRepository) *FlagUpsertUsecase {
	return &FlagUpsertUsecase{Repo: repo}
}

func (u *FlagUpsertUsecase) Do(ctx context.Context, f Flag) error {
	return u.Repo.Upsert(ctx, f)
}
