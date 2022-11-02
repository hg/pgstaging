package web

func processResult(req *requestContext, result <-chan error) {
	if err := <-result; err == nil {
		req.setResult("success", "")
	} else {
		req.setResult("error", err.Error())
	}
}
