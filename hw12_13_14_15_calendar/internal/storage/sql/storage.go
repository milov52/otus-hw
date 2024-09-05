package sqlstorage

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/milov52/hw12_13_14_15_calendar/internal/config"
	"github.com/milov52/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New() *Storage {
	return &Storage{
		pool: nil,
	}
}

func (s *Storage) Connect(ctx context.Context, cfg config.Config) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	s.pool = pool
	return s.pool.Ping(ctx)
}

func (s *Storage) Close(ctx context.Context) {
	s.pool.Close()
}

func (s *Storage) generateID() string {
	return uuid.New().String()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (uuid.UUID, error) {
	const op = "storage.sql.CreateEvent"

	builderInsert := sq.Insert("events").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "title", "start_time", "description").
		Values(s.generateID(), event.Title, event.StartTime, event.Description).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	var eventID uuid.UUID
	err = s.pool.QueryRow(ctx, query, args...).Scan(&eventID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return eventID, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id uuid.UUID, event storage.Event) error {
	const op = "storage.sql.UpdateEvent"

	builderUpdate := sq.Update("events").
		PlaceholderFormat(sq.Dollar).
		Set("title", event.Title).
		Set("start_time", event.StartTime).
		Set("description", event.Description).
		Where(sq.Eq{"id": id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	const op = "storage.sql.DeleteEvent"

	builderDelete := sq.Delete("events").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetEvents(ctx context.Context, date time.Time, offset int) ([]storage.Event, error) {
	const op = "storage.sql.GetEvents"

	startDate := date.Format(time.DateOnly)                     // Приводим к формату даты
	endDate := date.AddDate(0, 0, offset).Format(time.DateOnly) // Конечная дата

	builderSelect := sq.Select("id", "title", "start_time", "description").
		From("events").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Expr("start_time BETWEEN ? AND ?", startDate, endDate))

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.StartTime, &event.Description); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return events, nil
}
