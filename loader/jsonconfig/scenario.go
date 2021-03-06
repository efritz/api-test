package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/loader/util"
)

type Scenario struct {
	Name         string          `json:"name"`
	Enabled      *bool           `json:"enabled"`
	Dependencies json.RawMessage `json:"dependencies"`
	Parallel     bool            `json:"parallel"`
	Tests        []*Test         `json:"tests"`
}

func (s *Scenario) Translate(globalRequest *GlobalRequest) (*config.Scenario, error) {
	dependencies, err := util.UnmarshalStringList(s.Dependencies)
	if err != nil {
		return nil, err
	}

	tests := []*config.Test{}
	for _, jsonTest := range s.Tests {
		test, err := jsonTest.Translate(globalRequest)
		if err != nil {
			return nil, err
		}

		tests = append(tests, test)
	}

	enabled := true
	if s.Enabled != nil {
		enabled = *s.Enabled
	}

	scenario := &config.Scenario{
		Name:         s.Name,
		Enabled:      enabled,
		Dependencies: dependencies,
		Parallel:     s.Parallel,
		Tests:        tests,
	}

	return scenario, nil
}
