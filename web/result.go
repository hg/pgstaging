package web

import "github.com/hg/pgstaging/web/sessions"

func processResult(req *requestContext, result <-chan error) {
	if err := <-result; err == nil {
		req.setResult(sessions.StatusSuccess, "")
	} else {
		req.setResult(sessions.StatusError, err.Error())
	}
}
