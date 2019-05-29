package jsonconfig

import (
	"github.com/aphistic/sweet"
	"github.com/efritz/api-test/config"
	. "github.com/onsi/gomega"
)

type ConfigSuite struct{}

func (s *ConfigSuite) TestTranslate(t sweet.T) {
	config := &Config{
		Scenarios: []*Scenario{
			&Scenario{Name: "foo", Tests: []*Test{&Test{Request: &Request{URI: "/t1"}}}},
			&Scenario{Name: "bar", Tests: []*Test{&Test{Request: &Request{URI: "/t2"}}}},
			&Scenario{Name: "baz", Tests: []*Test{&Test{Request: &Request{URI: "/t3"}}}},
		},
	}

	scenarios, err := config.Translate(&GlobalRequest{
		BaseURL: "http://test.io",
	})

	Expect(err).To(BeNil())
	Expect(scenarios).To(HaveLen(3))
	Expect(scenarios[0].Name).To(Equal("foo"))
	Expect(scenarios[1].Name).To(Equal("bar"))
	Expect(scenarios[2].Name).To(Equal("baz"))
	Expect(scenarios[0].Tests).To(HaveLen(1))
	Expect(scenarios[1].Tests).To(HaveLen(1))
	Expect(scenarios[2].Tests).To(HaveLen(1))
	Expect(testExec(scenarios[0].Tests[0].Request.URL)).To(Equal("http://test.io/t1"))
	Expect(testExec(scenarios[1].Tests[0].Request.URL)).To(Equal("http://test.io/t2"))
	Expect(testExec(scenarios[2].Tests[0].Request.URL)).To(Equal("http://test.io/t3"))
}

func (s *ConfigSuite) TestTranslateOptions(t sweet.T) {
	options := &Options{
		ForceSequential: true,
	}

	translated, err := options.Translate()
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Options{
		ForceSequential: true,
	}))
}
