package sessions

import (
	"github.com/hg/pgstaging/util"
	"net/http"
	"sync"
)

type SessionID string

type Sessions struct {
	mu       sync.Mutex
	sessions map[SessionID]*Session
}

func (s *Sessions) sessionId(w http.ResponseWriter, r *http.Request) SessionID {
	cookie, err := r.Cookie("session")
	if err == nil {
		return SessionID(cookie.Value)
	}
	cookie = &http.Cookie{
		Name:     "session",
		Value:    util.RandomString(16),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	return SessionID(cookie.Value)
}

func (s *Sessions) makeSession(w http.ResponseWriter, r *http.Request) *Session {
	id := s.sessionId(w, r)

	s.mu.Lock()
	defer s.mu.Unlock()

	session := s.sessions[id]
	if session == nil {
		session = &Session{}
		s.sessions[id] = session
	}

	return session
}

func New() *Sessions {
	return &Sessions{
		sessions: make(map[SessionID]*Session),
	}
}

func (s *Sessions) Get(w http.ResponseWriter, r *http.Request) *Session {
	return s.makeSession(w, r)
}
