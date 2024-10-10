package integration

import (
	"log/slog"
	"os"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
	sqlstorage "github.com/milov52/hw12_13_14_15_calendar/internal/repository/event/sql"
	"github.com/milov52/hw12_13_14_15_calendar/internal/service/calendar"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type EventService interface {
	CreateEvent(ctx context.Context, event model.Event) (uuid.UUID, error)
	DayEventList(ctx context.Context, date time.Time) ([]model.Event, error)
	WeekEventList(ctx context.Context, date time.Time) ([]model.Event, error)
	MonthEventList(ctx context.Context, date time.Time) ([]model.Event, error)
}

type CalendarServiceSuite struct {
	suite.Suite
	svc  EventService
	pool *pgxpool.Pool
}

func TestServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(CalendarServiceSuite))
}

func (s *CalendarServiceSuite) SetupSuite() {
	pool, err := pgxpool.Connect(context.Background(), "postgres://test:test@localhost:5435/calendar")
	s.Require().NoError(err)

	repo := sqlstorage.New(pool)
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	s.svc = calendar.NewEventService(*logger, repo)
	s.pool = pool

}

func (s *CalendarServiceSuite) TearDownSuite() {
	_, _ = s.pool.Exec(context.Background(), "TRUNCATE TABLE event")
}

func (s *CalendarServiceSuite) TestCreateEvent() {
	const (
		eventTitle = "test create event"
	)

	startTime := time.Now()
	m := model.Event{
		Title:     eventTitle,
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "user1",
	}

	eventID, err := s.svc.CreateEvent(context.Background(), m)
	s.Require().NoError(err)
	s.Require().NotEmpty(eventID)

	createdEvent := s.getDirectItem(m.Title)

	s.Require().NoError(err)
	s.Require().NotEmpty(createdEvent)
	s.Require().Equal(eventID, createdEvent.ID)
	s.Require().Equal(m.Title, createdEvent.Title)
	s.Require().Equal(m.StartTime.Format(time.DateTime), createdEvent.StartTime.Format(time.DateTime))
	s.Require().Equal(m.Duration, createdEvent.Duration)
	s.Require().Equal(m.UserID, createdEvent.UserID)
}

func (s *CalendarServiceSuite) TestDayEventList() {
	startTime := time.Now().Truncate(24 * time.Hour).Add(9 * time.Hour)
	m1 := model.Event{
		Title:     "event 1",
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m2 := model.Event{
		Title:     "event 2",
		StartTime: startTime.Add(time.Hour * 2),
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m3 := model.Event{
		Title:     "event 3",
		StartTime: startTime.Add(time.Hour * 25),
		Duration:  time.Hour,
		UserID:    "user1",
	}
	_ = s.createDirectItem(m1)
	_ = s.createDirectItem(m2)
	_ = s.createDirectItem(m3)

	dayEvents, err := s.svc.DayEventList(context.Background(), time.Now())
	s.Require().NoError(err)
	s.Require().Equal(2, len(dayEvents))

}

func (s *CalendarServiceSuite) TestWeekEventList() {
	startTime := time.Now().Truncate(24 * time.Hour).Add(9 * time.Hour)
	m1 := model.Event{
		Title:     "event 1",
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m2 := model.Event{
		Title:     "event 2",
		StartTime: startTime.Add(time.Hour * 2),
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m3 := model.Event{
		Title:     "event 3",
		StartTime: startTime.Add(time.Hour * 25),
		Duration:  time.Hour,
		UserID:    "user1",
	}
	_ = s.createDirectItem(m1)
	_ = s.createDirectItem(m2)
	_ = s.createDirectItem(m3)

	dayEvents, err := s.svc.WeekEventList(context.Background(), time.Now())
	s.Require().NoError(err)
	s.Require().Equal(3, len(dayEvents))

}

func (s *CalendarServiceSuite) TestMonthEventList() {
	startTime := time.Now().Truncate(24 * time.Hour).Add(9 * time.Hour)
	m1 := model.Event{
		Title:     "event 1",
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m2 := model.Event{
		Title:     "event 2",
		StartTime: startTime.Add(time.Hour * 25),
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m3 := model.Event{
		Title:     "event 3",
		StartTime: startTime.Add(time.Hour * 24 * 30),
		Duration:  time.Hour,
		UserID:    "user1",
	}
	_ = s.createDirectItem(m1)
	_ = s.createDirectItem(m2)
	_ = s.createDirectItem(m3)

	dayEvents, err := s.svc.MonthEventList(context.Background(), time.Now())
	s.Require().NoError(err)
	s.Require().Equal(3, len(dayEvents))

}

func (s *CalendarServiceSuite) getDirectItem(title string) model.Event {
	query, args, err := sq.
		Select("id", "title", "start_time", "description", "duration", "user_id").
		From("event").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"title": title}).
		ToSql()

	if err != nil {
		s.Fail(err.Error())
	}

	rows, err := s.pool.Query(context.Background(), query, args...)
	if err != nil {
		s.Fail(err.Error())
	}
	defer rows.Close()

	item := model.Event{}
	for rows.Next() {
		scanErr := rows.Scan(&item.ID, &item.Title, &item.StartTime, &item.Description, &item.Duration, &item.UserID)
		if scanErr != nil {
			s.Fail(scanErr.Error())
		}
	}

	return item
}

func (s *CalendarServiceSuite) createDirectItem(event model.Event) uuid.UUID {
	query, args, err := sq.
		Insert("event").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "title", "start_time", "description", "duration", "notify_before", "user_id").
		Values(uuid.New().String(), event.Title, event.StartTime, event.Description,
			event.Duration, event.NotifyBefore, event.UserID).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		s.Fail(err.Error())
	}

	var eventID uuid.UUID
	err = s.pool.QueryRow(context.Background(), query, args...).Scan(&eventID)
	if err != nil {
		s.Fail(err.Error())
	}

	return eventID
}
