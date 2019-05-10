package runner

import (
	"net/http"
	"sync"
	"time"

	"github.com/efritz/api-test/config"
)

type ScenarioContext struct {
	Scenario        *config.Scenario
	Results         []*TestResult
	Pending         bool
	Running         bool
	Skipped         bool
	Context         map[string]interface{}
	forceSequential bool
}

func NewScenarioContext(scenario *config.Scenario, forceSequential bool) *ScenarioContext {
	return &ScenarioContext{
		Scenario:        scenario,
		Results:         []*TestResult{},
		Pending:         true,
		forceSequential: forceSequential,
	}
}

func (c *ScenarioContext) Duration() time.Duration {
	duration := time.Duration(0)
	for _, result := range c.Results {
		if result != nil {
			duration += result.Duration
		}
	}

	return duration
}

func (c *ScenarioContext) Resolved() bool {
	if !c.finished() {
		return false
	}

	for _, result := range c.Results {
		if result != nil && (result.Errored() || result.Failed()) {
			return false
		}
	}

	return true
}

func (c *ScenarioContext) Errored() bool {
	if !c.finished() {
		return false
	}

	for _, result := range c.Results {
		if result != nil && result.Errored() {
			return true
		}
	}

	return false
}

func (c *ScenarioContext) Failed() bool {
	if !c.finished() {
		return false
	}

	for _, result := range c.Results {
		if result != nil && result.Failed() {
			return true
		}
	}

	return false
}

func (c *ScenarioContext) finished() bool {
	if len(c.Scenario.Tests) != len(c.Results) {
		return false
	}

	for i, result := range c.Results {
		if result == nil && !c.Scenario.Tests[i].Disabled {
			return false
		}
	}

	return true
}

func (c *ScenarioContext) LastResult() *TestResult {
	for i := len(c.Results) - 1; i >= 0; i-- {
		if c.Results[i] != nil {
			return c.Results[i]
		}
	}

	return nil
}

func (c *ScenarioContext) LastTest() *config.Test {
	for i := len(c.Results) - 1; i >= 0; i-- {
		if c.Results[i] != nil {
			return c.Scenario.Tests[i]
		}
	}

	return nil
}

func (c *ScenarioContext) Run(client *http.Client) <-chan *TestResult {
	ch := make(chan *TestResult)

	go func() {
		defer close(ch)

		wg := sync.WaitGroup{}
		defer wg.Wait()

		for i, test := range c.Scenario.Tests {
			if test.Disabled {
				result := &TestResult{
					Index:    i,
					Disabled: true,
				}

				ch <- result
				continue
			}

			if !test.Enabled {
				result := &TestResult{
					Index:   i,
					Skipped: true,
				}

				ch <- result
				continue
			}

			if c.Scenario.Parallel && !c.forceSequential {
				wg.Add(1)
				go func(index int) {
					defer wg.Done()

					result, err := c.runTest(client, index, test)
					if err != nil {
						result = &TestResult{
							Index: i,
							Err:   err,
						}
					}

					ch <- result
				}(i)

				continue
			}

			result, err := c.runTest(client, i, test)
			if err != nil {
				result = &TestResult{
					Index: i,
					Err:   err,
				}
			}

			ch <- result

			if result.Err != nil || len(result.RequestMatchErrors) > 0 {
				break
			}
		}
	}()

	return ch
}

func (c *ScenarioContext) runTest(client *http.Client, index int, test *config.Test) (*TestResult, error) {
	started := time.Now()

	req, reqBody, err := buildRequest(test.Request, c.Context)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	duration := time.Now().Sub(started)

	respBody, extraction, errors, err := matchResponse(resp, test.Response)
	if err != nil {
		return nil, err
	}

	c.Context[test.Name] = extraction

	return &TestResult{
		Index:              index,
		Request:            req,
		RequestBody:        reqBody,
		Response:           resp,
		ResponseBody:       respBody,
		RequestMatchErrors: errors,
		Duration:           duration,
	}, nil
}
