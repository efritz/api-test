package runner

import (
	"net/http"
	"sync"
	"time"

	"github.com/efritz/api-test/config"
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
		scenario        *config.Scenario
		results         []*TestResult
		forceSequential bool
		running         bool
		submitted       bool
		mutex           sync.RWMutex
		waitGroup       sync.WaitGroup
	}
)

func NewScenarioRunner(scenario *config.Scenario, forceSequential bool) ScenarioRunner {
	return newScenarioRunner(scenario, forceSequential)
}

func newScenarioRunner(scenario *config.Scenario, forceSequential bool) *scenarioRunner {
	return &scenarioRunner{
		scenario:        scenario,
		results:         []*TestResult{},
		forceSequential: forceSequential,
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
			} else if r.scenario.Parallel && !r.forceSequential {
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
	test := r.scenario.Tests[index]

	r.mutex.RLock()
	req, reqBody, err := buildRequest(test.Request, context)
	r.mutex.RUnlock()

	if err != nil {
		return nil, err
	}

	for i := 0; i <= test.Retries; i++ {
		if i > 0 {
			time.Sleep(test.RetryInterval)
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
			break
		}
	}

	return
}
