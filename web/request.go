package web

import (
	"github.com/hg/pgstaging/web/l10n"
	"github.com/hg/pgstaging/web/sessions"
	"github.com/hg/pgstaging/worker"
	"github.com/hg/pgstaging/worker/command"
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

func (r *requestContext) processResult(result <-chan command.Result) {
	if re := <-result; re.Err == nil {
		r.addResult(sessions.StatusSuccess, "")
	} else {
		r.addResult(sessions.StatusError, re.Err.Error())
	}
}

func (r *requestContext) queueTask(action worker.Action, name, pass string) <-chan command.Result {
	result := r.srv.worker.Enqueue(action, name, pass)
	r.addResult(sessions.StatusQueued, "")
	return result
}
