package app

import (
	"errors"
	"fmt"
)

type Kind string

const (
	KindInvalid   Kind = "invalid"
	KindNotFound  Kind = "not_found"
	KindForbidden Kind = "forbidden"
	KindConflict  Kind = "conflict"
	KindInternal  Kind = "internal"
)

type E struct {
	Kind Kind   // семантика (класс) ошибки
	Op   string // операция/контекст, где произошла
	Err  error  // первопричина (wrapped)
	Msg  string // безопасное сообщение для клиента/лога
}

func (e *E) Error() string {
	// Короткая читаемая строка
	if e.Msg != "" {
		if e.Op != "" {
			return fmt.Sprintf("%s: %s: %s", e.Op, e.Kind, e.Msg)
		}
		return fmt.Sprintf("%s: %s", e.Kind, e.Msg)
	}
	if e.Op != "" {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Kind, e.Err)
	}
	return fmt.Sprintf("%s: %v", e.Kind, e.Err)
}

func (e *E) Unwrap() error { return e.Err }

// Фабрики
func Invalid(op, msg string) error   { return &E{Kind: KindInvalid, Op: op, Msg: msg} }
func NotFound(op, msg string) error  { return &E{Kind: KindNotFound, Op: op, Msg: msg} }
func Forbidden(op, msg string) error { return &E{Kind: KindForbidden, Op: op, Msg: msg} }
func Conflict(op, msg string) error  { return &E{Kind: KindConflict, Op: op, Msg: msg} }
func Internal(op string, err error) error {
	if err == nil {
		return &E{Kind: KindInternal, Op: op}
	}
	return &E{Kind: KindInternal, Op: op, Err: err}
}

// Wrap — оборачиваем чужую ошибку текущей операцией, сохраняя Kind если это уже app.E
func Wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	var ae *E
	if errors.As(err, &ae) {
		return &E{Kind: ae.Kind, Op: op, Err: err}
	}
	return &E{Kind: KindInternal, Op: op, Err: err}
}
