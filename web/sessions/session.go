package sessions

import (
	"sync"
	"time"
)

const keepLastEvents = 25

type Status string

const (
	StatusError   = Status("error")
	StatusSuccess = Status("success")
	StatusQueued  = Status("queued")
)

type Event struct {
	Created time.Time
	Status  Status
	Message string
}

type Session struct {
	mu     sync.Mutex
	events []Event
}

func (s *Session) AddEvent(status Status, message string) {
	s.mu.Lock()
	e := Event{
		Created: time.Now(),
		Status:  status,
		Message: message,
	}
	s.events = append([]Event{e}, s.events...)
	if len(s.events) > keepLastEvents {
		s.events = s.events[:keepLastEvents]
	}
	defer s.mu.Unlock()
}

func (s *Session) Events() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Event{}, s.events...)
}
