package runner

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/api-test/config"
	. "github.com/efritz/go-mockgen/matchers"
	. "github.com/onsi/gomega"
)

type RunnerSuite struct{}

func (s *RunnerSuite) TestRunDependencyChainContexts(t sweet.T) {
	mockRunner1 := newBasicMockRunner(nil, map[string]interface{}{"ret": "s1"})
	mockRunner2 := newBasicMockRunner(nil, map[string]interface{}{"ret": "s2"})
	mockRunner3 := newBasicMockRunner(nil, map[string]interface{}{"ret": "s3"})
	mockRunner4 := newBasicMockRunner(nil, map[string]interface{}{"ret": "s4"})

	factory := NewMockScenarioRunnerFactory()
	factory.NewFunc.PushReturn(mockRunner1)
	factory.NewFunc.PushReturn(mockRunner2)
	factory.NewFunc.PushReturn(mockRunner3)
	factory.NewFunc.PushReturn(mockRunner4)

	config := &config.Config{
		Options: &config.Options{},
		Scenarios: map[string]*config.Scenario{
			"s1": &config.Scenario{
				Name:    "s1",
				Enabled: true,
			},
			"s2": &config.Scenario{
				Name:         "s2",
				Dependencies: []string{"s1"},
				Enabled:      true,
			},
			"s3": &config.Scenario{
				Name:    "s3",
				Enabled: true,
			},
			"s4": &config.Scenario{
				Name:         "s4",
				Dependencies: []string{"s2", "s3"},
				Enabled:      true,
			},
		},
	}

	runner := NewRunner(config, WithScenarioRunnerFactory(factory))
	err := runner.Run()
	Expect(err).To(BeNil())

	Expect(mockRunner2.RunFunc).To(BeCalledOnceWith(
		BeAnything(),
		HaveKeyWithValue("s1", map[string]interface{}{"ret": "s1"}),
	))

	Expect(mockRunner4.RunFunc).To(BeCalledOnceWith(
		BeAnything(),
		And(
			HaveKeyWithValue("s2", map[string]interface{}{"ret": "s2"}),
			HaveKeyWithValue("s3", map[string]interface{}{"ret": "s3"}),
		),
	))
}

func (s *RunnerSuite) TestRunSkipsScenariosWithFailedDependency(t sweet.T) {
	failure := &TestResult{
		Err: fmt.Errorf("oops"),
	}

	mockRunner1 := newBasicMockRunner(nil, nil)
	mockRunner2 := newBasicMockRunner([]*TestResult{failure}, nil)
	mockRunner3 := newBasicMockRunner(nil, nil)

	factory := NewMockScenarioRunnerFactory()
	factory.NewFunc.PushReturn(mockRunner1)
	factory.NewFunc.PushReturn(mockRunner2)
	factory.NewFunc.PushReturn(mockRunner3)

	config := &config.Config{
		Options: &config.Options{},
		Scenarios: map[string]*config.Scenario{
			"s1": &config.Scenario{
				Name:    "s1",
				Enabled: true,
			},
			"s2": &config.Scenario{
				Name:         "s2",
				Dependencies: []string{"s1"},
				Enabled:      true,
				Tests:        []*config.Test{&config.Test{}},
			},
			"s3": &config.Scenario{
				Name:         "s3",
				Dependencies: []string{"s2"},
				Enabled:      true,
			},
		},
	}

	runner := NewRunner(config, WithScenarioRunnerFactory(factory))
	err := runner.Run()
	Expect(err).To(BeNil())
	Expect(mockRunner3.RunFunc).NotTo(BeCalled())
}

func (s *RunnerSuite) TestRunSkipsScenarioWithSkippedDependency(t sweet.T) {
	mockRunner1 := newBasicMockRunner(nil, nil)
	mockRunner2 := newBasicMockRunner(nil, nil)
	mockRunner3 := newBasicMockRunner(nil, nil)

	factory := NewMockScenarioRunnerFactory()
	factory.NewFunc.PushReturn(mockRunner1)
	factory.NewFunc.PushReturn(mockRunner2)
	factory.NewFunc.PushReturn(mockRunner3)

	config := &config.Config{
		Options: &config.Options{},
		Scenarios: map[string]*config.Scenario{
			"s1": &config.Scenario{
				Name:    "s1",
				Enabled: true,
			},
			"s2": &config.Scenario{
				Name:         "s2",
				Dependencies: []string{"s1"},
				Disabled:     true,
			},
			"s3": &config.Scenario{
				Name:         "s3",
				Dependencies: []string{"s2"},
				Enabled:      true,
			},
		},
	}

	runner := NewRunner(config, WithScenarioRunnerFactory(factory))
	err := runner.Run()
	Expect(err).To(BeNil())
	Expect(mockRunner2.RunFunc).NotTo(BeCalled())
	Expect(mockRunner3.RunFunc).NotTo(BeCalled())
}

func (s *RunnerSuite) TestRunDependencyOrder(t sweet.T) {
	sequence := make(chan string, 6)
	defer close(sequence)

	hook := func(name string) func() map[string]interface{} {
		return func() map[string]interface{} {
			sequence <- fmt.Sprintf("%s1", name)
			<-time.After(time.Millisecond * 25)
			sequence <- fmt.Sprintf("%s2", name)
			return nil
		}
	}

	mockRunner1 := newMockRunner(nil, hook("A"))
	mockRunner2 := newMockRunner(nil, hook("B"))
	mockRunner3 := newMockRunner(nil, hook("C"))

	factory := NewMockScenarioRunnerFactory()
	factory.NewFunc.PushReturn(mockRunner1)
	factory.NewFunc.PushReturn(mockRunner2)
	factory.NewFunc.PushReturn(mockRunner3)

	config := &config.Config{
		Options: &config.Options{},
		Scenarios: map[string]*config.Scenario{
			"s1": &config.Scenario{
				Name:    "s1",
				Enabled: true,
			},
			"s2": &config.Scenario{
				Name:         "s2",
				Dependencies: []string{"s1"},
				Enabled:      true,
			},
			"s3": &config.Scenario{
				Name:         "s3",
				Dependencies: []string{"s2"},
				Enabled:      true,
			},
		},
	}

	runner := NewRunner(config, WithScenarioRunnerFactory(factory))
	err := runner.Run()
	Expect(err).To(BeNil())

	Expect(sequence).To(Receive(Equal("A1")))
	Expect(sequence).To(Receive(Equal("A2")))
	Expect(sequence).To(Receive(Equal("B1")))
	Expect(sequence).To(Receive(Equal("B2")))
	Expect(sequence).To(Receive(Equal("C1")))
	Expect(sequence).To(Receive(Equal("C2")))
}

func (s *RunnerSuite) TestRunParallel(t sweet.T) {
	numTests := 5
	sequence := make(chan string, 6*numTests)
	defer close(sequence)

	hook := func() map[string]interface{} {
		sequence <- "A"
		<-time.After(time.Millisecond * 25)
		sequence <- "B"
		return nil
	}

	mockRunner1 := newMockRunner(nil, hook)
	mockRunner2 := newMockRunner(nil, hook)
	mockRunner3 := newMockRunner(nil, hook)

	for i := 0; i < numTests; i++ {
		factory := NewMockScenarioRunnerFactory()
		factory.NewFunc.PushReturn(mockRunner1)
		factory.NewFunc.PushReturn(mockRunner2)
		factory.NewFunc.PushReturn(mockRunner3)

		config := &config.Config{
			Options: &config.Options{},
			Scenarios: map[string]*config.Scenario{
				"s1": &config.Scenario{
					Name:    "s1",
					Enabled: true,
				},
				"s2": &config.Scenario{
					Name:    "s2",
					Enabled: true,
				},
				"s3": &config.Scenario{
					Name:    "s3",
					Enabled: true,
				},
			},
		}

		runner := NewRunner(config, WithScenarioRunnerFactory(factory))
		err := runner.Run()
		Expect(err).To(BeNil())
	}

	nonParallel := false
	for i := 0; i < numTests; i++ {
		str := ""
		for j := 0; j < 6; j++ {
			str += <-sequence
		}

		if str != "ABABAB" {
			nonParallel = true
		}
	}

	Expect(nonParallel).To(BeTrue())
}

func (s *RunnerSuite) TestRunForceSequential(t sweet.T) {
	sequence := make(chan string, 6)
	defer close(sequence)

	hook := func() map[string]interface{} {
		sequence <- "A"
		<-time.After(time.Millisecond * 25)
		sequence <- "B"
		return nil
	}

	mockRunner1 := newMockRunner(nil, hook)
	mockRunner2 := newMockRunner(nil, hook)
	mockRunner3 := newMockRunner(nil, hook)

	factory := NewMockScenarioRunnerFactory()
	factory.NewFunc.PushReturn(mockRunner1)
	factory.NewFunc.PushReturn(mockRunner2)
	factory.NewFunc.PushReturn(mockRunner3)

	config := &config.Config{
		Options: &config.Options{
			ForceSequential: true,
		},
		Scenarios: map[string]*config.Scenario{
			"s1": &config.Scenario{
				Name:    "s1",
				Enabled: true,
			},
			"s2": &config.Scenario{
				Name:    "s2",
				Enabled: true,
			},
			"s3": &config.Scenario{
				Name:    "s3",
				Enabled: true,
			},
		},
	}

	runner := NewRunner(config, WithScenarioRunnerFactory(factory))
	err := runner.Run()
	Expect(err).To(BeNil())

	Expect(sequence).To(Receive(Equal("A")))
	Expect(sequence).To(Receive(Equal("B")))
	Expect(sequence).To(Receive(Equal("A")))
	Expect(sequence).To(Receive(Equal("B")))
	Expect(sequence).To(Receive(Equal("A")))
	Expect(sequence).To(Receive(Equal("B")))
}

//
// Helpers

func newBasicMockRunner(results []*TestResult, context map[string]interface{}) *MockScenarioRunner {
	return newMockRunner(results, func() map[string]interface{} {
		return context
	})
}

func newMockRunner(results []*TestResult, hook func() map[string]interface{}) *MockScenarioRunner {
	ran := false
	errored := false
	failed := false

	for _, result := range results {
		if result.Errored() {
			errored = true
		}

		if result.Failed() {
			failed = true
		}
	}

	runner := NewMockScenarioRunner()
	runner.RunFunc.SetDefaultHook(func(*http.Client, map[string]interface{}) map[string]interface{} {
		defer func() {
			ran = true
		}()

		return hook()
	})

	runner.ResultsFunc.SetDefaultHook(func() []*TestResult { return results })
	runner.ErroredFunc.SetDefaultHook(func() bool { return ran && errored })
	runner.FailedFunc.SetDefaultHook(func() bool { return ran && failed })
	runner.ResolvedFunc.SetDefaultHook(func() bool { return ran && !runner.Errored() && !runner.Failed() })

	return runner
}
