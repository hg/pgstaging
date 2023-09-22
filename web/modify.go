package web

import (
	"fmt"
	"github.com/hg/pgstaging/web/util"
	"github.com/hg/pgstaging/worker"
	"net/http"
)

func parseAction(name string) (worker.Action, error) {
	switch name {
	case "start":
		return worker.ActionStart, nil
	case "stop":
		return worker.ActionStop, nil
	case "drop":
		return worker.ActionDrop, nil
	default:
		return 0, fmt.Errorf("invalid action %s", name)
	}
}

func serveModify(rc *requestContext) {
	if !rc.requireMethod(http.MethodPost) {
		return
	}

	name := util.NormalizeName(rc.request.PostFormValue("name"))

	if !util.IsDevName(name) {
		rc.bail(fmt.Sprintf("некорректное имя '%s'", name))
		return
	}

	pass := rc.request.PostFormValue("password")

	action, err := parseAction(rc.request.FormValue("action"))

	if err == nil {
		result := rc.queueTask(action, name, pass)
		rc.redirect("/")
		go rc.processResult(result)
	} else {
		rc.bail(err.Error())
	}
}
