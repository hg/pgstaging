package web

import (
	"github.com/hg/pgstaging/web/util"
	"github.com/hg/pgstaging/worker"
	"fmt"
	"net/http"
)

func serveCreate(rc *requestContext) {
	if !rc.requireMethod(http.MethodPost) {
		return
	}

	pass := rc.request.PostFormValue("password")

	if pass != "" && !util.IsOkPassword(pass) {
		rc.bail(fmt.Sprintf("некорректный пароль '%s'", pass))
		return
	}

	name := rc.request.PostFormValue("name")
	name = util.NormalizeName(name)

	if name == "" || len(name) > 32 {
		rc.bail(fmt.Sprintf("некорректное имя '%s'", name))
		return
	}

	name = util.AddPrefix(name)
	rc.queueTask(worker.ActionCreate, name, pass)
}
