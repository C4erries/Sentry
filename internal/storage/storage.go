package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	DSN string `mapstructure:"dsn"`
}

type Storage struct {
	DB     *sqlx.DB
	SQ     squirrel.StatementBuilderType
	Events *EventRepository
	Alerts *AlertRepository
}

func NewStorage(cfg *Config) (*Storage, error) {
	db, err := sqlx.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	s := &Storage{
		DB: db,
		SQ: sq,
	}

	s.Events = NewEventRepository(db, sq)
	s.Alerts = NewAlertRepository(db, sq)

	return s, nil
}
