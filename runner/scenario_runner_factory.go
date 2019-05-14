package runner

import "github.com/efritz/api-test/config"

type (
	ScenarioRunnerFactory interface {
		New(scenario *config.Scenario, forceSequential bool) ScenarioRunner
	}

	ScenarioRunnerFactoryFunc func(scenario *config.Scenario, forceSequential bool) ScenarioRunner
)

func (f ScenarioRunnerFactoryFunc) New(scenario *config.Scenario, forceSequential bool) ScenarioRunner {
	return f(scenario, forceSequential)
}
