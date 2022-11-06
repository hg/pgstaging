package web

import (
	"github.com/hg/pgstaging/web/sessions"
	"github.com/hg/pgstaging/web/util"
	"github.com/hg/pgstaging/worker"
	"fmt"
	"net/http"
)

func serveModify(rc *requestContext) {
	if !rc.requireMethod(http.MethodPost) {
		return
	}

	name := util.NormalizeName(rc.request.PostFormValue("name"))

	if !util.IsDevName(name) {
		rc.setResult(sessions.StatusError, fmt.Sprintf("некорректное имя '%s'", name))
		rc.redirect("/")
		return
	}

	action := rc.request.FormValue("action")

	var result <-chan error

	switch action {
	case "start":
		result = rc.srv.worker.Enqueue(worker.ActionStart, name)
	case "stop":
		result = rc.srv.worker.Enqueue(worker.ActionStop, name)
	case "drop":
		result = rc.srv.worker.Enqueue(worker.ActionDrop, name)
	default:
		rc.setResult(sessions.StatusError, fmt.Sprintf("неизвестное действие '%s'", action))
	}

	if result != nil {
		go processResult(rc, result)
		rc.setResult(sessions.StatusQueued, "")
	}

	rc.redirect("/")
}
