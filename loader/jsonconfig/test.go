package jsonconfig

import (
	"fmt"
	"strings"
	"time"

	"github.com/efritz/api-test/config"
)

type Test struct {
	Name          string    `json:"name"`
	Enabled       *bool     `json:"enabled"`
	Request       *Request  `json:"request"`
	Response      *Response `json:"response"`
	Retries       int       `json:"retries"`
	RetryInterval float64   `json:"retry-interval"`
}

func (t *Test) Translate(globalRequest *GlobalRequest) (*config.Test, error) {
	if t.Request == nil {
		t.Request = &Request{}
	}

	request, err := t.Request.Translate(globalRequest)
	if err != nil {
		return nil, err
	}

	if t.Response == nil {
		t.Response = &Response{}
	}

	response, err := t.Response.Translate()
	if err != nil {
		return nil, err
	}

	name := t.Name
	if name == "" {
		name = fmt.Sprintf(
			"%s %s",
			strings.ToUpper(request.Method),
			t.Request.URI,
		)
	}

	enabled := true
	if t.Enabled != nil {
		enabled = *t.Enabled
	}

	return &config.Test{
		Name:          name,
		Enabled:       enabled,
		Request:       request,
		Response:      response,
		Retries:       t.Retries,
		RetryInterval: getRetryInterval(t.RetryInterval),
	}, nil
}

func getRetryInterval(rawRetryInterval float64) time.Duration {
	if rawRetryInterval == 0 {
		return time.Second
	}

	return time.Duration(float64(time.Second) * rawRetryInterval)
}
