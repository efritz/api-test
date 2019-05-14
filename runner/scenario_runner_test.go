package runner

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/api-test/config"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

type ScenarioRunnerSuite struct{}

func (s *ScenarioRunnerSuite) TestRun(t sweet.T) {
	paths := []string{"/t1", "/t2", "/t3"}
	reqBodies := []string{`{"req": "r1"}`, `{"req": "r2"}`, `{"req": "r3"}`}
	respBodies := []string{`{"resp": "r1"}`, `{"resp": "r2"}`, `{"resp": "r3"}`}

	server := ghttp.NewServer()
	server.AppendHandlers(ghttp.RespondWith(http.StatusOK, respBodies[0]))
	server.AppendHandlers(ghttp.RespondWith(http.StatusOK, respBodies[1]))
	server.AppendHandlers(ghttp.RespondWith(http.StatusOK, respBodies[2]))

	scenario := &config.Scenario{
		Tests: []*config.Test{
			&config.Test{
				Request: &config.Request{
					URL:  testTemplate(server.URL() + paths[0]),
					Body: testTemplate(reqBodies[0]),
				},
				Response: &config.Response{
					Status: testPattern("2.."),
					Body:   testPattern(`r\d`),
				},
				Enabled: true,
			},
			&config.Test{
				Request: &config.Request{
					URL:  testTemplate(server.URL() + paths[1]),
					Body: testTemplate(reqBodies[1]),
				},
				Response: &config.Response{
					Status: testPattern("2.."),
					Body:   testPattern(`r\d`),
				},
				Enabled: true,
			},
			&config.Test{
				Request: &config.Request{
					URL:  testTemplate(server.URL() + paths[2]),
					Body: testTemplate(reqBodies[2]),
				},
				Response: &config.Response{
					Status: testPattern("2.."),
					Body:   testPattern(`r\d`),
				},
				Enabled: true,
			},
		},
	}

	runner := NewScenarioRunner(scenario, false)
	runner.Run(http.DefaultClient, map[string]interface{}{})
	Expect(runner.Resolved()).To(BeTrue())
	Expect(runner.Errored()).To(BeFalse())
	Expect(runner.Failed()).To(BeFalse())

	Expect(server.ReceivedRequests()).To(HaveLen(3))
	Expect(server.ReceivedRequests()[0].URL.Path).To(Equal(paths[0]))
	Expect(server.ReceivedRequests()[1].URL.Path).To(Equal(paths[1]))
	Expect(server.ReceivedRequests()[2].URL.Path).To(Equal(paths[2]))

	for i, result := range runner.Results() {
		Expect(result.Index).To(Equal(i))
		Expect(result.Disabled).To(BeFalse())
		Expect(result.Skipped).To(BeFalse())
		Expect(result.Request.URL.Path).To(Equal(paths[i]))
		Expect(result.RequestBody).To(Equal(reqBodies[i]))
		Expect(result.Response.StatusCode).To(Equal(http.StatusOK))
		Expect(result.ResponseBody).To(Equal(respBodies[i]))
		Expect(result.RequestMatchErrors).To(HaveLen(0))
		Expect(result.Err).To(BeNil())
	}
}

func (s *ScenarioRunnerSuite) TestRunFailure(t sweet.T) {
	paths := []string{"/t1", "/t2", "/t3"}
	reqBodies := []string{`{"req": "r1"}`, `{"req": "r2"}`, `{"req": "r3"}`}
	respBodies := []string{`{"resp": "r1"}`, `{"resp": "r2"}`, `{"resp": "r3"}`}

	server := ghttp.NewServer()
	server.AppendHandlers(ghttp.RespondWith(http.StatusOK, respBodies[0]))
	server.AppendHandlers(ghttp.RespondWith(http.StatusOK, respBodies[1]))
	server.AppendHandlers(ghttp.RespondWith(http.StatusOK, respBodies[2]))

	scenario := &config.Scenario{
		Tests: []*config.Test{
			&config.Test{
				Request: &config.Request{
					URL:  testTemplate(server.URL() + paths[0]),
					Body: testTemplate(reqBodies[0]),
				},
				Response: &config.Response{
					Status: testPattern("2.."),
					Body:   testPattern(`r\d`),
				},
				Enabled: true,
			},
			&config.Test{
				Request: &config.Request{
					URL:  testTemplate(server.URL() + paths[1]),
					Body: testTemplate(reqBodies[1]),
				},
				Response: &config.Response{
					Status: testPattern("2.."),
					Body:   testPattern(`missing pattern`),
				},
				Enabled: true,
			},
			&config.Test{
				Request: &config.Request{
					URL:  testTemplate(server.URL() + paths[2]),
					Body: testTemplate(reqBodies[2]),
				},
				Response: &config.Response{
					Status: testPattern("2.."),
					Body:   testPattern(`r\d`),
				},
				Enabled: true,
			},
		},
	}

	runner := NewScenarioRunner(scenario, false)
	runner.Run(http.DefaultClient, map[string]interface{}{})
	Expect(runner.Resolved()).To(BeFalse())
	Expect(runner.Errored()).To(BeFalse())
	Expect(runner.Failed()).To(BeTrue())

	results := runner.Results()
	Expect(results[0].RequestMatchErrors).To(HaveLen(0))
	Expect(results[1].RequestMatchErrors).To(HaveLen(1))
	Expect(results[1].RequestMatchErrors[0]).To(Equal(RequestMatchError{
		Type:     "Body",
		Expected: "missing pattern",
		Actual:   "<placeholder>",
	}))
}

func (s *ScenarioRunnerSuite) TestRunParallel(t sweet.T) {
	numTests := 20
	started := make(chan time.Time, numTests)
	responded := make(chan time.Time, numTests)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started <- time.Now()
		<-time.After(time.Millisecond * 25)
		responded <- time.Now()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	scenario := &config.Scenario{
		Tests:    []*config.Test{},
		Parallel: true,
	}

	for i := 0; i < numTests; i++ {
		scenario.Tests = append(scenario.Tests, &config.Test{
			Request: &config.Request{
				URL: testTemplate(ts.URL + fmt.Sprintf("/t%d", i+1)),
			},
			Response: &config.Response{
				Status: testPattern("2.."),
			},
			Enabled: true,
		})
	}

	runner := NewScenarioRunner(scenario, false)
	runner.Run(http.DefaultClient, map[string]interface{}{})

	maxStarted := <-started
	minResponded := <-responded

	for i := 1; i < numTests; i++ {
		if t := <-started; maxStarted.Before(t) {
			maxStarted = t
		}

		if t := <-responded; t.Before(minResponded) {
			minResponded = t
		}
	}

	Expect(maxStarted.Before(minResponded)).To(BeTrue())
}

func (s *ScenarioRunnerSuite) TestRunParallelFailure(t sweet.T) {
	paths := []string{"/t1", "/t2", "/t3"}
	reqBodies := []string{`{"req": "r1"}`, `{"req": "r2"}`, `{"req": "r3"}`}
	respBodies := []string{`{"resp": "r1"}`, `{"resp": "r2"}`, `{"resp": "r3"}`}

	server := ghttp.NewServer()
	server.AppendHandlers(ghttp.RespondWith(http.StatusOK, respBodies[0]))
	server.AppendHandlers(ghttp.RespondWith(http.StatusOK, respBodies[1]))
	server.AppendHandlers(ghttp.RespondWith(http.StatusOK, respBodies[2]))

	scenario := &config.Scenario{
		Tests: []*config.Test{
			&config.Test{
				Request: &config.Request{
					URL:  testTemplate(server.URL() + paths[0]),
					Body: testTemplate(reqBodies[0]),
				},
				Response: &config.Response{
					Status: testPattern("2.."),
					Body:   testPattern(`r\d`),
				},
				Enabled: true,
			},
			&config.Test{
				Request: &config.Request{
					URL:  testTemplate(server.URL() + paths[1]),
					Body: testTemplate(reqBodies[1]),
				},
				Response: &config.Response{
					Status: testPattern("2.."),
					Body:   testPattern(`missing pattern`),
				},
				Enabled: true,
			},
			&config.Test{
				Request: &config.Request{
					URL:  testTemplate(server.URL() + paths[2]),
					Body: testTemplate(reqBodies[2]),
				},
				Response: &config.Response{
					Status: testPattern("2.."),
					Body:   testPattern(`r\d`),
				},
				Enabled: true,
			},
		},
		Parallel: true,
	}

	runner := NewScenarioRunner(scenario, false)
	runner.Run(http.DefaultClient, map[string]interface{}{})

	results := runner.Results()
	Expect(results[0].RequestMatchErrors).To(HaveLen(0))
	Expect(results[1].RequestMatchErrors).To(HaveLen(1))
	Expect(results[1].RequestMatchErrors[0]).To(Equal(RequestMatchError{
		Type:     "Body",
		Expected: "missing pattern",
		Actual:   "<placeholder>",
	}))
	Expect(results[2].RequestMatchErrors).To(HaveLen(0))
}

func (s *ScenarioRunnerSuite) TestRunDisabled(t sweet.T) {
	scenario := &config.Scenario{
		Tests: []*config.Test{
			&config.Test{
				Disabled: true,
			},
		},
	}

	runner := NewScenarioRunner(scenario, false)
	runner.Run(http.DefaultClient, map[string]interface{}{})

	results := runner.Results()
	Expect(results).To(BeEmpty())
}

func (s *ScenarioRunnerSuite) TestRunFirstDisabled(t sweet.T) {
	server := ghttp.NewServer()
	server.AppendHandlers(ghttp.RespondWith(http.StatusOK, "ok"))

	scenario := &config.Scenario{
		Tests: []*config.Test{
			&config.Test{
				Disabled: true,
			},
			&config.Test{
				Request: &config.Request{
					URL:  testTemplate(server.URL()),
					Body: testTemplate(""),
				},
				Response: &config.Response{
					Status: testPattern("2.."),
					Body:   testPattern(`ok`),
				},
				Enabled: true,
			},
		},
	}

	runner := NewScenarioRunner(scenario, false)
	runner.Run(http.DefaultClient, map[string]interface{}{})

	results := runner.Results()
	Expect(results[0]).To(BeNil())
	Expect(results[1].RequestMatchErrors).To(HaveLen(0))
}

func (s *ScenarioRunnerSuite) TestRunNotEnabled(t sweet.T) {
	scenario := &config.Scenario{
		Tests: []*config.Test{
			&config.Test{
				Enabled: false,
			},
		},
	}

	runner := NewScenarioRunner(scenario, false)
	runner.Run(http.DefaultClient, map[string]interface{}{})

	results := runner.Results()
	Expect(results[0]).To(Equal(&TestResult{Index: 0, Skipped: true}))
}
