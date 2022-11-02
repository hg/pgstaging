package web

import (
	"github.com/hg/pgstaging/util"
	"github.com/hg/pgstaging/worker"
	"fmt"
	"net/http"
	"strings"
)

func serveModify(rc *requestContext) {
	if !rc.isMethod(http.MethodPost) {
		return
	}

	name := util.NormalizeName(rc.request.FormValue("name"))

	if len(name) < len("dev_") || !strings.HasPrefix(name, "dev_") {
		rc.setResult("error", fmt.Sprintf("некорректное имя '%s'", name))
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
		rc.setResult("error", fmt.Sprintf("неизвестное действие '%s'", action))
	}

	if result != nil {
		go processResult(rc, result)
	}

	rc.redirect("/")
}
