package storage

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type Storage struct {
	DB     *sql.DB
	SQ     squirrel.StatementBuilderType
	Events *EventRepository
	Alerts *AlertRepository
}

func NewStorage(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
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
