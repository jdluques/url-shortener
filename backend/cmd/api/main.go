package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jdluques/url-shortener/internal/application/usecases"
	"github.com/jdluques/url-shortener/internal/infrastructure/http"
	"github.com/jdluques/url-shortener/internal/infrastructure/http/handlers"
	"github.com/jdluques/url-shortener/internal/infrastructure/id"
	"github.com/jdluques/url-shortener/internal/infrastructure/postgres"
	"github.com/jdluques/url-shortener/internal/infrastructure/redis"
	"github.com/jdluques/url-shortener/internal/infrastructure/shortcode"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	cacheAddr := os.Getenv("CACHE_ADDRESS")
	port := os.Getenv("PORT")

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	urlRepo := postgres.NewURLRepository(db)
	cache := redis.NewRedisCache(cacheAddr)

	idGen, err := id.NewSnowFlakeGenerator(1)
	if err != nil {
		log.Fatal(err)
	}
	shortCodeGen := shortcode.NewBase62Generator()

	shortenURLUseCase := usecases.NewShortenURLUseCase(urlRepo, cache, idGen, shortCodeGen)

	shortenURLHandler := handlers.NewShortenURLHandler(*shortenURLUseCase)

	router := http.NewRouter(shortenURLHandler)
	server := http.NewServer(router)

	log.Fatal(server.Start(":" + port))
}
