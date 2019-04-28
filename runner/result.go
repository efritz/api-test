package runner

import (
	"net/http"
	"time"
)

type TestResult struct {
	Request      *http.Request
	RequestBody  string
	Response     *http.Response
	ResponseBody string
	MatchErrors  []MatchError
	Err          error
	Duration     time.Duration
}

func (r *TestResult) Failed() bool {
	return r.Err != nil || len(r.MatchErrors) > 0
}
