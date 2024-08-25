package memorystorage

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/milov52/hw12_13_14_15_calendar/internal/storage"
	"golang.org/x/net/context"
)

type Storage struct {
	byDay  map[string][]storage.Event
	events map[uuid.UUID]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		byDay:  make(map[string][]storage.Event),
		events: make(map[uuid.UUID]storage.Event),
	}
}

func (s *Storage) generateID() uuid.UUID {
	return uuid.New()
}

func (s *Storage) addToIndex(event storage.Event) {
	dayKey := event.StartTime.Format(time.DateOnly)
	s.byDay[dayKey] = append(s.byDay[dayKey], event)
}

func (s *Storage) removeFromIndex(event storage.Event) {
	dayKey := event.StartTime.Format(time.DateOnly)
	s.byDay[dayKey] = removeEventFromSlice(s.byDay[dayKey], event.ID)
}

func (s *Storage) isExistEvent(event storage.Event) bool {
	dayKey := event.StartTime.Format(time.DateOnly)
	events, ok := s.byDay[dayKey]
	if ok {
		for _, item := range events {
			if item.StartTime == event.StartTime {
				return true
			}
		}
	}
	return false
}

func removeEventFromSlice(events []storage.Event, eventID uuid.UUID) []storage.Event {
	for i, e := range events {
		if e.ID == eventID {
			return append(events[:i], events[i+1:]...)
		}
	}
	return events
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (uuid.UUID, error) {
	if s.isExistEvent(event) {
		return uuid.Nil, storage.ErrDateBusy
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = s.generateID()
	s.events[event.ID] = event
	s.addToIndex(event)
	return event.ID, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id uuid.UUID, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldEvent, ok := s.events[id]
	if !ok {
		return storage.ErrEventNotFound
	}

	// Удаляем старую версию события из индексов
	s.removeFromIndex(oldEvent)

	s.events[id] = event
	s.addToIndex(event)
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[id]
	if !ok {
		return storage.ErrEventNotFound
	}

	s.removeFromIndex(event)
	delete(s.events, id)
	return nil
}

func (s *Storage) GetEvents(ctx context.Context, startDate time.Time, offset int) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []storage.Event
	for i := 0; i < offset; i++ {
		day := startDate.AddDate(0, 0, i)
		dayKey := day.Format(time.DateOnly)
		if dayEvents, ok := s.byDay[dayKey]; ok {
			events = append(events, dayEvents...)
		}
	}
	if len(events) == 0 {
		return events, storage.ErrEventNotFound
	}
	return events, nil
}
