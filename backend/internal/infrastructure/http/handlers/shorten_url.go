package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/jdluques/url-shortener/internal/application/usecases"
	"github.com/jdluques/url-shortener/internal/infrastructure/logging"
)

type ShortenURLHandler struct {
	usecase usecases.ShortenURLUseCase
}

func NewShortenURLHandler(usecase usecases.ShortenURLUseCase) *ShortenURLHandler {
	return &ShortenURLHandler{usecase: usecase}
}

func (handler *ShortenURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context())

	var req struct {
		OriginalURL string        `json:"originalUrl"`
		TTL         time.Duration `json:"ttl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("invalid request body", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	shortURL, err := handler.usecase.Execute(ctx, req.OriginalURL, req.TTL)
	if err != nil {
		logger.Error("failed to shorten url", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"short_code": shortURL.ShortCode(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("failed to encode response", zap.Error(err))
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
