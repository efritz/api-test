package jsonconfig

import (
	"fmt"
	"strings"

	"github.com/efritz/api-test/config"
)

type (
	Request struct {
		URI     string            `json:"uri"`
		Method  string            `json:"method"`
		Body    string            `json:"body"`
		Headers map[string]string `json:"headers"`
		Auth    *BasicAuth        `json:"auth"`
	}
)

func (c *Request) Translate(globalRequest *GlobalRequest) (*config.Request, error) {
	method := c.Method
	if method == "" {
		method = "get"
	}

	url := c.URI
	headers := c.Headers
	jsonAuth := c.Auth

	if headers == nil {
		headers = map[string]string{}
	}

	if globalRequest != nil {
		if globalRequest.BaseURL != "" && isRelative(url) {
			url = fmt.Sprintf(
				"%s/%s",
				strings.TrimRight(globalRequest.BaseURL, "/"),
				strings.TrimLeft(url, "/"),
			)

			for key, val := range globalRequest.Headers {
				if _, ok := headers[key]; !ok {
					headers[key] = val
				}
			}

			if jsonAuth == nil {
				jsonAuth = globalRequest.Auth
			}
		}
	}

	auth, err := jsonAuth.Translate()
	if err != nil {
		return nil, err
	}

	return &config.Request{
		URI:     url,
		Method:  method,
		Headers: headers,
		Auth:    auth,
		Body:    c.Body,
	}, nil
}

func isRelative(uri string) bool {
	for _, prefix := range []string{"http://", "https://"} {
		if strings.HasPrefix(uri, prefix) {
			return false
		}
	}

	return true
}
