package runner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/efritz/api-test/config"
	jq "github.com/efritz/go-jq"
	"github.com/xeipuuv/gojsonschema"
)

type RequestMatchError struct {
	Type     string
	Expected string
	Actual   string
}

func matchResponse(resp *http.Response, expected *config.Response) (string, map[string]interface{}, []RequestMatchError, error) {
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, nil, err
	}

	headers := map[string]string{}
	for name, values := range resp.Header {
		headers[name] = values[0]
	}

	errors := []RequestMatchError{}
	strStatus := fmt.Sprintf("%d", resp.StatusCode)

	if match, _ := matchRegex(expected.Status, strStatus); !match {
		errors = append(errors, RequestMatchError{
			Type:     "Status Code",
			Expected: fmt.Sprintf("%s", expected.Status),
			Actual:   strStatus,
		})
	}

	context := map[string]interface{}{
		"status":  resp.StatusCode,
		"headers": headers,
	}

	for key, extractor := range expected.Extract {
		sourceType, source := getSource(extractor, resp, content)
		value, matchError := getValue(extractor, resp, sourceType, source)
		if matchError != nil {
			errors = append(errors, *matchError)
			continue
		}

		matchError, err := assert(extractor.Assert, sourceType, source)
		if err != nil {
			return "", nil, nil, err
		}

		if matchError != nil {
			errors = append(errors, *matchError)
			continue
		}

		context[key] = value
	}

	return string(content), context, errors, nil
}

func getSource(
	extractor *config.ValueExtractor,
	resp *http.Response,
	content []byte,
) (string, string) {
	if extractor.Header != "" {
		sourceType := fmt.Sprintf("Header '%s'", extractor.Header)
		return sourceType, resp.Header[extractor.Header][0]
	}

	return "body", string(content)
}

func getValue(
	extractor *config.ValueExtractor,
	resp *http.Response,
	sourceType string,
	source string,
) (interface{}, *RequestMatchError) {
	if extractor.JQ != "" {
		value, err := extract(source, extractor.JQ, extractor.IsList)
		if err != nil {
			return nil, &RequestMatchError{
				Type:     sourceType,
				Expected: fmt.Sprintf("%s", extractor.JQ),
				Actual:   source,
			}
		}

		return value, nil
	}

	if extractor.Pattern != nil {
		match, value := matchRegex(extractor.Pattern, source)
		if !match {
			return nil, &RequestMatchError{
				Type:     sourceType,
				Expected: fmt.Sprintf("%s", extractor.Pattern),
				Actual:   source,
			}
		}

		// TODO - collapse if no capture groups
		return value, nil
	}

	return nil, nil
}

func assert(
	assertion *config.ValueAssertion,
	sourceType string,
	value interface{},
) (*RequestMatchError, error) {
	if assertion == nil {
		return nil, nil
	}

	if assertion.Pattern != nil {
		strValue := fmt.Sprintf("%s", value)

		if match, _ := matchRegex(assertion.Pattern, strValue); !match {
			return &RequestMatchError{
				Type:     sourceType,
				Expected: fmt.Sprintf("%s", assertion.Pattern),
				Actual:   strValue,
			}, nil
		}
	}

	if assertion.Schema != nil {
		result, err := assertion.Schema.Validate(gojsonschema.NewGoLoader(value))
		if err != nil {
			return nil, err
		}

		if !result.Valid() {
			validationErrors := []string{}
			for _, validationError := range result.Errors() {
				validationErrors = append(
					validationErrors,
					validationError.String(),
				)
			}

			// TODO - indent as well
			serialized, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}

			// TODO - use validation errors
			return &RequestMatchError{
				Type:     sourceType,
				Expected: fmt.Sprintf("%s", assertion.Pattern),
				Actual:   string(serialized),
			}, nil
		}
	}

	return nil, nil
}

func extract(content, expr string, all bool) (interface{}, error) {
	var payload interface{}
	if err := json.Unmarshal([]byte(content), &payload); err != nil {
		return nil, err
	}

	results, err := jq.Run(expr, payload)
	if err != nil {
		return nil, err
	}

	if all {
		return results, nil
	}

	if len(results) > 0 {
		return results[0], nil
	}

	return nil, nil
}

func matchRegex(re *regexp.Regexp, val string) (bool, []string) {
	if re == nil {
		return true, nil
	}

	if !re.MatchString(val) {
		return false, nil
	}

	return true, re.FindStringSubmatch(val)
}
