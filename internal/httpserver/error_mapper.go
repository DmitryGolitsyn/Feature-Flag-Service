package httpserver

import (
	"errors"
	"ffs-tutorial/internal/app"
	"net/http"
)

func MapErrorToStatus(err error) int {
	switch {

	case errors.Is(err, app.ErrInvalid):
		return http.StatusBadRequest // 400

	case errors.Is(err, app.ErrForbidden):
		return http.StatusForbidden // 403

	case errors.Is(err, app.ErrNotFound):
		return http.StatusNotFound // 404

	case errors.Is(err, app.ErrInternal):
		return http.StatusInternalServerError // 500

	default:
		return http.StatusInternalServerError // fallback
	}
}
