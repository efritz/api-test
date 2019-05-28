package jsonconfig

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/efritz/api-test/config"
)

type Response struct {
	Status  json.RawMessage            `json:"status"`
	Extract map[string]json.RawMessage `json:"extract"`
}

var patternOK = regexp.MustCompile("2..")

func (r *Response) Translate() (*config.Response, error) {
	status, err := unmarshalStatus(r.Status)
	if err != nil {
		return nil, err
	}

	extractors := map[string]*config.ValueExtractor{}
	for key, value := range r.Extract {
		extractor, err := unmarshalValueExtractor(value)
		if err != nil {
			return nil, err
		}

		extractors[key] = extractor
	}

	return &config.Response{
		Status:  status,
		Extract: extractors,
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
