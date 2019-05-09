package loader

import (
	"bytes"
	"fmt"
	tmpl "text/template"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type LoaderSuite struct{}

func (s *LoaderSuite) TestLoad(t sweet.T) {
	loaded, err := NewLoader().Load("./test-configs/basic.yaml")
	Expect(err).To(BeNil())

	for _, scenario := range loaded.Scenarios {
		for i, test := range scenario.Tests {
			Expect(testExec(test.Request.URL)).To(Equal(fmt.Sprintf("http://test.io/t%d", i+1)))
		}
	}

	Expect(loaded.Scenarios["s3"].Parallel).To(BeTrue())
	Expect(loaded.Scenarios["s2"].Tests[2].Enabled).To(BeFalse())
}

func (s *LoaderSuite) TestLoadIncludes(t sweet.T) {
	loaded, err := NewLoader().Load("./test-configs/include.yaml")
	Expect(err).To(BeNil())

	keys := []string{}
	for name := range loaded.Scenarios {
		keys = append(keys, name)
	}

	Expect(keys).To(ConsistOf("s1", "s2", "s3", "s4", "s5"))
}

func (s *LoaderSuite) TestLoadDuplicateScenario(t sweet.T) {
	_, err := NewLoader().Load("./test-configs/duplicate-scenario.yaml")
	Expect(err).To(MatchError("scenario 's1' defined more than once"))
}

func (s *LoaderSuite) TestLoadDuplicateTest(t sweet.T) {
	_, err := NewLoader().Load("./test-configs/duplicate-test.yaml")
	Expect(err).To(MatchError("test 's3/t1' defined more than once"))
}

func (s *LoaderSuite) TestLoadUnknownScenarioReference(t sweet.T) {
	_, err := NewLoader().Load("./test-configs/unknown-reference.yaml")
	Expect(err).To(MatchError("unknown scenario 's3' referenced in scenario 's1'"))
}

func (s *LoaderSuite) TestLoadDependencyCycles(t sweet.T) {
	_, err := NewLoader().Load("./test-configs/cycle.yaml")
	Expect(err).NotTo(BeNil()) // TODO
}

func (s *LoaderSuite) TestLoadInvalidSchemaConfig(t sweet.T) {
	_, err := NewLoader().Load("test-configs/invalid-config.yaml")
	Expect(err).To(MatchError("failed to validate config test-configs/invalid-config.yaml: Additional property tasks is not allowed"))
}

func (s *LoaderSuite) TestLoadInvalidSchemaInclude(t sweet.T) {
	_, err := NewLoader().Load("test-configs/invalid-include-parent.yaml")
	Expect(err).To(MatchError("failed to validate config test-configs/invalid-include.yaml: Additional property tasks is not allowed"))
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
