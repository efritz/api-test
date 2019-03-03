package jsonconfig

import "github.com/efritz/api-test/config"

type (
	Config struct {
		GlobalRequest *GlobalRequest `json:"global-request"`
		Tests         []*Test        `json:"tests"`
	}

	GlobalRequest struct {
		BaseURL string            `json:"base-url"`
		Auth    *BasicAuth        `json:"auth"`
		Headers map[string]string `json:"headers"`
	}
)

func (c *Config) Translate() (*config.Config, error) {
	tests := []*config.Test{}
	for _, jsonTest := range c.Tests {
		test, err := jsonTest.Translate(c.GlobalRequest)
		if err != nil {
			return nil, err
		}

		tests = append(tests, test)
	}

	config := &config.Config{
		Tests: tests,
	}

	return config, nil
}
