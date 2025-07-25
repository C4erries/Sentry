package logger

import (
	"fmt"
	"log/slog"
	"strings"

	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
)

type Config struct {
	Level string `mapstructure:"level"`
}

func Init(cfg *Config) error {
	level := parseLevel(cfg.Level)

	var zapLogger *zap.Logger
	var err error

	// Включаем zap development mode только для debug
	if level <= slog.LevelDebug {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		return fmt.Errorf("failed to init zap: %w", err)
	}

	options := slogzap.Option{
		Level:  level,
		Logger: zapLogger,
	}

	slog.SetDefault(slog.New(options.NewZapHandler()))

	return nil
}

// parseLevel преобразует строку в slog.Level
func parseLevel(levelStr string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(levelStr)) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
