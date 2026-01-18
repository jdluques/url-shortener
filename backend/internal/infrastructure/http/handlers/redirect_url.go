package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/jdluques/url-shortener/internal/application/usecases"
	"github.com/jdluques/url-shortener/internal/infrastructure/logging"
)

type RedirectHandler struct {
	usecase usecases.RedirectUseCase
}

func NewRedirectHandler(usecase usecases.RedirectUseCase) *RedirectHandler {
	return &RedirectHandler{usecase: usecase}
}

func (handler *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context())

	code := chi.URLParam(r, "code")

	originalURL, err := handler.usecase.Execute(r.Context(), code)
	if err != nil {
		logger.Error("failed to redirect URL", zap.Error(err))
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
