package config

import "time"

type Test struct {
	Name          string
	Enabled       bool
	Disabled      bool
	Request       *Request
	Response      *Response
	Retries       int
	RetryInterval time.Duration
}
