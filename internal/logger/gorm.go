package logger

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

// SlogGormLogger — адаптер GORM Logger, использующий slog.Logger
type SlogGormLogger struct {
	logger     *slog.Logger
	level      logger.LogLevel
	slowThresh time.Duration
}

// New создаёт новый адаптер GORM Logger на slog
func NewSlogGormLogger(logger *slog.Logger, level logger.LogLevel) *SlogGormLogger {
	return &SlogGormLogger{
		logger:     logger,
		level:      level,
		slowThresh: 200 * time.Millisecond,
	}
}

// LogMode — смена уровня логирования (копия с новым уровнем)
func (l *SlogGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &SlogGormLogger{
		logger:     l.logger,
		level:      level,
		slowThresh: l.slowThresh,
	}
}

// Info логирует информационное сообщение
func (l *SlogGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		l.logger.InfoContext(ctx, msg, data...)
	}
}

// Warn логирует предупреждение
func (l *SlogGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Warn {
		l.logger.WarnContext(ctx, msg, data...)
	}
}

// Error логирует ошибку
func (l *SlogGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Error {
		l.logger.ErrorContext(ctx, msg, data...)
	}
}

// Trace логирует SQL-запрос
func (l *SlogGormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (string, int64),
	err error,
) {
	if l.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	attrs := []slog.Attr{
		slog.String("sql", sql),
		slog.Int64("rows", rows),
		slog.Duration("elapsed", elapsed),
	}

	switch {
	case err != nil && l.level >= logger.Error:
		attrs = append(attrs, slog.String("error", err.Error()))
		l.logger.LogAttrs(ctx, slog.LevelError, "gorm error", attrs...)

	case elapsed > l.slowThresh && l.level >= logger.Warn:
		l.logger.LogAttrs(ctx, slog.LevelWarn, "gorm slow query", attrs...)

	case l.level >= logger.Info:
		l.logger.LogAttrs(ctx, slog.LevelInfo, "gorm query", attrs...)
	}
}
