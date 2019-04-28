package jsonconfig

import (
	"fmt"
	"strings"

	"github.com/efritz/api-test/config"
)

type Test struct {
	Name     string    `json:"name"`
	Request  *Request  `json:"request"`
	Response *Response `json:"response"`
}

func (c *Test) Translate(globalRequest *GlobalRequest) (*config.Test, error) {
	if c.Request == nil {
		c.Request = &Request{}
	}

	request, err := c.Request.Translate(globalRequest)
	if err != nil {
		return nil, err
	}

	if c.Response == nil {
		c.Response = &Response{}
	}

	response, err := c.Response.Translate()
	if err != nil {
		return nil, err
	}

	name := c.Name
	if name == "" {
		name = fmt.Sprintf(
			"%s %s",
			strings.ToUpper(request.Method),
			c.Request.URI,
		)
	}

	return &config.Test{
		Name:     name,
		Request:  request,
		Response: response,
	}, nil
}
