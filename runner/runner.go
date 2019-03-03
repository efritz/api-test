package runner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/logging"
)

type Runner struct {
	config *config.Config
	logger logging.Logger
}

func NewRunner(config *config.Config, logger logging.Logger) *Runner {
	return &Runner{
		config: config,
		logger: logger,
	}
}

func (r *Runner) Run() error {
	for _, test := range r.config.Tests {
		r.logger.Info("Starting test %s", test.Name)

		req, err := buildRequest(test.Request)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if err := checkResponse(resp, test.Response); err != nil {
			r.logger.Error("Failure: %s", err.Error())
		}
	}

	return nil
}

func buildRequest(prototype *config.Request) (*http.Request, error) {
	req, err := http.NewRequest(
		strings.ToUpper(prototype.Method),
		prototype.URI,
		bytes.NewReader([]byte(prototype.Body)),
	)

	if err != nil {
		return nil, err
	}

	return req, err
}

func checkResponse(resp *http.Response, expected *config.Response) error {
	defer resp.Body.Close()

	if !expected.Status.MatchString(fmt.Sprintf("%d", resp.StatusCode)) {
		return fmt.Errorf("status code mismatch '%d' != '%s'", resp.StatusCode, expected.Status)
	}

	for key, pattern := range expected.Headers {
		value := resp.Header.Get(key)

		if !pattern.MatchString(value) {
			return fmt.Errorf("header value mismatch '%s' != '%s'", value, pattern)
		}
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !expected.Body.Match(content) {
		return fmt.Errorf("body mismatch '%s' != '%s'", content, expected.Body)
	}

	return nil
}
