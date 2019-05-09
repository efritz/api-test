package jsonconfig

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/loader/util"
)

type Response struct {
	Status     json.RawMessage            `json:"status"`
	Headers    map[string]json.RawMessage `json:"headers"`
	Body       string                     `json:"body"`
	Extract    string                     `json:"extract"`
}

var patternOK = regexp.MustCompile("2..")

func (r *Response) Translate() (*config.Response, error) {
	status, err := unmarshalStatus(r.Status)
	if err != nil {
		return nil, err
	}

	headers := map[string][]*regexp.Regexp{}
	for name, raw := range r.Headers {
		values, err := util.UnmarshalStringList(raw)
		if err != nil {
			return nil, err
		}

		patterns := []*regexp.Regexp{}
		for _, value := range values {
			pattern, err := regexp.Compile(value)
			if err != nil {
				return nil, fmt.Errorf("illegal header regex")
			}

			patterns = append(patterns, pattern)
		}

		headers[name] = patterns
	}

	body, err := regexp.Compile(r.Body)
	if err != nil {
		return nil, fmt.Errorf("illegal body regex")
	}

	return &config.Response{
		Status:     status,
		Headers:    headers,
		Body:       body,
		Extract:    r.Extract,
	}, nil
}

func unmarshalStatus(data json.RawMessage) (*regexp.Regexp, error) {
	if len(data) == 0 {
		return patternOK, nil
	}

	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		return regexp.Compile(fmt.Sprintf("%d", num))
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		pattern, err := regexp.Compile(str)
		if err != nil {
			return nil, fmt.Errorf("illegal status regex")
		}

		return pattern, nil
	}

	return nil, fmt.Errorf("status value is neither string nor int")
}
