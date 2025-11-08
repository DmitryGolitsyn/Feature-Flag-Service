package httpserver

import (
	"errors"
	"ffs-tutorial/internal/app"
	"net/http"
)

// mapAppErrorToHTTP переводит ошибку usecase в HTTP-статус.
// Важно: использовать errors.Is, т.к. usecase возвращают обёрнутые ошибки.
func httpStatusFor(err error) int {
	// Новое: сначала пробуем типизированную ошибку
	var ae *app.E
	if errors.As(err, &ae) {
		switch ae.Kind {
		case app.KindInvalid:
			return http.StatusBadRequest // 400
		case app.KindForbidden:
			return http.StatusForbidden // 403
		case app.KindNotFound:
			return http.StatusNotFound // 404
		case app.KindConflict:
			return http.StatusConflict // 409
		default:
			return http.StatusInternalServerError // 500
		}
	}

	// Старые sentinel'ы — на переходный период
	switch {

	case errors.Is(err, app.ErrInvalid):
		return http.StatusBadRequest // 400

	case errors.Is(err, app.ErrForbidden):
		return http.StatusForbidden // 403

	case errors.Is(err, app.ErrNotFound):
		return http.StatusNotFound // 404

	case errors.Is(err, app.ErrConflict):
		return http.StatusConflict // 409

	case errors.Is(err, app.ErrNotImplemented):
		return http.StatusNotImplemented // 501

	case errors.Is(err, app.ErrInternal):
		return http.StatusInternalServerError // 500

	default:
		return http.StatusInternalServerError // fallback
	}
}
