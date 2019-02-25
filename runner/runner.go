package runner

import (
	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/logging"
)

type Runner struct {
	config *config.Config
	logger logging.Logger
}

func NewRunner(config *config.Config, logger logging.Logger) *Runner {
	return &Runner{
		config: config,
		logger: logger,
	}
}

func (r *Runner) Run() error {
	// TODO
	r.logger.Info("running tests with config: %#v", r.config)
	return nil
}
