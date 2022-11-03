package sessions

import (
	"sync"
	"time"
)

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
	m      sync.Mutex
	vals   map[string]string
	events []Event
}

func (s *Session) AddEvent(status Status, message string) {
	s.m.Lock()
	e := Event{
		Created: time.Now(),
		Status:  status,
		Message: message,
	}
	s.events = append([]Event{e}, s.events...)
	s.m.Unlock()
}

func (s *Session) Events() []Event {
	s.m.Lock()
	defer s.m.Unlock()
	return s.events
}
