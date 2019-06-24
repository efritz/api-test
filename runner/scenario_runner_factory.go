package runner

import (
	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/logging"
)

type (
	ScenarioRunnerFactory interface {
		New(scenario *config.Scenario, logger logging.Logger, forceSequential bool) ScenarioRunner
	}

	ScenarioRunnerFactoryFunc func(scenario *config.Scenario, logger logging.Logger, forceSequential bool) ScenarioRunner
)

func (f ScenarioRunnerFactoryFunc) New(scenario *config.Scenario, logger logging.Logger, forceSequential bool) ScenarioRunner {
	return f(scenario, logger, forceSequential)
}
