package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/jdluques/url-shortener/internal/infrastructure/http/handlers"
	"github.com/jdluques/url-shortener/internal/infrastructure/http/middleware"
)

func NewRouter(
	logger *zap.Logger,
	shortenURLHandler *handlers.ShortenURLHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)

	r.Use(middleware.WithContextLogger(logger))
	r.Use(middleware.WithContextLogger(logger))

	r.Use(chimw.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/shorten", shortenURLHandler.ServeHTTP)
	})

	return r
}
