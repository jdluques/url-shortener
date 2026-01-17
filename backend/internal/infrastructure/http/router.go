package http

import (
	"net/http"

	"github.com/jdluques/url-shortener/internal/infrastructure/http/handlers"
)

func NewRouter(
	shortenURLHandler *handlers.ShortenURLHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("api/v1/shorten", shortenURLHandler)

	return mux
}
