package middleware

import (
	"context"
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/jdluques/url-shortener/internal/infrastructure/logging"
)

func WithContextLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqLogger := logger.With(
				zap.String("request_id", chimw.GetReqID(r.Context())),
			)

			ctx := context.WithValue(r.Context(), logging.LoggerKey(), reqLogger)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
