package web

import (
	"github.com/hg/pgstaging/web/sessions"
	"github.com/hg/pgstaging/web/util"
	"github.com/hg/pgstaging/worker"
	"fmt"
	"net/http"
)

func serveCreate(rc *requestContext) {
	if !rc.requireMethod(http.MethodPost) {
		return
	}

	name := rc.request.PostFormValue("name")
	name = util.NormalizeName(name)

	if name == "" || len(name) > 32 {
		rc.setResult(sessions.StatusError, fmt.Sprintf("некорректное имя '%s'", name))
	} else {
		name = util.AddPrefix(name)
		result := rc.srv.worker.Enqueue(worker.ActionCreate, name)
		rc.setResult(sessions.StatusQueued, "")
		go processResult(rc, result)
	}

	rc.redirect("/")
}
