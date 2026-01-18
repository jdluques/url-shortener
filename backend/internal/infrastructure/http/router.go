package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"go.uber.org/zap"

	"github.com/jdluques/url-shortener/internal/infrastructure/http/handlers"
	"github.com/jdluques/url-shortener/internal/infrastructure/http/middleware"
)

func NewRouter(
	allowedOrigins []string,
	logger *zap.Logger,
	shortenURLHandler *handlers.ShortenURLHandler,
	redirectHandler *handlers.RedirectHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)

	r.Use(middleware.WithContextLogger(logger))
	r.Use(middleware.WithContextLogger(logger))

	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestSize(1 << 20))

	r.Use(chimw.NoCache)
	r.Use(chimw.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(chimw.SetHeader("X-Frame-Options", "DENY"))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.With(httprate.Limit(
			10,
			1*time.Minute,
			httprate.WithKeyByIP(),
		)).Post("/shorten", shortenURLHandler.ServeHTTP)

		r.With(httprate.Limit(
			300,
			1*time.Minute,
			httprate.WithKeyByIP(),
		)).Get("/{code}", redirectHandler.ServeHTTP)
		r.Post("/shorten", shortenURLHandler.ServeHTTP)
		r.Get("/{code}", redirectHandler.ServeHTTP)
	})

	return r
}
