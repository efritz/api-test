package runner

//go:generate go-mockgen -f  github.com/efritz/api-test/runner -i ScenarioRunner -o scenario_runner_mock_test.go
//go:generate go-mockgen -f  github.com/efritz/api-test/runner -i ScenarioRunnerFactory -o scenario_runner_factory_mock_test.go

import (
	"regexp"
	"testing"
	tmpl "text/template"

	"github.com/aphistic/sweet"
	junit "github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&ReportSuite{})
		s.AddSuite(&RequestSuite{})
		s.AddSuite(&ResponseSuite{})
		s.AddSuite(&ResultSuite{})
		s.AddSuite(&RunnerSuite{})
		s.AddSuite(&ScenarioRunnerSuite{})
	})
}

//
// Helpers

func testTemplate(pattern string) *tmpl.Template {
	return tmpl.Must(tmpl.New("").Parse(pattern))
}

func testPattern(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}
