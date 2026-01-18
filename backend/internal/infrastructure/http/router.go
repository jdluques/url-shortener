package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/jdluques/url-shortener/internal/infrastructure/http/handlers"
)

func NewRouter(
	shortenURLHandler *handlers.ShortenURLHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/shorten", shortenURLHandler.ServeHTTP)
	})

	return r
}
