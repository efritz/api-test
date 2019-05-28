package runner

import (
	"github.com/efritz/api-test/logging"
)

type RunnerConfigFunc func(*Runner)

func WithLogger(logger logging.Logger) RunnerConfigFunc {
	return func(c *Runner) { c.logger = logger }
}

func WithEnvironment(env []string) RunnerConfigFunc {
	return func(c *Runner) { c.env = env }
}

func WithJUnitReportPath(path string) RunnerConfigFunc {
	return func(c *Runner) { c.junitReportPath = path }
}

func WithScenarioRunnerFactory(factory ScenarioRunnerFactory) RunnerConfigFunc {
	return func(c *Runner) { c.scenarioRunnerFactory = factory }
}
