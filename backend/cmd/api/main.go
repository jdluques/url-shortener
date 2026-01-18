package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

	db, err := postgres.NewPostgresDatabaseConnection(logger, envVars.DatabaseSource)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	urlRepo := postgres.NewURLRepository(db)
	cache := redis.NewRedisCache(envVars.CacheAddress)

	idGen, err := id.NewSnowFlakeGenerator(envVars.NodeId)
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

	go func() {
		logger.Info("http server started")
		if err := server.Start(); err != nil {
			logger.Fatal("http server stopped with error", zap.Error(err))
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	logger.Info("shutting server down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server shutdown failed", zap.Error(err))
	}

	logger.Info("server gracefully stopped")
}
