package storage

import (
	"fmt"
	"log/slog"

	"github.com/c4erries/Sentry/internal/logger"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	lg "gorm.io/gorm/logger"
)

type Config struct {
	DSN string `mapstructure:"dsn"`
}

type Storage struct {
	DB     *gorm.DB
	Events *EventRepository
	Alerts *AlertRepository
}

func NewStorage(cfg *Config) (*Storage, error) {
	gormLogger := logger.NewSlogGormLogger(slog.Default(), lg.Info)
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	s := &Storage{
		DB: db,
	}

	s.Events = NewEventRepository(db)
	s.Alerts = NewAlertRepository(db)

	return s, nil
}

// TODO?: реализовать получение sql.DB раньше, а тут лишь вызывать закрытие
func (s *Storage) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get *sql.DB: %v", err)
	}
	return sqlDB.Close()
}
