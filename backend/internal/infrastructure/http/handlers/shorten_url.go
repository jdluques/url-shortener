package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jdluques/url-shortener/internal/application/usecases"
)

type ShortenURLHandler struct {
	usecase usecases.ShortenURLUseCase
}

func NewShortenURLHandler(usecase usecases.ShortenURLUseCase) *ShortenURLHandler {
	return &ShortenURLHandler{usecase: usecase}
}

func (handler *ShortenURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string        `json:"url"`
		TTL time.Duration `json:"ttl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	shortURL, err := handler.usecase.Execute(ctx, req.URL, req.TTL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"short_code": shortURL.ShortCode(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
