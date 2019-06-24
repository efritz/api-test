package runner

import (
	"github.com/efritz/api-test/logging"
)

type RunnerConfigFunc func(*Runner)

func WithLogger(logger logging.Logger) RunnerConfigFunc {
	return func(r *Runner) { r.logger = logger }
}

func WithVerbosityLevel(verbosityLevel logging.VerbosityLevel) RunnerConfigFunc {
	return func(r *Runner) { r.verbosityLevel = verbosityLevel }
}

func WithEnvironment(env []string) RunnerConfigFunc {
	return func(r *Runner) { r.env = env }
}

func WithJUnitReportPath(path string) RunnerConfigFunc {
	return func(r *Runner) { r.junitReportPath = path }
}

func WithScenarioRunnerFactory(factory ScenarioRunnerFactory) RunnerConfigFunc {
	return func(r *Runner) { r.scenarioRunnerFactory = factory }
}
