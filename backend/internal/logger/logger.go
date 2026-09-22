package logger

import (
	"errors"
	"log/slog"
	"os"

	"github.com/P3rCh1/immersive-images/backend/internal/config"
)

var (
	ErrInvalidLogLevel  = errors.New("invalid log level")
	ErrInvalidLogFormat = errors.New("invalid log format")
)

var levels = map[string]slog.Leveler{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func New() (*slog.Logger, error) {
	level, ok := levels[config.Config.Logger.Level]
	if !ok {
		return nil, ErrInvalidLogLevel
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch config.Config.Logger.Format {
	case "text":
		return slog.New(slog.NewTextHandler(os.Stdout, opts)), nil

	case "json":
		return slog.New(slog.NewJSONHandler(os.Stdout, opts)), nil
	}

	return nil, ErrInvalidLogFormat
}
