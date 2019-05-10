package jsonconfig

import (
	"github.com/efritz/api-test/config"
)

type Override struct {
	Options *Options `json:"options"`
}

func (c *Override) Translate() (*config.Override, error) {
	options, err := c.Options.Translate()
	if err != nil {
		return nil, err
	}

	return &config.Override{
		Options: options,
	}, nil
}
