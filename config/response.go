package config

import "regexp"

type Response struct {
	Status     *regexp.Regexp
	Headers    map[string][]*regexp.Regexp
	Body       *regexp.Regexp
	Extract    string
	Assertions map[string]*regexp.Regexp
	// TODO - json schema body
}
