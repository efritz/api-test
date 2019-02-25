package loader

import (
	"github.com/efritz/api-test/config"
)

type jsonConfig struct {
	Tests []string `json:"tests"`
}

func (c *jsonConfig) Translate() (*config.Config, error) {
	config := &config.Config{
		Tests: c.Tests,
	}

	return config, nil
}
