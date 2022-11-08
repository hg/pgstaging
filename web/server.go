package web

import (
	"github.com/hg/pgstaging/web/sessions"
	"github.com/hg/pgstaging/worker"
	"html/template"
	"net/http"
)

type server struct {
	sessions *sessions.Sessions
	worker   *worker.Client
	tpl      *template.Template
}

func (srv *server) wrap(handler func(*requestContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(&requestContext{
			srv:     srv,
			writer:  w,
			request: r,
			session: srv.sessions.Get(w, r),
		})
	}
}
