package jsonconfig

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ScenarioSuite struct{}

func (s *ScenarioSuite) TestTranslate(t sweet.T) {
	enabled := true

	scenario := &Scenario{
		Name:         "scenario",
		Enabled:      &enabled,
		Dependencies: []byte(`"dep"`),
		Parallel:     true,
		Tests: []*Test{
			&Test{Request: &Request{URI: "/t1"}},
			&Test{Request: &Request{URI: "/t2"}},
			&Test{Request: &Request{URI: "/t3"}},
		},
	}

	translated, err := scenario.Translate(&GlobalRequest{
		BaseURL: "http://test.io",
	})

	Expect(err).To(BeNil())
	Expect(translated.Name).To(Equal("scenario"))
	Expect(translated.Enabled).To(BeTrue())
	Expect(translated.Dependencies).To(ConsistOf("dep"))
	Expect(translated.Parallel).To(BeTrue())
	Expect(translated.Tests).To(HaveLen(3))
	Expect(testExec(translated.Tests[0].Request.URL)).To(Equal("http://test.io/t1"))
	Expect(testExec(translated.Tests[1].Request.URL)).To(Equal("http://test.io/t2"))
	Expect(testExec(translated.Tests[2].Request.URL)).To(Equal("http://test.io/t3"))
}

func (s *ScenarioSuite) TestTranslateStringLists(t sweet.T) {
	scenario := &Scenario{
		Dependencies: []byte(`["dep1", "dep2", "dep3"]`),
	}

	translated, err := scenario.Translate(&GlobalRequest{
		BaseURL: "http://test.io",
	})

	Expect(err).To(BeNil())
	Expect(translated.Dependencies).To(ConsistOf("dep1", "dep2", "dep3"))
}

func (s *ScenarioSuite) TestTranslateDefaultEnabled(t sweet.T) {
	scenario := &Scenario{Enabled: nil}
	translated, err := scenario.Translate(nil)
	Expect(err).To(BeNil())
	Expect(translated.Enabled).To(BeTrue())
}
