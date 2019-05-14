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
		// TODO - form
		// TODO - file
	}

	BasicAuth struct {
		Username *tmpl.Template
		Password *tmpl.Template
	}
)
