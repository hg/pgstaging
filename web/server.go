package web

import (
	"github.com/hg/pgstaging/web/l10n"
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
		langs := toIsoCodes(r.Header.Get("Accept-Language"))

		handler(&requestContext{
			srv:     srv,
			writer:  w,
			request: r,
			session: srv.sessions.Get(w, r),
			locale:  l10n.TryLangs(langs),
		})
	}
}
