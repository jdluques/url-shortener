package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jdluques/url-shortener/internal/application/usecases"
	"github.com/jdluques/url-shortener/internal/infrastructure/http"
	"github.com/jdluques/url-shortener/internal/infrastructure/http/handlers"
	"github.com/jdluques/url-shortener/internal/infrastructure/id"
	"github.com/jdluques/url-shortener/internal/infrastructure/logging"
	"github.com/jdluques/url-shortener/internal/infrastructure/postgres"
	"github.com/jdluques/url-shortener/internal/infrastructure/redis"
	"github.com/jdluques/url-shortener/internal/infrastructure/shortcode"
	"go.uber.org/zap"
)

func main() {
	env := os.Getenv("ENV")
	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")
	databaseURL := os.Getenv("DATABASE_URL")
	cacheAddr := os.Getenv("CACHE_ADDRESS")

	serverAddr := serverHost + ":" + serverPort

	logger, err := logging.NewLogger(env)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("starting url-shortener",
		zap.String("env", env),
		zap.String("addr", serverAddr),
	)

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		logger.Fatal("failed to open database", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	urlRepo := postgres.NewURLRepository(db)
	cache := redis.NewRedisCache(cacheAddr)

	idGen, err := id.NewSnowFlakeGenerator(1)
	if err != nil {
		logger.Fatal("failed to create id generator", zap.Error(err))
	}
	shortCodeGen := shortcode.NewBase62Generator()

	shortenURLUseCase := usecases.NewShortenURLUseCase(urlRepo, cache, idGen, shortCodeGen)

	shortenURLHandler := handlers.NewShortenURLHandler(*shortenURLUseCase)

	router := http.NewRouter(logger, shortenURLHandler)
	server := http.NewServer(router, serverAddr)

	logger.Info("http server started")

	if err := server.Start(); err != nil {
		logger.Fatal("http server stopped with error", zap.Error(err))
	}
}
