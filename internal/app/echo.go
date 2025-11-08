package app

import (
	"context"
)

type EchoUsecase struct{}

func NewEchoUsecase() *EchoUsecase {
	return &EchoUsecase{}
}

func (u *EchoUsecase) Do(ctx context.Context, msg string) (string, error) {
	if len(msg) == 0 {
		return "", WithOp("EchoUsecase.Do", Invalid("msg is empty", nil))
	}
	// будущая бизнес-логика (например, запись в БД)
	// сейчас просто возвращаем echo
	return msg, nil
}
