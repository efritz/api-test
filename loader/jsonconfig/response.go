package jsonconfig

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/efritz/api-test/config"
)

type Response struct {
	Status  *json.RawMessage  `json:"status"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}

var patternOK = regexp.MustCompile("2..")

func (c *Response) Translate() (*config.Response, error) {
	status, err := unmarshalStatus(c.Status)
	if err != nil {
		return nil, err
	}

	headers := map[string]*regexp.Regexp{}
	for key, value := range c.Headers {
		pattern, err := regexp.Compile(value)
		if err != nil {
			return nil, err
		}

		headers[key] = pattern
	}

	body, err := regexp.Compile(c.Body)
	if err != nil {
		return nil, err
	}

	return &config.Response{
		Status:  status,
		Headers: headers,
		Body:    body,
	}, nil
}

func unmarshalStatus(data *json.RawMessage) (*regexp.Regexp, error) {
	if data == nil {
		return patternOK, nil
	}

	var num int
	if err := json.Unmarshal(*data, &num); err == nil {
		return regexp.Compile(fmt.Sprintf("%d", num))
	}

	var str string
	if err := json.Unmarshal(*data, &str); err == nil {
		return regexp.Compile(str)
	}

	return nil, fmt.Errorf("status value is neither string nor int")
}
