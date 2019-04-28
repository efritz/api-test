package runner

import (
	"bytes"
	"net/http"
	"strings"
	tmpl "text/template"

	"github.com/efritz/api-test/config"
)

func buildRequest(prototype *config.Request, context map[string]interface{}) (*http.Request, string, error) {
	url, err := applyTemplate(prototype.URL, context)
	if err != nil {
		return nil, "", err
	}

	body, err := applyTemplate(prototype.Body, context)
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequest(
		strings.ToUpper(prototype.Method),
		url,
		bytes.NewReader([]byte(body)),
	)

	if err != nil {
		return nil, "", err
	}

	for key, values := range prototype.Headers {
		for _, template := range values {
			value, err := applyTemplate(template, context)
			if err != nil {
				return nil, "", err
			}

			req.Header.Add(key, value)
		}
	}

	return req, body, err
}

func applyTemplate(t *tmpl.Template, args map[string]interface{}) (string, error) {
	if t == nil {
		return "", nil
	}

	buffer := &bytes.Buffer{}
	if err := t.Execute(buffer, args); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
