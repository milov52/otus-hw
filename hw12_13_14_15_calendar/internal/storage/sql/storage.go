package sqlstorage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/milov52/hw12_13_14_15_calendar/internal/config"
	"github.com/milov52/hw12_13_14_15_calendar/internal/storage"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Storage struct {
	dsn  string
	pool *pgxpool.Pool
}

func New() *Storage {
	return &Storage{
		dsn:  "",
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
	builderInsert := sq.Insert("events").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "title", "start_time", "description").
		Values(s.generateID(), event.Title, event.StartTime, event.Description).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to build query: %v", err)
	}

	var eventID uuid.UUID
	err = s.pool.QueryRow(ctx, query, args...).Scan(&eventID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert event: %v", err)
	}

	return eventID, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id uuid.UUID, event storage.Event) error {
	builderUpdate := sq.Update("events").
		PlaceholderFormat(sq.Dollar).
		Set("title", event.Title).
		Set("start_time", event.StartTime).
		Set("description", event.Description).
		Where(sq.Eq{"id": id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update event: %v", err)
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	builderDelete := sq.Delete("events").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}
	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete event: %v", err)
	}
	return nil
}

func (s *Storage) GetEvents(ctx context.Context, date time.Time, offset int) ([]storage.Event, error) {
	builderSelect := sq.Select("id", "title", "start_time", "description").
		From("events").
		PlaceholderFormat(sq.Dollar).
		Where(sq.GtOrEq{"start_time": date.AddDate(0, 0, offset)})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %v", err)
	}
	defer rows.Close()
	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.StartTime, &event.Description); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %v", err)
	}
	return events, nil
}
