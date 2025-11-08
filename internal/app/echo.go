package app

import (
	"context"
	"strings"
)

type EchoUsecase struct{}

func NewEchoUsecase() *EchoUsecase {
	return &EchoUsecase{}
}

func (u *EchoUsecase) Do(ctx context.Context, msg string) (string, error) {
	const op = "app.EchoUsecase.Do"

	// Простейшая бизнес-валидация
	if strings.TrimSpace(msg) == "" {
		return "", Invalid(op, "msg is empty")
	}
	// Если тут будет вызов БД/внешнего API:
	// if err := repo.Save(...); err != nil {
	//     return "", Wrap(op, err) // сохраним Kind, либо превратим в internal
	// }

	return msg, nil
}
