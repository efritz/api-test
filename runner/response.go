package runner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/efritz/api-test/config"
	jq "github.com/efritz/go-jq"
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

	errors := []RequestMatchError{}

	match, statusGroups := matchRegex(expected.Status, fmt.Sprintf("%d", resp.StatusCode))
	if !match {
		errors = append(errors, RequestMatchError{
			Type:     "Status Code",
			Expected: fmt.Sprintf("%s", expected.Status),
			Actual:   fmt.Sprintf("%d", resp.StatusCode),
		})
	}

	headerGroups := map[string][]string{}
	for key, patterns := range expected.Headers {
		// TODO - not sure how to order these
		for _, pattern := range patterns {
			value := resp.Header.Get(key)

			match, groups := matchRegex(pattern, value)
			if !match {
				errors = append(errors, RequestMatchError{
					Type:     fmt.Sprintf("Header '%s'", key),
					Expected: fmt.Sprintf("%s", pattern),
					Actual:   value,
				})
			}

			headerGroups[key] = groups
		}
	}

	match, bodyGroups := matchRegex(expected.Body, string(content))
	if !match {
		errors = append(errors, RequestMatchError{
			Type:     "Body",
			Expected: fmt.Sprintf("%s", expected.Body),
			Actual:   string(content),
		})
	}

	context := map[string]interface{}{
		// TODO - rename these (keep it symmetric in derision as well)
		"statusGroups": statusGroups,
		"headerGroups": headerGroups,
		"bodyGroups":   bodyGroups,
	}

	for key, expr := range expected.Extract {
		value, err := extract(content, expr, false)
		if err != nil {
			return "", nil, nil, err
		}

		context[key] = value
	}

	for key, expr := range expected.ExtractList {
		value, err := extract(content, expr, true)
		if err != nil {
			return "", nil, nil, err
		}

		context[key] = value
	}

	extractionGroups := map[string][]string{}
	for key, pattern := range expected.Assertions {
		rawValue := context[key]
		strValue := fmt.Sprintf("%v", rawValue)

		match, groups := matchRegex(pattern, strValue)
		if !match {
			errors = append(errors, RequestMatchError{
				Type:     key,
				Expected: fmt.Sprintf("%s", pattern),
				Actual:   strValue,
			})
		}

		extractionGroups[key] = groups
	}

	context["extractionGroups"] = extractionGroups

	return string(content), context, errors, nil
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

func extract(content []byte, expr string, all bool) (interface{}, error) {
	var payload interface{}
	if err := json.Unmarshal(content, &payload); err != nil {
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
