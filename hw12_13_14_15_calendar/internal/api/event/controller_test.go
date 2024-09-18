package event_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	event2 "github.com/milov52/hw12_13_14_15_calendar/internal/api/event"
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
	"github.com/milov52/hw12_13_14_15_calendar/internal/service/event"
	servicepb "github.com/milov52/hw12_13_14_15_calendar/pkg/api/event/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreateEvent(ctx context.Context, evt model.Event) (uuid.UUID, error) {
	args := m.Called(ctx, evt)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockStorage) UpdateEvent(ctx context.Context, id uuid.UUID, evt model.Event) error {
	args := m.Called(ctx, id, evt)
	return args.Error(0)
}

func (m *MockStorage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStorage) GetEvents(ctx context.Context, date time.Time, offset int) ([]model.Event, error) {
	args := m.Called(ctx, date, offset)
	return args.Get(0).([]model.Event), args.Error(1)
}

func TestCreateEventGRPC(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mockRepo := new(MockStorage)
	mockService := event.NewEventService(*logger, mockRepo)
	controller := event2.NewEventController(mockService)

	mockRepo.On("CreateEvent", mock.Anything,
		mock.AnythingOfType("model.Event")).Return(uuid.New(), nil)

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

	mockRepo.AssertExpectations(t)
}

func TestUpdateEventGRPC(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mockRepo := new(MockStorage)
	mockService := event.NewEventService(*logger, mockRepo)
	controller := event2.NewEventController(mockService)

	mockRepo.On("UpdateEvent", mock.Anything, mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("model.Event")).Return(nil)

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

	mockRepo.AssertExpectations(t)
}

func TestDeleteEventGRPC(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mockRepo := new(MockStorage)
	mockService := event.NewEventService(*logger, mockRepo)
	controller := event2.NewEventController(mockService)

	mockRepo.On("DeleteEvent", mock.Anything,
		mock.AnythingOfType("uuid.UUID")).Return(nil)

	req := &servicepb.DeleteRequest{
		UUID: uuid.New().String(),
	}

	_, err := controller.DeleteEvent(context.Background(), req)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDayEventList(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mockRepo := new(MockStorage)
	mockService := event.NewEventService(*logger, mockRepo)
	controller := event2.NewEventController(mockService)

	mockRepo.On("GetEvents", mock.Anything,
		mock.AnythingOfType("time.Time"), mock.AnythingOfType("int")).Return([]model.Event{}, nil)

	req := &servicepb.GetRequest{
		Date: timestamppb.New(time.Now()),
	}
	resp, err := controller.GetDayEventList(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	mockRepo.AssertExpectations(t)
}
