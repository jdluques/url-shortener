package logging

import (
	"context"

	"go.uber.org/zap"
)

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

func LoggerKey() loggerKeyType {
	return loggerKey
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.NewNop()
	}
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger
	}
	return zap.NewNop()
}
