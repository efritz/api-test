package loader

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	tmpl "text/template"

	"github.com/aphistic/sweet"
	"github.com/efritz/api-test/config"
	. "github.com/onsi/gomega"
)

type LoaderSuite struct{}

func (s *LoaderSuite) TestLoad(t sweet.T) {
	loaded, err := Load("./test-configs/basic.yaml", nil)
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
	loaded, err := Load("./test-configs/include.yaml", nil)
	Expect(err).To(BeNil())

	keys := []string{}
	for name := range loaded.Scenarios {
		keys = append(keys, name)
	}

	Expect(keys).To(ConsistOf("s1", "s2", "s3", "s4", "s5"))
}

func (s *LoaderSuite) TestLoadOverrideFile(t sweet.T) {
	name := "api-test.override.yaml"
	content, _ := ioutil.ReadFile("./test-configs/override.yaml")
	ioutil.WriteFile(name, content, os.ModePerm)
	defer os.RemoveAll(name)

	loaded, err := Load("./test-configs/basic.yaml", nil)
	Expect(err).To(BeNil())
	Expect(loaded.Options.ForceSequential).To(BeTrue())
}

func (s *LoaderSuite) TestLoadOverrideCommandLine(t sweet.T) {
	loaded, err := Load("./test-configs/basic.yaml", &config.Override{
		Options: &config.Options{
			ForceSequential: true,
		},
	})

	Expect(err).To(BeNil())
	Expect(loaded.Options.ForceSequential).To(BeTrue())
}

func (s *LoaderSuite) TestLoadDuplicateScenario(t sweet.T) {
	_, err := Load("./test-configs/duplicate-scenario.yaml", nil)
	Expect(err).To(MatchError("scenario 's1' defined more than once"))
}

func (s *LoaderSuite) TestLoadDuplicateTest(t sweet.T) {
	_, err := Load("./test-configs/duplicate-test.yaml", nil)
	Expect(err).To(MatchError("test 's3/t1' defined more than once"))
}

func (s *LoaderSuite) TestLoadUnknownScenarioReference(t sweet.T) {
	_, err := Load("./test-configs/unknown-reference.yaml", nil)
	Expect(err).To(MatchError("unknown scenario 's3' referenced in scenario 's1'"))
}

func (s *LoaderSuite) TestLoadIncludeCycles(t sweet.T) {
	// Must terminate
	_, err := Load("./test-configs/cyclic-includes-parent.yaml", nil)
	Expect(err).To(BeNil())
}

func (s *LoaderSuite) TestLoadDependencyCycles(t sweet.T) {
	_, err := Load("./test-configs/cyclic-dependencies.yaml", nil)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(HavePrefix("scenario dependencies are cyclic"))
}

func (s *LoaderSuite) TestLoadInvalidSchemaConfig(t sweet.T) {
	_, err := Load("test-configs/invalid-config.yaml", nil)
	Expect(err).To(MatchError("failed to validate input test-configs/invalid-config.yaml: Additional property foobar is not allowed"))
}

func (s *LoaderSuite) TestLoadInvalidSchemaInclude(t sweet.T) {
	_, err := Load("test-configs/invalid-include-parent.yaml", nil)
	Expect(err).To(MatchError("failed to validate input test-configs/invalid-include.yaml: Additional property foobar is not allowed"))
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
