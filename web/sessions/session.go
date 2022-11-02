package sessions

import (
	"sync"
)

type Session struct {
	m    sync.Mutex
	vals map[string]interface{}
}

func (s *Session) Get(key string) interface{} {
	s.m.Lock()
	defer s.m.Unlock()
	return s.vals[key]
}

func (s *Session) Set(key string, value interface{}) {
	s.m.Lock()
	s.vals[key] = value
	s.m.Unlock()
}
