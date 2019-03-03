package jsonconfig

import "github.com/efritz/api-test/config"

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

	return &config.Test{
		Name:     c.Name,
		Request:  request,
		Response: response,
	}, nil
}
