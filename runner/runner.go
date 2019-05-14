package runner

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/logging"
	"github.com/efritz/pentimento"
)

type (
	Runner struct {
		config                *config.Config
		logger                logging.Logger
		junitReportPath       string
		scenarioRunnerFactory ScenarioRunnerFactory
		client                *http.Client
		names                 []string
		contexts              map[string]*ScenarioContext
		halt                  chan struct{}
		sequenceMutex         sync.Mutex
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
		halt:                  make(chan struct{}),
		names:                 []string{},
		contexts:              map[string]*ScenarioContext{},
	}

	for _, f := range runnerConfigs {
		f(r)
	}

	for name := range config.Scenarios {
		r.names = append(r.names, name)
	}

	sort.Strings(r.names)

	for _, name := range r.names {
		scenario := config.Scenarios[name]

		r.contexts[name] = &ScenarioContext{
			Scenario: scenario,
			Runner:   r.scenarioRunnerFactory.New(scenario, config.Options.ForceSequential),
			Pending:  true,
		}
	}

	return r
}

func (r *Runner) Run() error {
	started := time.Now()
	r.submitReady()

	go func() {
		defer close(r.halt)
		r.wg.Wait()
	}()

	pentimento.PrintProgress(
		r.progressUpdater,
		pentimento.WithWriter(logging.Writer(r.logger)),
	)

	displaySummary(r.logger, r.contexts, started)

	if r.junitReportPath != "" {
		content, err := formatJUnitReport(r.contexts)
		if err != nil {
			return err
		}

		content = append([]byte(xml.Header), content...)

		if err := ioutil.WriteFile(r.junitReportPath, content, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (r *Runner) progressUpdater(p *pentimento.Printer) error {
outer:
	for {
		select {
		case <-r.halt:
			break outer
		case <-time.After(time.Millisecond * 100):
		}

		displayProgress(r.logger, r.names, r.contexts, p)
	}

	displayProgress(r.logger, r.names, r.contexts, p)
	return nil
}

func (r *Runner) submitReady() {
	skipped := false

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
			skipped = true
			context.Skipped = true
			context.Pending = false
			continue
		}

		if shouldSubmit {
			context.Pending = false

			r.wg.Add(1)
			go r.submit(context)
		}
	}

	if skipped {
		r.submitReady()
	}
}

func (r *Runner) submit(context *ScenarioContext) {
	defer r.wg.Done()

	if r.config.Options.ForceSequential {
		r.sequenceMutex.Lock()
		defer r.sequenceMutex.Unlock()
	}

	context.Context = context.Runner.Run(r.client, r.makeContext(context))
	r.submitReady()
}

func (r *Runner) makeContext(context *ScenarioContext) map[string]interface{} {
	// TODO - make a whitelist instead
	envMap := map[string]string{}
	for _, pair := range os.Environ() {
		parts := strings.Split(pair, "=")
		envMap[parts[0]] = parts[1]
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
