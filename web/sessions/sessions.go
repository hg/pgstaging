package sessions

import (
	"github.com/hg/pgstaging/util"
	"net/http"
	"sync"
)

type HttpHandler func(http.ResponseWriter, *http.Request)
type SessionHandler func(http.ResponseWriter, *http.Request, *Session)

type Sessions struct {
	m        sync.Mutex
	sessions map[string]*Session
}

func (s *Sessions) sessionId(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("session")
	if err == nil {
		return cookie.Value
	}
	cookie = &http.Cookie{
		Name:     "session",
		Value:    util.RandomString(16),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	return cookie.Value
}

func (s *Sessions) makeSession(w http.ResponseWriter, r *http.Request) *Session {
	id := s.sessionId(w, r)

	s.m.Lock()

	session := s.sessions[id]
	if session == nil {
		session = &Session{
			vals: make(map[string]interface{}),
		}
		s.sessions[id] = session
	}

	s.m.Unlock()
	return session
}

func (s *Sessions) Wrap(handler SessionHandler) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.makeSession(w, r)
		handler(w, r, session)
	}
}

func New() *Sessions {
	return &Sessions{
		sessions: make(map[string]*Session),
	}
}

func (s *Sessions) Get(w http.ResponseWriter, r *http.Request) *Session {
	return s.makeSession(w, r)
}
