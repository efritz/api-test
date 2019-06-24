// Code generated by github.com/efritz/go-mockgen 0.1.0; DO NOT EDIT.
// This file was generated by robots at
// 2019-06-24T09:16:11-05:00
// using the command
// $ go-mockgen -f github.com/efritz/api-test/runner -i ScenarioRunnerFactory -o scenario_runner_factory_mock_test.go

package runner

import (
	config "github.com/efritz/api-test/config"
	logging "github.com/efritz/api-test/logging"
	"sync"
)

// MockScenarioRunnerFactory is a mock impelementation of the
// ScenarioRunnerFactory interface (from the package
// github.com/efritz/api-test/runner) used for unit testing.
type MockScenarioRunnerFactory struct {
	// NewFunc is an instance of a mock function object controlling the
	// behavior of the method New.
	NewFunc *ScenarioRunnerFactoryNewFunc
}

// NewMockScenarioRunnerFactory creates a new mock of the
// ScenarioRunnerFactory interface. All methods return zero values for all
// results, unless overwritten.
func NewMockScenarioRunnerFactory() *MockScenarioRunnerFactory {
	return &MockScenarioRunnerFactory{
		NewFunc: &ScenarioRunnerFactoryNewFunc{
			defaultHook: func(*config.Scenario, logging.Logger, bool) ScenarioRunner {
				return nil
			},
		},
	}
}

// NewMockScenarioRunnerFactoryFrom creates a new mock of the
// MockScenarioRunnerFactory interface. All methods delegate to the given
// implementation, unless overwritten.
func NewMockScenarioRunnerFactoryFrom(i ScenarioRunnerFactory) *MockScenarioRunnerFactory {
	return &MockScenarioRunnerFactory{
		NewFunc: &ScenarioRunnerFactoryNewFunc{
			defaultHook: i.New,
		},
	}
}

// ScenarioRunnerFactoryNewFunc describes the behavior when the New method
// of the parent MockScenarioRunnerFactory instance is invoked.
type ScenarioRunnerFactoryNewFunc struct {
	defaultHook func(*config.Scenario, logging.Logger, bool) ScenarioRunner
	hooks       []func(*config.Scenario, logging.Logger, bool) ScenarioRunner
	history     []ScenarioRunnerFactoryNewFuncCall
	mutex       sync.Mutex
}

// New delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockScenarioRunnerFactory) New(v0 *config.Scenario, v1 logging.Logger, v2 bool) ScenarioRunner {
	r0 := m.NewFunc.nextHook()(v0, v1, v2)
	m.NewFunc.appendCall(ScenarioRunnerFactoryNewFuncCall{v0, v1, v2, r0})
	return r0
}

// SetDefaultHook sets function that is called when the New method of the
// parent MockScenarioRunnerFactory instance is invoked and the hook queue
// is empty.
func (f *ScenarioRunnerFactoryNewFunc) SetDefaultHook(hook func(*config.Scenario, logging.Logger, bool) ScenarioRunner) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// New method of the parent MockScenarioRunnerFactory instance inovkes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *ScenarioRunnerFactoryNewFunc) PushHook(hook func(*config.Scenario, logging.Logger, bool) ScenarioRunner) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ScenarioRunnerFactoryNewFunc) SetDefaultReturn(r0 ScenarioRunner) {
	f.SetDefaultHook(func(*config.Scenario, logging.Logger, bool) ScenarioRunner {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ScenarioRunnerFactoryNewFunc) PushReturn(r0 ScenarioRunner) {
	f.PushHook(func(*config.Scenario, logging.Logger, bool) ScenarioRunner {
		return r0
	})
}

func (f *ScenarioRunnerFactoryNewFunc) nextHook() func(*config.Scenario, logging.Logger, bool) ScenarioRunner {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ScenarioRunnerFactoryNewFunc) appendCall(r0 ScenarioRunnerFactoryNewFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ScenarioRunnerFactoryNewFuncCall objects
// describing the invocations of this function.
func (f *ScenarioRunnerFactoryNewFunc) History() []ScenarioRunnerFactoryNewFuncCall {
	f.mutex.Lock()
	history := make([]ScenarioRunnerFactoryNewFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ScenarioRunnerFactoryNewFuncCall is an object that describes an
// invocation of method New on an instance of MockScenarioRunnerFactory.
type ScenarioRunnerFactoryNewFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 *config.Scenario
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 logging.Logger
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 bool
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 ScenarioRunner
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ScenarioRunnerFactoryNewFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ScenarioRunnerFactoryNewFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}
