package config

import (
	tmpl "text/template"
)

type (
	Request struct {
		URL     *tmpl.Template
		Method  string
		Auth    *BasicAuth
		Headers map[string][]*tmpl.Template
		Body    *tmpl.Template
	}

	BasicAuth struct {
		Username string
		Password string
	}
)
