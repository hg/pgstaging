package sessions

import (
	"sync"
	"time"
)

type Event struct {
	Created time.Time
	Status  string
	Message string
}

type Session struct {
	m      sync.Mutex
	vals   map[string]string
	events []Event
}

func (s *Session) AddEvent(e Event) {
	s.m.Lock()
	s.events = append([]Event{e}, s.events...)
	s.m.Unlock()
}

func (s *Session) Events() []Event {
	s.m.Lock()
	defer s.m.Unlock()
	return s.events
}
