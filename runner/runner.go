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

type Runner struct {
	config          *config.Config
	logger          logging.Logger
	junitReportPath string
	client          *http.Client
	names           []string
	contexts        map[string]*ScenarioContext
	halt            chan struct{}
	sequenceMutex   sync.Mutex
	wg              sync.WaitGroup
}

func NewRunner(config *config.Config, logger logging.Logger, junitReportPath string) *Runner {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	names := []string{}
	contexts := map[string]*ScenarioContext{}

	for name, scenario := range config.Scenarios {
		names = append(names, name)
		contexts[name] = NewScenarioContext(scenario)
	}

	sort.Strings(names)

	return &Runner{
		config:          config,
		logger:          logger,
		junitReportPath: junitReportPath,
		client:          client,
		names:           names,
		contexts:        contexts,
		halt:            make(chan struct{}),
	}
}

func (r *Runner) Run() error {
	started := time.Now()
	r.submitReady()

	go func() {
		defer close(r.halt)
		r.wg.Wait()
	}()

	pentimento.PrintProgress(func(p *pentimento.Printer) error {
	outer:
		for {
			select {
			case <-r.halt:
				break outer
			case <-time.After(time.Millisecond * 100):
			}

			displayProgress(r.names, r.contexts, p)
		}

		displayProgress(r.names, r.contexts, p)
		return nil
	})

	displaySummary(r.contexts, started, r.logger)

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

func (r *Runner) submitReady() {
	skipped := false

	for _, context := range r.contexts {
		if !context.Pending {
			continue
		}

		shouldSkip := false
		shouldSubmit := true

		for _, dependency := range context.Scenario.Dependencies {
			if r.contexts[dependency].Skipped || r.contexts[dependency].Failed() {
				shouldSkip = true
			}

			if !r.contexts[dependency].Resolved() {
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

	context.Running = true
	r.prepareContext(context)

	for result := range context.Run(r.client) {
		context.Results = append(context.Results, result)
	}

	context.Running = false
	r.submitReady()
}

func (r *Runner) prepareContext(context *ScenarioContext) {
	context.Context = map[string]interface{}{}

	envMap := map[string]string{}
	for _, pair := range os.Environ() {
		parts := strings.Split(pair, "=")
		envMap[parts[0]] = parts[1]
	}

	context.Context["env"] = envMap

	for _, dependency := range r.getAllDependencies(context) {
		context.Context[dependency] = r.contexts[dependency].Context
	}
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
