package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const maxRetries = 5
const retryDelay = 2 * time.Second

func NewPostgresDatabaseConnection(logger *zap.Logger, databaseURL string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	attempt := 1
	for ; attempt <= maxRetries; attempt++ {
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			logger.Warn("failed to open database at attempt "+strconv.Itoa(attempt), zap.Error(err))
		}

		if err := db.Ping(); err != nil {
			logger.Warn("failed to ping database at attempt "+strconv.Itoa(attempt), zap.Error(err))
		}

		if attempt < maxRetries {
			backoff := time.Duration(attempt) * retryDelay
			logger.Info("retrying in " + backoff.String() + "...")
			time.Sleep(backoff)
		} else {
			return nil, fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
		}
	}

	return db, nil
}
