package app

import (
	"context"
)

// PingUsecase — пример бизнес-логики.
// Пока он просто возвращает true, но позже здесь появится:
//   - валидация
//   - БД
//   - Kafka-события
//   - трейсинг
//   - error wrapping
type PingUsecase struct{}

func NewPingUsecase() *PingUsecase {
	return &PingUsecase{}
}

func (u *PingUsecase) Ping(ctx context.Context) (bool, error) {
	// сейчас usecase ничего не делает
	return true, nil
}
