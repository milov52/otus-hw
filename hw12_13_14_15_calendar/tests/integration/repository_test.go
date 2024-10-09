//go:build integration

package integration

import (
	"log"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
	sqlstorage "github.com/milov52/hw12_13_14_15_calendar/internal/repository/event/sql"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type Repository interface {
	CreateEvent(ctx context.Context, event model.Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, event model.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvents(ctx context.Context, date time.Time, offset int) ([]model.Event, error)
}

type IntegrationSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	r    Repository
}

func TestStorageIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) SetupSuite() {
	pool, err := pgxpool.Connect(context.Background(), "postgres://test:test@pg:5432/calendar")
	if err != nil {
		log.Fatal(err)
	}
	s.pool = pool
}

func (s *IntegrationSuite) TearDownSuite() {
	_, _ = s.pool.Exec(context.Background(), "TRUNCATE TABLE event")
}

func (s *IntegrationSuite) SetupTest() {
	s.r = sqlstorage.New(s.pool)
}

func (s *IntegrationSuite) TestCreateEvent() {
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

	eventID, err := s.r.CreateEvent(context.Background(), m)
	s.Require().NoError(err)
	s.Require().NotEmpty(eventID)

	dbItem := s.getDirectItem(eventTitle)
	s.Require().NotEmpty(dbItem)
	s.Require().Equal(eventID, dbItem.ID)

	s.Require().Equal(m.StartTime.Format(time.DateTime), dbItem.StartTime.Format(time.DateTime))
	s.Require().Equal(m.Duration, dbItem.Duration)
	s.Require().Equal(m.UserID, dbItem.UserID)
}

func (s *IntegrationSuite) TestGetEvents() {
	const (
		eventTitle = "test get events"
	)

	startTime := time.Now()

	m := model.Event{
		Title:     eventTitle,
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "1000",
	}
	events, err := s.r.GetEvents(context.Background(), startTime, 0)
	oldLen := len(events)

	dbID := s.createDirectItem(m)

	events, err = s.r.GetEvents(context.Background(), startTime, 0)
	createdEvent := events[len(events)-1]
	s.Require().NoError(err)
	s.Require().NotEmpty(events)
	s.Require().Equal(oldLen+1, len(events))
	s.Require().Equal(dbID, createdEvent.ID)
	s.Require().Equal(m.StartTime.Format(time.DateTime), createdEvent.StartTime.Format(time.DateTime))
	s.Require().Equal(m.Duration, createdEvent.Duration)
	s.Require().Equal(m.UserID, createdEvent.UserID)
}

func (s *IntegrationSuite) TestGetNotFound() {
	const (
		eventTitle = "test not found events"
	)

	startTime := time.Now()

	m := model.Event{
		Title:     eventTitle,
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "1000",
	}
	_ = s.createDirectItem(m)

	events, err := s.r.GetEvents(context.Background(), startTime.Add(time.Hour*60), 0)
	s.Require().NoError(err)
	s.Require().Empty(events)
}

func (s *IntegrationSuite) TestUpdateEvent() {
	const (
		eventTitleOld = "old event title"
		eventTitleNew = "new event title"
	)

	startTime := time.Now()
	m1 := model.Event{
		Title:       eventTitleOld,
		Description: "some description",
		StartTime:   startTime,
		Duration:    time.Hour,
		UserID:      "1000",
	}
	m2 := model.Event{
		Title:       eventTitleNew,
		Description: "updated description",
		StartTime:   startTime,
		Duration:    time.Hour,
	}

	dbID := s.createDirectItem(m1)
	err := s.r.UpdateEvent(context.Background(), dbID, m2)
	updatedEvent := s.getDirectItem(m2.Title)

	s.Require().NoError(err)
	s.Require().Equal(dbID, updatedEvent.ID)
	s.Require().Equal(m2.Title, updatedEvent.Title)
	s.Require().Equal(m2.Duration, updatedEvent.Duration)
	s.Require().Equal(m2.Description, updatedEvent.Description)
}

func (s *IntegrationSuite) TestDeleteEvent() {
	const (
		eventTitle = "test delete event"
	)
	startTime := time.Now()
	m := model.Event{
		Title:     eventTitle,
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "1000",
	}
	dbID := s.createDirectItem(m)
	err := s.r.DeleteEvent(context.Background(), dbID)
	s.Require().NoError(err)

	events := s.getDirectItem(eventTitle)
	s.Require().Empty(events)
}

func (s *IntegrationSuite) createDirectItem(event model.Event) uuid.UUID {
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

func (s *IntegrationSuite) getDirectItem(title string) model.Event {
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
