package config

import "regexp"

type Response struct {
	Status     *regexp.Regexp
	Headers    map[string][]*regexp.Regexp
	Body       *regexp.Regexp
	Extract    string
	// TODO - json schema body
}
