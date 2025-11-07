package httpserver

import "context"

// маленький безопасный хелпер
func requestIDFromContext(ctx context.Context) string {
	if v := ctx.Value(reqIDKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
