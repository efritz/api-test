// Code generated by github.com/efritz/go-mockgen 0.1.0; DO NOT EDIT.
// This file was generated by robots at
// 2019-06-20T08:59:16-05:00
// using the command
// $ go-mockgen -f github.com/efritz/api-test/runner -i ScenarioRunner -o scenario_runner_mock_test.go

package runner

import (
	"net/http"
	"sync"
	"time"
)

// MockScenarioRunner is a mock impelementation of the ScenarioRunner
// interface (from the package github.com/efritz/api-test/runner) used for
// unit testing.
type MockScenarioRunner struct {
	// DurationFunc is an instance of a mock function object controlling the
	// behavior of the method Duration.
	DurationFunc *ScenarioRunnerDurationFunc
	// ErroredFunc is an instance of a mock function object controlling the
	// behavior of the method Errored.
	ErroredFunc *ScenarioRunnerErroredFunc
	// FailedFunc is an instance of a mock function object controlling the
	// behavior of the method Failed.
	FailedFunc *ScenarioRunnerFailedFunc
	// ResolvedFunc is an instance of a mock function object controlling the
	// behavior of the method Resolved.
	ResolvedFunc *ScenarioRunnerResolvedFunc
	// ResultsFunc is an instance of a mock function object controlling the
	// behavior of the method Results.
	ResultsFunc *ScenarioRunnerResultsFunc
	// RunFunc is an instance of a mock function object controlling the
	// behavior of the method Run.
	RunFunc *ScenarioRunnerRunFunc
}

// NewMockScenarioRunner creates a new mock of the ScenarioRunner interface.
// All methods return zero values for all results, unless overwritten.
func NewMockScenarioRunner() *MockScenarioRunner {
	return &MockScenarioRunner{
		DurationFunc: &ScenarioRunnerDurationFunc{
			defaultHook: func() time.Duration {
				return 0
			},
		},
		ErroredFunc: &ScenarioRunnerErroredFunc{
			defaultHook: func() bool {
				return false
			},
		},
		FailedFunc: &ScenarioRunnerFailedFunc{
			defaultHook: func() bool {
				return false
			},
		},
		ResolvedFunc: &ScenarioRunnerResolvedFunc{
			defaultHook: func() bool {
				return false
			},
		},
		ResultsFunc: &ScenarioRunnerResultsFunc{
			defaultHook: func() []*TestResult {
				return nil
			},
		},
		RunFunc: &ScenarioRunnerRunFunc{
			defaultHook: func(*http.Client, map[string]interface{}) map[string]interface{} {
				return nil
			},
		},
	}
}

// NewMockScenarioRunnerFrom creates a new mock of the MockScenarioRunner
// interface. All methods delegate to the given implementation, unless
// overwritten.
func NewMockScenarioRunnerFrom(i ScenarioRunner) *MockScenarioRunner {
	return &MockScenarioRunner{
		DurationFunc: &ScenarioRunnerDurationFunc{
			defaultHook: i.Duration,
		},
		ErroredFunc: &ScenarioRunnerErroredFunc{
			defaultHook: i.Errored,
		},
		FailedFunc: &ScenarioRunnerFailedFunc{
			defaultHook: i.Failed,
		},
		ResolvedFunc: &ScenarioRunnerResolvedFunc{
			defaultHook: i.Resolved,
		},
		ResultsFunc: &ScenarioRunnerResultsFunc{
			defaultHook: i.Results,
		},
		RunFunc: &ScenarioRunnerRunFunc{
			defaultHook: i.Run,
		},
	}
}

// ScenarioRunnerDurationFunc describes the behavior when the Duration
// method of the parent MockScenarioRunner instance is invoked.
type ScenarioRunnerDurationFunc struct {
	defaultHook func() time.Duration
	hooks       []func() time.Duration
	history     []ScenarioRunnerDurationFuncCall
	mutex       sync.Mutex
}

// Duration delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockScenarioRunner) Duration() time.Duration {
	r0 := m.DurationFunc.nextHook()()
	m.DurationFunc.appendCall(ScenarioRunnerDurationFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Duration method of
// the parent MockScenarioRunner instance is invoked and the hook queue is
// empty.
func (f *ScenarioRunnerDurationFunc) SetDefaultHook(hook func() time.Duration) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Duration method of the parent MockScenarioRunner instance inovkes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *ScenarioRunnerDurationFunc) PushHook(hook func() time.Duration) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ScenarioRunnerDurationFunc) SetDefaultReturn(r0 time.Duration) {
	f.SetDefaultHook(func() time.Duration {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ScenarioRunnerDurationFunc) PushReturn(r0 time.Duration) {
	f.PushHook(func() time.Duration {
		return r0
	})
}

func (f *ScenarioRunnerDurationFunc) nextHook() func() time.Duration {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ScenarioRunnerDurationFunc) appendCall(r0 ScenarioRunnerDurationFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ScenarioRunnerDurationFuncCall objects
// describing the invocations of this function.
func (f *ScenarioRunnerDurationFunc) History() []ScenarioRunnerDurationFuncCall {
	f.mutex.Lock()
	history := make([]ScenarioRunnerDurationFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ScenarioRunnerDurationFuncCall is an object that describes an invocation
// of method Duration on an instance of MockScenarioRunner.
type ScenarioRunnerDurationFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 time.Duration
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ScenarioRunnerDurationFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ScenarioRunnerDurationFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ScenarioRunnerErroredFunc describes the behavior when the Errored method
// of the parent MockScenarioRunner instance is invoked.
type ScenarioRunnerErroredFunc struct {
	defaultHook func() bool
	hooks       []func() bool
	history     []ScenarioRunnerErroredFuncCall
	mutex       sync.Mutex
}

// Errored delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockScenarioRunner) Errored() bool {
	r0 := m.ErroredFunc.nextHook()()
	m.ErroredFunc.appendCall(ScenarioRunnerErroredFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Errored method of
// the parent MockScenarioRunner instance is invoked and the hook queue is
// empty.
func (f *ScenarioRunnerErroredFunc) SetDefaultHook(hook func() bool) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Errored method of the parent MockScenarioRunner instance inovkes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ScenarioRunnerErroredFunc) PushHook(hook func() bool) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ScenarioRunnerErroredFunc) SetDefaultReturn(r0 bool) {
	f.SetDefaultHook(func() bool {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ScenarioRunnerErroredFunc) PushReturn(r0 bool) {
	f.PushHook(func() bool {
		return r0
	})
}

func (f *ScenarioRunnerErroredFunc) nextHook() func() bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ScenarioRunnerErroredFunc) appendCall(r0 ScenarioRunnerErroredFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ScenarioRunnerErroredFuncCall objects
// describing the invocations of this function.
func (f *ScenarioRunnerErroredFunc) History() []ScenarioRunnerErroredFuncCall {
	f.mutex.Lock()
	history := make([]ScenarioRunnerErroredFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ScenarioRunnerErroredFuncCall is an object that describes an invocation
// of method Errored on an instance of MockScenarioRunner.
type ScenarioRunnerErroredFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 bool
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ScenarioRunnerErroredFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ScenarioRunnerErroredFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ScenarioRunnerFailedFunc describes the behavior when the Failed method of
// the parent MockScenarioRunner instance is invoked.
type ScenarioRunnerFailedFunc struct {
	defaultHook func() bool
	hooks       []func() bool
	history     []ScenarioRunnerFailedFuncCall
	mutex       sync.Mutex
}

// Failed delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockScenarioRunner) Failed() bool {
	r0 := m.FailedFunc.nextHook()()
	m.FailedFunc.appendCall(ScenarioRunnerFailedFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Failed method of the
// parent MockScenarioRunner instance is invoked and the hook queue is
// empty.
func (f *ScenarioRunnerFailedFunc) SetDefaultHook(hook func() bool) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Failed method of the parent MockScenarioRunner instance inovkes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ScenarioRunnerFailedFunc) PushHook(hook func() bool) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ScenarioRunnerFailedFunc) SetDefaultReturn(r0 bool) {
	f.SetDefaultHook(func() bool {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ScenarioRunnerFailedFunc) PushReturn(r0 bool) {
	f.PushHook(func() bool {
		return r0
	})
}

func (f *ScenarioRunnerFailedFunc) nextHook() func() bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ScenarioRunnerFailedFunc) appendCall(r0 ScenarioRunnerFailedFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ScenarioRunnerFailedFuncCall objects
// describing the invocations of this function.
func (f *ScenarioRunnerFailedFunc) History() []ScenarioRunnerFailedFuncCall {
	f.mutex.Lock()
	history := make([]ScenarioRunnerFailedFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ScenarioRunnerFailedFuncCall is an object that describes an invocation of
// method Failed on an instance of MockScenarioRunner.
type ScenarioRunnerFailedFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 bool
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ScenarioRunnerFailedFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ScenarioRunnerFailedFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ScenarioRunnerResolvedFunc describes the behavior when the Resolved
// method of the parent MockScenarioRunner instance is invoked.
type ScenarioRunnerResolvedFunc struct {
	defaultHook func() bool
	hooks       []func() bool
	history     []ScenarioRunnerResolvedFuncCall
	mutex       sync.Mutex
}

// Resolved delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockScenarioRunner) Resolved() bool {
	r0 := m.ResolvedFunc.nextHook()()
	m.ResolvedFunc.appendCall(ScenarioRunnerResolvedFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Resolved method of
// the parent MockScenarioRunner instance is invoked and the hook queue is
// empty.
func (f *ScenarioRunnerResolvedFunc) SetDefaultHook(hook func() bool) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Resolved method of the parent MockScenarioRunner instance inovkes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *ScenarioRunnerResolvedFunc) PushHook(hook func() bool) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ScenarioRunnerResolvedFunc) SetDefaultReturn(r0 bool) {
	f.SetDefaultHook(func() bool {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ScenarioRunnerResolvedFunc) PushReturn(r0 bool) {
	f.PushHook(func() bool {
		return r0
	})
}

func (f *ScenarioRunnerResolvedFunc) nextHook() func() bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ScenarioRunnerResolvedFunc) appendCall(r0 ScenarioRunnerResolvedFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ScenarioRunnerResolvedFuncCall objects
// describing the invocations of this function.
func (f *ScenarioRunnerResolvedFunc) History() []ScenarioRunnerResolvedFuncCall {
	f.mutex.Lock()
	history := make([]ScenarioRunnerResolvedFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ScenarioRunnerResolvedFuncCall is an object that describes an invocation
// of method Resolved on an instance of MockScenarioRunner.
type ScenarioRunnerResolvedFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 bool
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ScenarioRunnerResolvedFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ScenarioRunnerResolvedFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ScenarioRunnerResultsFunc describes the behavior when the Results method
// of the parent MockScenarioRunner instance is invoked.
type ScenarioRunnerResultsFunc struct {
	defaultHook func() []*TestResult
	hooks       []func() []*TestResult
	history     []ScenarioRunnerResultsFuncCall
	mutex       sync.Mutex
}

// Results delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockScenarioRunner) Results() []*TestResult {
	r0 := m.ResultsFunc.nextHook()()
	m.ResultsFunc.appendCall(ScenarioRunnerResultsFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Results method of
// the parent MockScenarioRunner instance is invoked and the hook queue is
// empty.
func (f *ScenarioRunnerResultsFunc) SetDefaultHook(hook func() []*TestResult) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Results method of the parent MockScenarioRunner instance inovkes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ScenarioRunnerResultsFunc) PushHook(hook func() []*TestResult) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ScenarioRunnerResultsFunc) SetDefaultReturn(r0 []*TestResult) {
	f.SetDefaultHook(func() []*TestResult {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ScenarioRunnerResultsFunc) PushReturn(r0 []*TestResult) {
	f.PushHook(func() []*TestResult {
		return r0
	})
}

func (f *ScenarioRunnerResultsFunc) nextHook() func() []*TestResult {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ScenarioRunnerResultsFunc) appendCall(r0 ScenarioRunnerResultsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ScenarioRunnerResultsFuncCall objects
// describing the invocations of this function.
func (f *ScenarioRunnerResultsFunc) History() []ScenarioRunnerResultsFuncCall {
	f.mutex.Lock()
	history := make([]ScenarioRunnerResultsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ScenarioRunnerResultsFuncCall is an object that describes an invocation
// of method Results on an instance of MockScenarioRunner.
type ScenarioRunnerResultsFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []*TestResult
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ScenarioRunnerResultsFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ScenarioRunnerResultsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ScenarioRunnerRunFunc describes the behavior when the Run method of the
// parent MockScenarioRunner instance is invoked.
type ScenarioRunnerRunFunc struct {
	defaultHook func(*http.Client, map[string]interface{}) map[string]interface{}
	hooks       []func(*http.Client, map[string]interface{}) map[string]interface{}
	history     []ScenarioRunnerRunFuncCall
	mutex       sync.Mutex
}

// Run delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockScenarioRunner) Run(v0 *http.Client, v1 map[string]interface{}) map[string]interface{} {
	r0 := m.RunFunc.nextHook()(v0, v1)
	m.RunFunc.appendCall(ScenarioRunnerRunFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Run method of the
// parent MockScenarioRunner instance is invoked and the hook queue is
// empty.
func (f *ScenarioRunnerRunFunc) SetDefaultHook(hook func(*http.Client, map[string]interface{}) map[string]interface{}) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Run method of the parent MockScenarioRunner instance inovkes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ScenarioRunnerRunFunc) PushHook(hook func(*http.Client, map[string]interface{}) map[string]interface{}) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ScenarioRunnerRunFunc) SetDefaultReturn(r0 map[string]interface{}) {
	f.SetDefaultHook(func(*http.Client, map[string]interface{}) map[string]interface{} {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ScenarioRunnerRunFunc) PushReturn(r0 map[string]interface{}) {
	f.PushHook(func(*http.Client, map[string]interface{}) map[string]interface{} {
		return r0
	})
}

func (f *ScenarioRunnerRunFunc) nextHook() func(*http.Client, map[string]interface{}) map[string]interface{} {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ScenarioRunnerRunFunc) appendCall(r0 ScenarioRunnerRunFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ScenarioRunnerRunFuncCall objects
// describing the invocations of this function.
func (f *ScenarioRunnerRunFunc) History() []ScenarioRunnerRunFuncCall {
	f.mutex.Lock()
	history := make([]ScenarioRunnerRunFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ScenarioRunnerRunFuncCall is an object that describes an invocation of
// method Run on an instance of MockScenarioRunner.
type ScenarioRunnerRunFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 *http.Client
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 map[string]interface{}
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 map[string]interface{}
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ScenarioRunnerRunFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ScenarioRunnerRunFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}
