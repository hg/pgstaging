package web

import (
	"github.com/hg/pgstaging/web/util"
	"github.com/hg/pgstaging/worker"
	"fmt"
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
		rc.queueTask(action, name, pass)
	} else {
		rc.bail(err.Error())
	}
}
