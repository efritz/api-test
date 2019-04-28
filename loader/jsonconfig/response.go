package jsonconfig

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/loader/util"
)

type Response struct {
	Status  *json.RawMessage           `json:"status"`
	Headers map[string]json.RawMessage `json:"headers"`
	Body    string                     `json:"body"`
	Extract string                     `json:"extract"`
}

var patternOK = regexp.MustCompile("2..")

func (c *Response) Translate() (*config.Response, error) {
	status, err := unmarshalStatus(c.Status)
	if err != nil {
		return nil, err
	}

	headers := map[string][]*regexp.Regexp{}
	for name, raw := range c.Headers {
		values, err := util.UnmarshalStringList(raw)
		if err != nil {
			return nil, err
		}

		patterns := []*regexp.Regexp{}
		for _, value := range values {
			pattern, err := regexp.Compile(value)
			if err != nil {
				return nil, err
			}

			patterns = append(patterns, pattern)
		}

		headers[name] = patterns
	}

	body, err := regexp.Compile(c.Body)
	if err != nil {
		return nil, err
	}

	return &config.Response{
		Status:  status,
		Headers: headers,
		Body:    body,
		Extract: c.Extract,
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
