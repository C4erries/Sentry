package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/c4erries/Sentry/internal/model"
	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
)

type AlertRepository struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func NewAlertRepository(db *sqlx.DB, sq squirrel.StatementBuilderType) *AlertRepository {
	return &AlertRepository{db: db, sq: sq}
}

func (r *AlertRepository) Save(ctx context.Context, alert *model.Alert) error {
	var err error
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	data, err := json.Marshal(alert.Data)
	if err != nil {
		return fmt.Errorf("marshal alert data: %w", err)
	}

	query, args, err := squirrel.
		Insert("alerts").
		Columns("id", "rule", "level", "event_id", "detected_at", "data").
		Values(alert.ID, alert.Rule, alert.Level, alert.Event.ID, alert.DetectedAt, data).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert alert: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec insert alert: %w", err)
	}
	/*
		builder := squirrel.Insert("alert_events").
			Columns("alert_id", "event_id").
			PlaceholderFormat(squirrel.Dollar)

		builder = builder.Values(alert.ID, alert.Event.ID)

		query, args, err = builder.ToSql()
		if err != nil {
			return fmt.Errorf("failed to build insert query: %w", err)
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to execute insert query: %w", err)
		}
	*/
	return err
}

func (r *AlertRepository) AllEventsExist(ctx context.Context, eventIDs []string) (bool, error) {
	if len(eventIDs) == 0 {
		return true, nil // нет требований
	}

	// Преобразуем []string -> []any для передачи в $1::uuid[]
	args := make([]any, len(eventIDs))
	for i, id := range eventIDs {
		args[i] = id
	}

	// Строим SQL-запрос
	query := `
		SELECT COUNT(*) 
		FROM events 
		WHERE id = ANY($1::uuid[])
	`

	var found int
	err := r.db.QueryRowContext(ctx, query, args).Scan(&found)
	if err != nil {
		return false, fmt.Errorf("checking events exist: %w", err)
	}

	// Сравниваем: если найдено меньше, чем ожидаем, значит чего-то не хватает
	return found == len(eventIDs), nil
}

func (r *AlertRepository) FindByID(ctx context.Context, alertID string) (*model.Alert, error) {
	var alert model.Alert

	queryBuilder := r.sq.Select("Distinct a.*").
		From("alerts a").
		Where(squirrel.Eq{"a.id": alertID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.SelectContext(ctx, &alert, query, args...)

	return &alert, err
}

// TODO: json_agg
func (r *AlertRepository) FindByUser(ctx context.Context, userID string) ([]*model.Alert, error) {
	builder := squirrel.Select().
		Columns(
			"a.id            AS alert_id",
			"a.rule          AS alert_rule",
			"a.level         AS alert_level",
			"a.detected_at   AS alert_detected_at",
			"a.data          AS alert_data",

			// event
			"e.id            AS event_id",
			"e.user_id       AS event_user_id",
			"e.event_type    AS event_type",
			"e.created_at    AS event_timestamp",
			"e.ip            AS event_ip",
			"e.geo_country   AS event_geo_country",
			"e.data          AS event_data",
		).
		From("alerts a").
		Join("events e ON e.id = a.event_id").
		Where(squirrel.Eq{"e.user_id": userID}).
		OrderBy("a.detected_at DESC").
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute joined alert query: %w", err)
	}
	defer rows.Close()

	var result []*model.Alert
	var parseErrors *multierror.Error

	for rows.Next() {
		var raw AlertDB
		err := rows.Scan(
			&raw.ID,
			&raw.Rule,
			&raw.Level,
			&raw.DetectedAt,
			&raw.RawData,

			&raw.Event.ID,
			&raw.Event.UserID,
			&raw.Event.EventType,
			&raw.Event.Timestamp,
			&raw.Event.IP,
			&raw.Event.GeoCountry,
			&raw.Event.RawData,
		)
		if err != nil {
			parseErrors = multierror.Append(parseErrors, fmt.Errorf("failed to scan row: %w", err))
			continue
		}
		alert, err := raw.ToAlert()
		if err != nil {
			parseErrors = multierror.Append(parseErrors, fmt.Errorf("failed to convert raw alert: %w", err))
			continue
		}
		result = append(result, &alert)
	}

	slog.Debug("[FindByUser]", slog.Int("count", len(result)), slog.String("userID", userID))
	return result, parseErrors.ErrorOrNil()
}
