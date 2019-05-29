package jsonconfig

import (
	"bytes"
	"fmt"
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

		s.AddSuite(&AuthSuite{})
		s.AddSuite(&ConfigSuite{})
		s.AddSuite(&ExtractorSuite{})
		s.AddSuite(&OverrideSuite{})
		s.AddSuite(&RequestSuite{})
		s.AddSuite(&ResponseSuite{})
		s.AddSuite(&ScenarioSuite{})
		s.AddSuite(&TestSuite{})
	})
}

//
// Helpers

func testExec(template *tmpl.Template) string {
	buffer := bytes.NewBuffer(nil)
	if err := template.Execute(buffer, nil); err != nil {
		panic(fmt.Sprintf("failed to execute template (%s)", err.Error()))
	}

	return buffer.String()
}
