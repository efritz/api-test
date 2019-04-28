package runner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/efritz/api-test/config"
	"github.com/efritz/go-jq"
)

type MatchError struct {
	Type     string
	Expected string
	Actual   string
}

func matchResponse(resp *http.Response, expected *config.Response) (*http.Response, string, map[string]interface{}, []MatchError, error) {
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", nil, nil, err
	}

	errors := []MatchError{}

	match, statusGroups := matchRegex(expected.Status, fmt.Sprintf("%d", resp.StatusCode))
	if !match {
		errors = append(errors, MatchError{
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
				errors = append(errors, MatchError{
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
		errors = append(errors, MatchError{
			Type:     "body",
			Expected: fmt.Sprintf("%s", expected.Body),
			Actual:   "<placeholder>", // string(content),
		})
	}

	context := map[string]interface{}{
		// TODO - rename these (keep it symmetric in derision as well)
		"statusGroups": statusGroups,
		"headerGroups": headerGroups,
		"bodyGroups":   bodyGroups,
	}

	if expected.Extract != "" {
		var payload interface{}
		if err := json.Unmarshal(content, &payload); err != nil {
			return nil, "", nil, nil, err
		}

		results, err := jq.Run(expected.Extract, payload)
		if err != nil {
			return nil, "", nil, nil, err
		}

		if len(results) != 1 {
			return nil, "", nil, nil, fmt.Errorf("extraction expects a single object")
		}

		resultMap, ok := results[0].(map[string]interface{})
		if !ok {
			return nil, "", nil, nil, fmt.Errorf("extraction expects a single object")
		}

		for k, v := range resultMap {
			context[k] = v
		}
	}

	return resp, string(content), context, errors, nil
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
