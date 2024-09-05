package memorystorage

import (
	"testing"
	"time"

	"github.com/milov52/hw12_13_14_15_calendar/internal/storage"
	"golang.org/x/net/context"
)

func TestStorage_CreateEvent(t *testing.T) {
	testStorage := New()

	event := storage.Event{
		Title:     "Test Event",
		StartTime: time.Now(),
		Duration:  time.Hour,
		UserID:    "user1",
	}

	_, err := testStorage.CreateEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(testStorage.events) != 1 {
		t.Fatalf("expected 1 event in storage, got %d", len(testStorage.events))
	}

	dayKey := event.StartTime.Format(time.DateOnly)
	if len(testStorage.byDay[dayKey]) != 1 {
		t.Fatalf("expected 1 event in day index, got %d", len(testStorage.byDay[dayKey]))
	}
}

func TestStorage_UpdateEvent(t *testing.T) {
	testStorage := New()

	event := storage.Event{
		Title:     "Test Event",
		StartTime: time.Now(),
		Duration:  time.Hour,
		UserID:    "user1",
	}

	id, err := testStorage.CreateEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedEvent := storage.Event{
		Title:     "Updated Event",
		StartTime: event.StartTime,
		Duration:  2 * time.Hour,
		UserID:    "user1",
	}

	err = testStorage.UpdateEvent(context.Background(), id, updatedEvent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	storedEvent := testStorage.events[id]
	if err != nil {
		t.Fatalf("expected event with ID %s to exist", event.ID)
	}

	if storedEvent.Title != "Updated Event" {
		t.Errorf("expected title 'Updated Event', got %s", storedEvent.Title)
	}

	dayKey := updatedEvent.StartTime.Format(time.DateOnly)
	if len(testStorage.byDay[dayKey]) != 1 {
		t.Fatalf("expected 1 event in day index after update, got %d", len(testStorage.byDay[dayKey]))
	}
}

func TestStorage_DeleteEvent(t *testing.T) {
	testStorage := New()

	event := storage.Event{
		Title:     "Test Event",
		StartTime: time.Now(),
		Duration:  time.Hour,
		UserID:    "user1",
	}

	id, err := testStorage.CreateEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = testStorage.DeleteEvent(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(testStorage.events) != 0 {
		t.Fatalf("expected 0 events in storage after deletion, got %d", len(testStorage.events))
	}

	dayKey := event.StartTime.Format(time.DateOnly)
	if len(testStorage.byDay[dayKey]) != 0 {
		t.Fatalf("expected 0 events in day index after deletion, got %d", len(testStorage.byDay[dayKey]))
	}
}

func TestStorage_GetEvents(t *testing.T) {
	testStorage := New()

	// Создаем несколько событий
	event1 := storage.Event{
		Title:     "Event 1",
		StartTime: time.Now(),
		Duration:  time.Hour,
		UserID:    "user1",
	}
	event2 := storage.Event{
		Title:     "Event 2",
		StartTime: time.Now().AddDate(0, 0, 1), // Завтра
		Duration:  2 * time.Hour,
		UserID:    "user2",
	}
	event3 := storage.Event{
		Title:     "Event 3",
		StartTime: time.Now().AddDate(0, 0, 2), // Послезавтра
		Duration:  time.Hour,
		UserID:    "user3",
	}
	ctx := context.Background()
	_, _ = testStorage.CreateEvent(ctx, event1)
	_, _ = testStorage.CreateEvent(ctx, event2)
	_, _ = testStorage.CreateEvent(ctx, event3)

	// Получаем события на 2 дня вперед
	events, err := testStorage.GetEvents(ctx, time.Now(), 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	// Проверяем, что события возвращаются в правильном порядке
	if events[0].Title != "Event 1" || events[1].Title != "Event 2" {
		t.Errorf("events not returned in expected order")
	}
}
