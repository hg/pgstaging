package web

import (
	"fmt"
	"github.com/hg/pgstaging/web/util"
	"github.com/hg/pgstaging/worker"
	"net/http"
	"strconv"
)

func serveCreatePullRequest(rc *requestContext) {
	if !rc.requireMethod(http.MethodPost) {
		return
	}

	prRaw := rc.request.PostFormValue("pr")
	pr, err := strconv.Atoi(prRaw)

	if err != nil || pr <= 0 {
		http.Error(rc.writer, "bad pull request number", http.StatusBadRequest)
		return
	}

	name := util.AddPrefix(fmt.Sprintf("pr_%d", pr))
	result := rc.queueTask(worker.ActionForceCreate, name, "")

	if re := <-result; re.Err == nil {
		fmt.Fprint(rc.writer, re.Data) // port
	} else {
		http.Error(rc.writer, err.Error(), http.StatusInternalServerError)
	}
}

func serveCreate(rc *requestContext) {
	if !rc.requireMethod(http.MethodPost) {
		return
	}

	pass := rc.request.PostFormValue("password")

	if pass != "" && !util.IsOkPassword(pass) {
		rc.bail(fmt.Sprintf("invalid password '%s'", pass))
		return
	}

	name := rc.request.PostFormValue("name")
	name = util.NormalizeName(name)

	if name == "" || len(name) > 32 {
		rc.bail(fmt.Sprintf("invalid name '%s'", name))
		return
	}

	name = util.AddPrefix(name)
	result := rc.queueTask(worker.ActionCreate, name, pass)
	rc.redirect("/")
	go rc.processResult(result)
}
