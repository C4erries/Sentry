package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/c4erries/Sentry/internal/model"
)

type AlertRepository struct {
	db *sql.DB
	sq squirrel.StatementBuilderType
}

func NewAlertRepository(db *sql.DB, sq squirrel.StatementBuilderType) *AlertRepository {
	return &AlertRepository{db: db, sq: sq}
}

func (r *AlertRepository) Save(ctx context.Context, alert *model.Alert) error {
	data, err := json.Marshal(alert.Data)
	if err != nil {
		return fmt.Errorf("marshal alert data: %w", err)
	}

	query, args, err := squirrel.
		Insert("alerts").
		Columns("id", "rule", "level", "detected_at", "data").
		Values(alert.ID, alert.Rule, alert.Level, alert.DetectedAt, data).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert alert: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec insert alert: %w", err)
	}

	return nil
}
