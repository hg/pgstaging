package web

import (
	"github.com/hg/pgstaging/web/sessions"
	"github.com/hg/pgstaging/web/util"
	"github.com/hg/pgstaging/worker"
	"fmt"
	"net/http"
)

func serveCreate(rc *requestContext) {
	if !rc.isMethod(http.MethodPost) {
		return
	}

	name := rc.request.PostFormValue("name")
	name = util.NormalizeName(name)

	if name == "" || len(name) > 32 {
		rc.setResult(sessions.StatusError, fmt.Sprintf("некорректное имя '%s'", name))
		rc.redirect("/")
		return
	}

	name = util.AddPrefix(name)

	result := rc.srv.worker.Enqueue(worker.ActionCreate, name)
	go processResult(rc, result)

	rc.setResult(sessions.StatusQueued, "")
	rc.redirect("/")
}
