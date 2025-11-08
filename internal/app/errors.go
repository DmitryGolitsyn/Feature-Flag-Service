package app

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid_input")
	// ErrForbidden = errors.New("forbidden")
	// ErrInternal  = errors.New("internal_error")
)

//// Invalid оборачивает как "invalid_input"
//func Invalid(msg string, cause error) error {
//	if cause != nil {
//		return fmt.Errorf("%w: %s: %v", ErrInvalid, msg, cause)
//	}
//	return fmt.Errorf("%w: %s", ErrInvalid, msg)
//}
//
//// NotFound оборачивает как "not_found"
//func NotFound(entity string, cause error) error {
//	if cause != nil {
//		return fmt.Errorf("%w: %s: %v", ErrNotFound, entity, cause)
//	}
//	return fmt.Errorf("%w: %s", ErrNotFound, entity)
//}
//
//// Forbidden оборачивает как "forbidden"
//func Forbidden(reason string, cause error) error {
//	if cause != nil {
//		return fmt.Errorf("%w: %s: %v", ErrForbidden, reason, cause)
//	}
//	return fmt.Errorf("%w: %s", ErrForbidden, reason)
//}
//
//// Internal оборачивает как "internal_error"
//func Internal(msg string, cause error) error {
//	if cause != nil {
//		return fmt.Errorf("%w: %s: %v", ErrInternal, msg, cause)
//	}
//	return fmt.Errorf("%w: %s", ErrInternal, msg)
//}

//// WithOp добавляет "операцию" (контекст) поверх любой ошибки
//// пример: return WithOp("Echo.Usecase.Do", Invalid("empty msg", nil))
//func WithOp(op string, err error) error {
//	if err == nil {
//		return nil
//	}
//	return fmt.Errorf("%s: %w", op, err)
//}
