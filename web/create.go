package web

import (
	"github.com/hg/pgstaging/util"
	"github.com/hg/pgstaging/worker"
	"fmt"
	"net/http"
)

func serveCreate(rc *requestContext) {
	if !rc.isMethod(http.MethodPost) {
		return
	}

	name := util.NormalizeName(rc.request.PostFormValue("name"))

	if name == "" || len(name) > 32 {
		rc.setResult("error", fmt.Sprintf("некорректное имя '%s'", name))
		rc.redirect("/")
		return
	}

	result := rc.srv.worker.Enqueue(worker.ActionCreate, name)
	go processResult(rc, result)

	rc.setResult("queued", "")
	rc.redirect("/")
}
