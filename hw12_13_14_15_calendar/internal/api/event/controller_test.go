package event_test

import (
	"context"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	event2 "github.com/milov52/hw12_13_14_15_calendar/internal/api/event"
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
	"github.com/milov52/hw12_13_14_15_calendar/internal/service/event"
	servicepb "github.com/milov52/hw12_13_14_15_calendar/pkg/api/event/v1"
)

type MockEventService struct {
	event.Service
}

func (m *MockEventService) CreateEvent(ctx context.Context, evt model.Event) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (m *MockEventService) UpdateEvent(ctx context.Context, id uuid.UUID, event model.Event) error {
	return nil
}

func (m *MockEventService) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *MockEventService) DayEventList(ctx context.Context, date time.Time) ([]model.Event, error) {
	return []model.Event{}, nil
}
func (m *MockEventService) WeekEventList(ctx context.Context, date time.Time) ([]model.Event, error) {
	return []model.Event{}, nil
}
func (m *MockEventService) MonthEventList(ctx context.Context, date time.Time) ([]model.Event, error) {
	return []model.Event{}, nil
}

func TestCreateEventGRPC(t *testing.T) {
	mockService := &MockEventService{}
	controller := event2.NewEventController(mockService)

	req := &servicepb.CreateRequest{
		Event: &servicepb.EventInfo{
			Title:     "Test Event",
			StartTime: timestamppb.New(time.Now()),
			Duration:  durationpb.New(time.Hour),
			UserId:    "user1",
		},
	}
	resp, err := controller.CreateEvent(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.UUID)
}

func TestUpdateEventGRPC(t *testing.T) {
	mockService := &MockEventService{}
	controller := event2.NewEventController(mockService)

	req := &servicepb.UpdateRequest{
		UUID: uuid.New().String(),
		Event: &servicepb.EventInfo{
			Title:    "Test Event2",
			Duration: durationpb.New(time.Hour * 2),
			UserId:   "user1",
		},
	}

	_, err := controller.UpdateEvent(context.Background(), req)
	require.NoError(t, err)
}

func TestDeleteEventGRPC(t *testing.T) {
	mockService := &MockEventService{}
	controller := event2.NewEventController(mockService)

	req := &servicepb.DeleteRequest{
		UUID: uuid.New().String(),
	}

	_, err := controller.DeleteEvent(context.Background(), req)
	require.NoError(t, err)
}

func TestDayEventList(t *testing.T) {
	mockService := &MockEventService{}
	controller := event2.NewEventController(mockService)

	req := &servicepb.GetRequest{
		Date: timestamppb.New(time.Now()),
	}
	resp, err := controller.GetDayEventList(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
}
