package runner

import (
	"net/http"
	"time"
)

type TestResult struct {
	Request            *http.Request
	RequestBody        string
	Response           *http.Response
	ResponseBody       string
	RequestMatchErrors []RequestMatchError
	Err                error
	Duration           time.Duration
}

func (r *TestResult) Failed() bool {
	return r.Err != nil || len(r.RequestMatchErrors) > 0
}
