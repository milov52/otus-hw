package sqlstorage

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/milov52/hw12_13_14_15_calendar/internal/config"
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
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

func (s *Storage) CreateEvent(ctx context.Context, event model.Event) (uuid.UUID, error) {
	const op = "repository.sql.CreateEvent"

	builderInsert := sq.Insert("event").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "title", "start_time", "description", "duration", "notify_before", "user_id").
		Values(s.generateID(), event.Title, event.StartTime, event.Description, event.Duration, event.NotifyBefore, event.UserID).
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

func (s *Storage) UpdateEvent(ctx context.Context, id uuid.UUID, event model.Event) error {
	const op = "repository.sql.UpdateEvent"

	builderUpdate := sq.Update("event").
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
	const op = "repository.sql.DeleteEvent"

	builderDelete := sq.Delete("event").
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

func (s *Storage) GetEvents(ctx context.Context, date time.Time, offset int) ([]model.Event, error) {
	const op = "repository.sql.GetEvents"

	startDate := date.Format(time.DateOnly)                     // Приводим к формату даты
	endDate := date.AddDate(0, 0, offset).Format(time.DateOnly) // Конечная дата

	builderSelect := sq.Select("id", "title", "start_time", "description").
		From("event").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Expr("start_time BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59"))

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var events []model.Event
	for rows.Next() {
		var event model.Event
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

func (s *Storage) GetNotifications(ctx context.Context, date time.Time) ([]model.Notification, error) {
	const op = "repository.sql.GetNotifications"

	// Преобразуем time.Time в строку в формате "YYYY-MM-DD HH:MM:SS"
	dateString := date.Format("2006-01-02 15:04:05")

	// Построение SQL-запроса с использованием правильных интервалов и дат
	builderSelect := sq.Select("id", "title", "start_time", "user_id").
		From("event").
		PlaceholderFormat(sq.Dollar).
		Where("notify_before IS NOT NULL").
		Where("sent = FALSE").
		Where("start_time - notify_before <= ?", dateString). // Здесь SQL обработает вычитание интервала
		Where("start_time > ?", dateString)

	// Преобразуем запрос в SQL строку
	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build SQL query: %w", op, err)
	}

	// Выполняем запрос к базе данных
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}
	defer rows.Close()

	var notifications []model.Notification

	// Проход по результатам
	for rows.Next() {
		var notification model.Notification
		if err := rows.Scan(&notification.EventID, &notification.Title, &notification.Date, &notification.UserID); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return notifications, nil
}
