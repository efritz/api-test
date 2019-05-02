package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/api-test/config"
)

type (
	BaseConfig struct {
		Scenarios []*Scenario     `json:"scenarios"`
		Includes  json.RawMessage `json:"include"`
	}

	MainConfig struct {
		*BaseConfig
		Options       *Options       `json:"options"`
		GlobalRequest *GlobalRequest `json:"global-request"`
	}

	Options struct {
		ForceSequential bool
	}

	GlobalRequest struct {
		BaseURL string                     `json:"base-url"`
		Auth    *BasicAuth                 `json:"auth"`
		Headers map[string]json.RawMessage `json:"headers"`
	}
)

func (c *BaseConfig) Translate(globalRequest *GlobalRequest) ([]*config.Scenario, error) {
	scenarios := []*config.Scenario{}
	for _, jsonScenario := range c.Scenarios {
		scenario, err := jsonScenario.Translate(globalRequest)
		if err != nil {
			return nil, err
		}

		scenarios = append(scenarios, scenario)
	}

	return scenarios, nil
}

func (o *Options) Translate() (*config.Options, error) {
	return &config.Options{
		ForceSequential: o.ForceSequential,
	}, nil
}
