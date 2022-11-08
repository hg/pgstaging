package web

import (
	"github.com/hg/pgstaging/web/l10n"
	"github.com/hg/pgstaging/web/sessions"
	"github.com/hg/pgstaging/worker"
	"net/http"
)

type requestContext struct {
	srv     *server
	writer  http.ResponseWriter
	request *http.Request
	session *sessions.Session
	locale  *l10n.Locale
}

func (r *requestContext) requireMethod(method string) bool {
	if r.request.Method != method {
		http.Error(r.writer, "unsupported method "+r.request.Method, http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func (r *requestContext) redirect(url string) {
	http.Redirect(r.writer, r.request, url, http.StatusFound)
}

func (r *requestContext) addResult(status sessions.Status, message string) {
	r.session.AddEvent(status, message)
}

func (r *requestContext) bail(message string) {
	r.addResult(sessions.StatusError, message)
	r.redirect("/")
}

func processResult(req *requestContext, result <-chan error) {
	if err := <-result; err == nil {
		req.addResult(sessions.StatusSuccess, "")
	} else {
		req.addResult(sessions.StatusError, err.Error())
	}
}

func (r *requestContext) queueTask(action worker.Action, name, pass string) {
	result := r.srv.worker.Enqueue(action, name, pass)
	r.addResult(sessions.StatusQueued, "")
	r.redirect("/")
	go processResult(r, result)
}
