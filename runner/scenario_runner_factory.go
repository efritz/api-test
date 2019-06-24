package runner

import (
	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/logging"
	"golang.org/x/sync/semaphore"
)

type (
	ScenarioRunnerFactory interface {
		New(scenario *config.Scenario, logger logging.Logger, testSemaphore *semaphore.Weighted) ScenarioRunner
	}

	ScenarioRunnerFactoryFunc func(scenario *config.Scenario, logger logging.Logger, testSemaphore *semaphore.Weighted) ScenarioRunner
)

func (f ScenarioRunnerFactoryFunc) New(scenario *config.Scenario, logger logging.Logger, testSemaphore *semaphore.Weighted) ScenarioRunner {
	return f(scenario, logger, testSemaphore)
}
