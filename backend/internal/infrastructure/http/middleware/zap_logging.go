package middleware

import (
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func ZapLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
			h.ServeHTTP(ww, r)

			logger.Info("http_request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Int("bytes", ww.BytesWritten()),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_ip", r.RemoteAddr),
				zap.String("request_id", chimw.GetReqID(r.Context())),
			)
		})
	}
}
