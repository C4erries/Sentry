package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/c4erries/Sentry/internal/model"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func NewEventRepository(db *sqlx.DB, sq squirrel.StatementBuilderType) *EventRepository {
	return &EventRepository{db: db, sq: sq}
}

func (r *EventRepository) Save(ctx context.Context, event *model.Event) error {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("marshal event data: %w", err)
	}

	query, args, err := squirrel.
		Insert("events").
		Columns("id", "user_id", "event_type", "created_at", "geo_country", "ip", "data").
		Values(event.ID, event.UserID, event.EventType.String(), event.Timestamp, event.GeoCountry, event.IP, data).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert event: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec insert event: %w", err)
	}

	return nil
}
