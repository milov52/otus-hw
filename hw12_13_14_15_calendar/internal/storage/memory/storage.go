package memorystorage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/milov52/hw12_13_14_15_calendar/internal/storage"
	"sync"
	"time"
)

type Storage struct {
	byDay  map[string][]storage.Event
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		byDay:  make(map[string][]storage.Event),
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) generateID() string {
	return uuid.New().String()
}

func (s *Storage) addToIndex(event storage.Event) {
	dayKey := event.StartTime.Format("2006-01-02")
	s.byDay[dayKey] = append(s.byDay[dayKey], event)
}

func (s *Storage) removeFromIndex(event storage.Event) {
	dayKey := event.StartTime.Format("2006-01-02")
	s.byDay[dayKey] = removeEventFromSlice(s.byDay[dayKey], event.ID)
}

func removeEventFromSlice(events []storage.Event, eventID string) []storage.Event {
	for i, e := range events {
		if e.ID == eventID {
			return append(events[:i], events[i+1:]...)
		}
	}
	return events
}

func (s *Storage) CreateEvent(event storage.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = s.generateID()
	s.events[event.ID] = event
	s.addToIndex(event)
	return event.ID, nil
}

func (s *Storage) UpdateEvent(id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldEvent, ok := s.events[id]
	if !ok {
		return fmt.Errorf("event with id not found")
	}

	// Удаляем старую версию события из индексов
	s.removeFromIndex(oldEvent)

	s.events[id] = event
	s.addToIndex(event)
	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[id]
	if !ok {
		return fmt.Errorf("event with id not foundd")
	}

	s.removeFromIndex(event)
	delete(s.events, id)
	return nil
}

func (s *Storage) GetEvents(startDate time.Time, offset int) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []storage.Event
	for i := 0; i < offset; i++ {
		day := startDate.AddDate(0, 0, i)
		dayKey := day.Format("2006-01-02")
		if dayEvents, ok := s.byDay[dayKey]; ok {
			events = append(events, dayEvents...)
		}
	}

	return events, nil
}
