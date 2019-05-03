package jsonconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	tmpl "text/template"

	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/loader/util"
	"github.com/google/uuid"
)

type (
	Request struct {
		URI      string                     `json:"uri"`
		Method   string                     `json:"method"`
		Auth     *BasicAuth                 `json:"auth"`
		Headers  map[string]json.RawMessage `json:"headers"`
		Body     string                     `json:"body"`
		JSONBody json.RawMessage            `json:"json-body"`
		// TODO - templated JSON body
		// TODO - form
		// TODO - file
	}
)

func (c *Request) Translate(globalRequest *GlobalRequest) (*config.Request, error) {
	method := sanitizeMethod(c.Method)
	url := sanitizeURL(c.URI, globalRequest)
	jsonAuth := sanitizeAuth(c.Auth, globalRequest)
	headers, err := sanitizeHeaders(c.Headers, globalRequest)
	if err != nil {
		return nil, err
	}

	urlTemplate, err := compile(url)
	if err != nil {
		return nil, fmt.Errorf("illegal uri template (%s)", err.Error())
	}

	headerTemplates := map[string][]*tmpl.Template{}
	for name, values := range headers {
		templates := []*tmpl.Template{}
		for _, value := range values {
			template, err := compile(value)
			if err != nil {
				return nil, fmt.Errorf("illegal header template (%s)", err.Error())
			}

			templates = append(templates, template)
		}

		headerTemplates[name] = templates
	}

	auth, err := jsonAuth.Translate()
	if err != nil {
		return nil, err
	}

	if c.Body != "" && c.JSONBody != nil {
		return nil, fmt.Errorf("multiple bodies supplied")
	}

	var bodyTemplate *tmpl.Template

	if c.Body != "" {
		bodyTemplate, err = compile(c.Body)
		if err != nil {
			return nil, fmt.Errorf("illegal body template (%s)", err.Error())
		}
	}

	if c.JSONBody != nil {
		bodyTemplate, err = compile(string(c.JSONBody))
		if err != nil {
			return nil, fmt.Errorf("illegal json body template (%s)", err.Error())
		}
	}

	return &config.Request{
		URL:     urlTemplate,
		Method:  method,
		Headers: headerTemplates,
		Auth:    auth,
		Body:    bodyTemplate,
	}, nil
}

func sanitizeMethod(method string) string {
	if method == "" {
		return "get"
	}

	return method
}

func sanitizeURL(uri string, globalRequest *GlobalRequest) string {
	if globalRequest == nil || globalRequest.BaseURL == "" || !isRelative(uri) {
		return uri
	}

	return fmt.Sprintf("%s/%s", strings.TrimRight(globalRequest.BaseURL, "/"), strings.TrimLeft(uri, "/"))
}

func sanitizeHeaders(rawHeaders map[string]json.RawMessage, globalRequest *GlobalRequest) (map[string][]string, error) {
	headers := map[string][]string{}
	for name, raw := range rawHeaders {
		values, err := util.UnmarshalStringList(raw)
		if err != nil {
			return nil, err
		}

		headers[name] = values
	}

	if globalRequest != nil {
		for name, raw := range globalRequest.Headers {
			values, err := util.UnmarshalStringList(raw)
			if err != nil {
				return nil, err
			}

			if _, ok := headers[name]; !ok {
				headers[name] = values
			}
		}
	}

	return headers, nil
}

func sanitizeAuth(auth *BasicAuth, globalRequest *GlobalRequest) *BasicAuth {
	if auth == nil && globalRequest != nil {
		return globalRequest.Auth
	}

	return auth
}

func isRelative(uri string) bool {
	for _, prefix := range []string{"http://", "https://"} {
		if strings.HasPrefix(uri, prefix) {
			return false
		}
	}

	return true
}

func compile(template string) (*tmpl.Template, error) {
	funcs := tmpl.FuncMap{
		"uuid": func() string { return uuid.New().String() },
		"file": func(path string) string {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				// TODO - ?
				return "<failed to read file>"
			}

			return string(content)
		},
	}

	return tmpl.New("").Funcs(funcs).Parse(template)
}
