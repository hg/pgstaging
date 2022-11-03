package web

import (
	"github.com/hg/pgstaging/web/sessions"
	"net/http"
)

type requestContext struct {
	srv     *server
	writer  http.ResponseWriter
	request *http.Request
	session *sessions.Session
}

func (r *requestContext) isMethod(method string) bool {
	if r.request.Method != method {
		http.Error(r.writer, "unsupported method "+r.request.Method, http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func (r *requestContext) redirect(url string) {
	http.Redirect(r.writer, r.request, url, http.StatusFound)
}

func (r *requestContext) setResult(status sessions.Status, message string) {
	r.session.AddEvent(status, message)
}
