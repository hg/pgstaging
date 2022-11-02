package sessions

import (
	"sync"
)

type Session struct {
	m    sync.Mutex
	vals map[string]string
}

func (s *Session) Get(key string) string {
	s.m.Lock()
	defer s.m.Unlock()
	return s.vals[key]
}

func (s *Session) Set(key string, value string) {
	s.m.Lock()
	s.vals[key] = value
	s.m.Unlock()
}
