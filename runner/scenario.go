package runner

import (
	"net/http"
	"time"

	"github.com/efritz/api-test/config"
)

type ScenarioContext struct {
	Scenario *config.Scenario
	Results  []*TestResult
	Pending  bool
	Running  bool
	Skipped  bool
	Context  map[string]interface{}
}

func NewScenarioContext(scenario *config.Scenario) *ScenarioContext {
	return &ScenarioContext{
		Scenario: scenario,
		Results:  []*TestResult{},
		Pending:  true,
	}
}

func (c *ScenarioContext) Duration() time.Duration {
	duration := time.Duration(0)
	for _, result := range c.Results {
		duration += result.Duration
	}

	return duration
}

func (c *ScenarioContext) Resolved() bool {
	return len(c.Results) == len(c.Scenario.Tests) && !c.Failed()
}

func (c *ScenarioContext) Failed() bool {
	if result := c.LastResult(); result != nil {
		return result.Failed()
	}

	return false
}

func (c *ScenarioContext) LastResult() *TestResult {
	if len(c.Results) == 0 {
		return nil
	}

	return c.Results[len(c.Results)-1]
}

func (c *ScenarioContext) LastTest() *config.Test {
	if len(c.Results) == 0 {
		return nil
	}

	return c.Scenario.Tests[len(c.Results)-1]
}

func (c *ScenarioContext) Run(client *http.Client) <-chan *TestResult {
	ch := make(chan *TestResult)

	go func() {
		defer close(ch)

		for _, test := range c.Scenario.Tests {
			result, err := c.runTest(client, test)
			if err != nil {
				result = &TestResult{
					Err: err,
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

func (c *ScenarioContext) runTest(client *http.Client, test *config.Test) (*TestResult, error) {
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
		Request:            req,
		RequestBody:        reqBody,
		Response:           resp,
		ResponseBody:       respBody,
		RequestMatchErrors: errors,
		Duration:           duration,
	}, nil
}
