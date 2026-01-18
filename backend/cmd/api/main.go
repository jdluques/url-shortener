package main

import (
	"database/sql"
	"log"
	"strconv"

	"go.uber.org/zap"

	"github.com/jdluques/url-shortener/internal/application/usecases"
	"github.com/jdluques/url-shortener/internal/config"
	"github.com/jdluques/url-shortener/internal/infrastructure/http"
	"github.com/jdluques/url-shortener/internal/infrastructure/http/handlers"
	"github.com/jdluques/url-shortener/internal/infrastructure/id"
	"github.com/jdluques/url-shortener/internal/infrastructure/logging"
	"github.com/jdluques/url-shortener/internal/infrastructure/postgres"
	"github.com/jdluques/url-shortener/internal/infrastructure/redis"
	"github.com/jdluques/url-shortener/internal/infrastructure/shortcode"
)

func main() {
	envVars, err := config.LoadEnvVars()
	if err != nil {
		log.Fatalf("failed to load env vars: %v", err)
	}

	serverAddr := envVars.ServerHost + ":" + strconv.Itoa(envVars.ServerPort)

	logger, err := logging.NewLogger(envVars.Env)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("starting url-shortener",
		zap.String("env", envVars.Env),
		zap.String("addr", serverAddr),
	)

	db, err := sql.Open("postgres", envVars.DatabaseSource)
	if err != nil {
		logger.Fatal("failed to open database", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	urlRepo := postgres.NewURLRepository(db)
	cache := redis.NewRedisCache(envVars.CacheAddress)

	idGen, err := id.NewSnowFlakeGenerator(1)
	if err != nil {
		logger.Fatal("failed to create id generator", zap.Error(err))
	}
	shortCodeGen := shortcode.NewBase62Generator()

	shortenURLUseCase := usecases.NewShortenURLUseCase(urlRepo, cache, idGen, shortCodeGen)
	redirectUseCase := usecases.NewRedirectUseCase(urlRepo, cache)

	shortenURLHandler := handlers.NewShortenURLHandler(*shortenURLUseCase)
	redirectHandler := handlers.NewRedirectHandler(*redirectUseCase)

	router := http.NewRouter(envVars.AllowedOrigins, logger, shortenURLHandler, redirectHandler)
	server := http.NewServer(router, serverAddr)

	logger.Info("http server started")

	if err := server.Start(); err != nil {
		logger.Fatal("http server stopped with error", zap.Error(err))
	}
}
