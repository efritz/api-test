package runner

import (
	"net/http"
	"time"
)

type TestResult struct {
	Index              int
	Disabled           bool
	Skipped            bool
	Request            *http.Request
	RequestBody        string
	Response           *http.Response
	ResponseBody       string
	RequestMatchErrors []RequestMatchError
	Err                error
	Duration           time.Duration
}

func (r *TestResult) Errored() bool {
	return r.Err != nil
}

func (r *TestResult) Failed() bool {
	return len(r.RequestMatchErrors) > 0
}
