package runner

import (
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&DisplaySuite{})
		s.AddSuite(&ReportSuite{})
		s.AddSuite(&RequestSuite{})
		s.AddSuite(&ResponseSuite{})
		s.AddSuite(&ResultSuite{})
		s.AddSuite(&RunnerSuite{})
		s.AddSuite(&ScenarioSuite{})
	})
}
