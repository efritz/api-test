package runner

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/logging"
	"github.com/efritz/pentimento"
	"golang.org/x/sync/semaphore"
)

type (
	Runner struct {
		config                *config.Config
		env                   []string
		logger                logging.Logger
		verbose               bool
		junitReportPath       string
		scenarioRunnerFactory ScenarioRunnerFactory
		client                *http.Client
		names                 []string
		contexts              map[string]*ScenarioContext
		submitMutex           sync.Mutex
		ctx                   context.Context
		sequenceSemaphore     *semaphore.Weighted
		wg                    sync.WaitGroup
	}

	ScenarioContext struct {
		Scenario *config.Scenario
		Runner   ScenarioRunner
		Pending  bool
		Skipped  bool
		Context  map[string]interface{}
	}
)

func NewRunner(
	config *config.Config,
	runnerConfigs ...RunnerConfigFunc,
) *Runner {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	r := &Runner{
		config:                config,
		logger:                logging.NilLogger,
		junitReportPath:       "",
		scenarioRunnerFactory: ScenarioRunnerFactoryFunc(NewScenarioRunner),
		client:                client,
		names:                 []string{},
		contexts:              map[string]*ScenarioContext{},
		ctx:                   context.Background(),
		sequenceSemaphore:     makeScenarioSemaphore(config.Options.ForceSequential, len(config.Scenarios)),
	}

	for _, f := range runnerConfigs {
		f(r)
	}

	var scenarioLogger logging.Logger = logging.NilLogger
	if r.verbose {
		scenarioLogger = r.logger
	}

	for name := range config.Scenarios {
		r.names = append(r.names, name)
	}

	sort.Strings(r.names)

	for _, name := range r.names {
		scenario := config.Scenarios[name]

		r.contexts[name] = &ScenarioContext{
			Scenario: scenario,
			Runner: r.scenarioRunnerFactory.New(
				scenario,
				scenarioLogger,
				makeTestSemaphore(
					config.Options.ForceSequential,
					config.Options.MaxParallelism,
				),
			),
			Pending: true,
		}
	}

	return r
}

func (r *Runner) Run() error {
	r.submitReady()
	r.waitForTests()
	return r.writeReport()
}

func (r *Runner) submitReady() {
	r.submitMutex.Lock()
	defer r.submitMutex.Unlock()

	for r.submitReadyLocked() {
	}
}

func (r *Runner) submitReadyLocked() bool {
	for _, context := range r.contexts {
		if context.Scenario.Disabled || !context.Pending {
			continue
		}

		shouldSkip := false
		shouldSubmit := true

		if !context.Scenario.Enabled {
			shouldSkip = true
		}

		for _, dependency := range context.Scenario.Dependencies {
			if r.contexts[dependency].Skipped {
				shouldSkip = true
			}

			if r.contexts[dependency].Runner.Errored() || r.contexts[dependency].Runner.Failed() {
				shouldSkip = true
			}

			if !r.contexts[dependency].Runner.Resolved() {
				shouldSubmit = false
			}
		}

		if shouldSkip {
			context.Skipped = true
			context.Pending = false
			return true
		}

		if shouldSubmit {
			context.Pending = false

			r.wg.Add(1)
			go r.submit(context)
		}
	}

	return false
}

func (r *Runner) submit(context *ScenarioContext) {
	defer r.wg.Done()

	_ = r.sequenceSemaphore.Acquire(r.ctx, 1)
	defer r.sequenceSemaphore.Release(1)

	context.Context = context.Runner.Run(r.client, r.makeContext(context))
	r.submitReady()
}

func (r *Runner) makeContext(context *ScenarioContext) map[string]interface{} {
	envMap := map[string]string{}
	for _, name := range r.env {
		envMap[name] = os.Getenv(name)
	}

	newContext := map[string]interface{}{
		"env": envMap,
	}

	for _, dependency := range r.getAllDependencies(context) {
		newContext[dependency] = r.contexts[dependency].Context
	}

	return newContext
}

func (r *Runner) getAllDependencies(context *ScenarioContext) []string {
	dependencies := map[string]bool{}
	for _, dependency := range context.Scenario.Dependencies {
		dependencies[dependency] = true

		for _, dependency := range r.getAllDependencies(r.contexts[dependency]) {
			dependencies[dependency] = true
		}
	}

	flattened := []string{}
	for dependency := range dependencies {
		flattened = append(flattened, dependency)
	}

	sort.Strings(flattened)
	return flattened
}

func (r *Runner) waitForTests() {
	started := time.Now()

	if r.verbose {
		r.wg.Wait()
	} else {
		pentimento.PrintProgress(
			r.progressUpdater,
			pentimento.WithWriter(logging.Writer(r.logger)),
		)
	}

	displaySummary(r.logger, r.contexts, started)
}

func (r *Runner) progressUpdater(p *pentimento.Printer) error {
	halt := make(chan struct{})
	defer close(halt)

	go func() {
	outer:
		for {
			select {
			case <-halt:
				break outer
			case <-time.After(time.Millisecond * 100):
			}

			displayProgress(r.logger, r.names, r.contexts, p)
		}

		displayProgress(r.logger, r.names, r.contexts, p)
	}()

	r.wg.Wait()
	return nil
}

func (r *Runner) writeReport() error {
	if r.junitReportPath == "" {
		return nil
	}

	content, err := formatJUnitReport(r.contexts)
	if err != nil {
		return err
	}

	content = append([]byte(xml.Header), content...)

	if err := ioutil.WriteFile(r.junitReportPath, content, 0644); err != nil {
		return err
	}

	return nil
}

//
// Helpers

func makeScenarioSemaphore(forceSequential bool, weight int) *semaphore.Weighted {
	if forceSequential {
		return semaphore.NewWeighted(1)
	}

	return semaphore.NewWeighted(int64(weight))
}

func makeTestSemaphore(forceSequential bool, maxParallelism int) *semaphore.Weighted {
	if forceSequential {
		return semaphore.NewWeighted(1)
	}

	if maxParallelism > 0 {
		return semaphore.NewWeighted(int64(maxParallelism))
	}

	return nil
}
