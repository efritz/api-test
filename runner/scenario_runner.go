package runner

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/logging"
	"golang.org/x/sync/semaphore"
)

type (
	ScenarioRunner interface {
		Run(client *http.Client, context map[string]interface{}) map[string]interface{}
		Results() []*TestResult
		Errored() bool
		Failed() bool
		Resolved() bool
		Duration() time.Duration
	}

	scenarioRunner struct {
		scenario       *config.Scenario
		logger         logging.Logger
		verbosityLevel logging.VerbosityLevel
		results        []*TestResult
		running        bool
		submitted      bool
		ctx            context.Context
		testSemaphore  *semaphore.Weighted
		mutex          sync.RWMutex
		waitGroup      sync.WaitGroup
	}
)

func NewScenarioRunner(
	scenario *config.Scenario,
	logger logging.Logger,
	verbosityLevel logging.VerbosityLevel,
	testSemaphore *semaphore.Weighted,
) ScenarioRunner {
	return newScenarioRunner(scenario, logger, verbosityLevel, testSemaphore)
}

func newScenarioRunner(
	scenario *config.Scenario,
	logger logging.Logger,
	verbosityLevel logging.VerbosityLevel,
	testSemaphore *semaphore.Weighted,
) *scenarioRunner {
	return &scenarioRunner{
		scenario:       scenario,
		logger:         logger,
		verbosityLevel: verbosityLevel,
		results:        []*TestResult{},
		ctx:            context.Background(),
		testSemaphore:  testSemaphore,
	}
}

func (r *scenarioRunner) Run(client *http.Client, context map[string]interface{}) map[string]interface{} {
	r.running = true
	r.submitted = true

	for result := range r.run(client, context) {
		if result.Disabled {
			continue
		}

		for len(r.results) <= result.Index {
			r.results = append(r.results, nil)
		}

		r.results[result.Index] = result
	}

	r.running = false
	return context
}

func (r *scenarioRunner) run(client *http.Client, context map[string]interface{}) <-chan *TestResult {
	ch := make(chan *TestResult)

	go func() {
		defer func() {
			r.waitGroup.Wait()
			close(ch)
		}()

		for i, test := range r.scenario.Tests {
			if test.Disabled {
				r.handleDisabled(i, ch)
			} else if !test.Enabled {
				r.handleSkipped(i, ch)
			} else if r.scenario.Parallel {
				r.runAsync(client, context, i, ch)
			} else if !r.runSync(client, context, i, ch) {
				break
			}
		}
	}()

	return ch
}

func (r *scenarioRunner) handleDisabled(index int, ch chan *TestResult) {
	ch <- &TestResult{
		Index:    index,
		Disabled: true,
	}
}

func (r *scenarioRunner) handleSkipped(index int, ch chan *TestResult) {
	ch <- &TestResult{
		Index:   index,
		Skipped: true,
	}
}

func (r *scenarioRunner) runSync(
	client *http.Client,
	context map[string]interface{},
	index int,
	ch chan *TestResult,
) bool {
	result, err := r.runTest(client, context, index)
	if err != nil {
		result = &TestResult{
			Index: index,
			Err:   err,
		}
	}

	ch <- result
	return result.Err == nil && len(result.RequestMatchErrors) == 0
}

func (r *scenarioRunner) runAsync(
	client *http.Client,
	context map[string]interface{},
	index int,
	ch chan *TestResult,
) {
	r.waitGroup.Add(1)

	go func() {
		defer r.waitGroup.Done()

		result, err := r.runTest(client, context, index)
		if err != nil {
			result = &TestResult{
				Index: index,
				Err:   err,
			}
		}

		ch <- result
	}()
}

func (r *scenarioRunner) Results() []*TestResult {
	results := make([]*TestResult, len(r.results))
	copy(results, r.results)
	return results
}

func (r *scenarioRunner) Errored() bool {
	if !r.submitted || r.running {
		return false
	}

	for _, result := range r.results {
		if result != nil && result.Errored() {
			return true
		}
	}

	return false
}

func (r *scenarioRunner) Failed() bool {
	if !r.submitted || r.running {
		return false
	}

	for _, result := range r.results {
		if result != nil && result.Failed() {
			return true
		}
	}

	return false
}

func (r *scenarioRunner) Resolved() bool {
	if !r.submitted || r.running {
		return false
	}

	return !r.Errored() && !r.Failed()
}

func (r *scenarioRunner) Duration() time.Duration {
	if r.running {
		return 0
	}

	duration := time.Duration(0)
	for _, result := range r.results {
		if result != nil {
			duration += result.Duration
		}
	}

	return duration
}

func (r *scenarioRunner) runTest(
	client *http.Client,
	context map[string]interface{},
	index int,
) (testResult *TestResult, err error) {
	if r.testSemaphore != nil {
		_ = r.testSemaphore.Acquire(r.ctx, 1)
		defer r.testSemaphore.Release(1)
	}

	test := r.scenario.Tests[index]

	// TODO - should build each time?
	r.mutex.RLock()
	req, reqBody, err := buildRequest(test.Request, context)
	r.mutex.RUnlock()

	if err != nil {
		return nil, err
	}

	prefix := logging.NewPrefix(r.scenario.Name, test.Name)

	for i := 0; i <= test.Retries; i++ {
		if i > 0 {
			r.logger.Log(
				prefix,
				r.logger.Colorize(
					logging.ColorWarn,
					"Previous attempt failed, retrying in %s", test.RetryInterval,
				),
			)

			time.Sleep(test.RetryInterval)
		}

		r.logger.Log(prefix, "Attempting request...")

		if r.verbosityLevel == logging.VerbosityLevelRequestResponse {
			r.logger.Log(prefix, formatRequest(req, reqBody, r.logger.Colorized()))
		}

		started := time.Now()
		resp, err := client.Do(req)
		duration := time.Now().Sub(started)

		if err != nil {
			return nil, err
		}

		respBody, extraction, errors, err := matchResponse(resp, test.Response)
		if err != nil {
			return nil, err
		}

		if r.verbosityLevel == logging.VerbosityLevelRequestResponse {
			r.logger.Log(prefix, formatResponse(resp, respBody, r.logger.Colorized()))
		}

		r.mutex.Lock()
		context[test.Name] = extraction
		r.mutex.Unlock()

		testResult = &TestResult{
			Index:              index,
			Request:            req,
			RequestBody:        reqBody,
			Response:           resp,
			ResponseBody:       respBody,
			RequestMatchErrors: errors,
			Duration:           duration,
		}

		if len(errors) == 0 {
			r.logger.Log(
				prefix,
				r.logger.Colorize(
					logging.ColorInfo,
					"Response matched expectation - test passed",
				),
			)

			break
		} else {

			if i == test.Retries {
				r.logger.Log(
					prefix,
					r.logger.Colorize(
						logging.ColorError,
						"Response did not match expectation - test failed",
					),
				)
			} else {
				r.logger.Log(
					prefix,
					r.logger.Colorize(
						logging.ColorWarn,
						"Response did not match expectation",
					),
				)
			}
		}
	}

	return
}
