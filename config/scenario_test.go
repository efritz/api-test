package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ScenarioSuite struct{}

func (s *ScenarioSuite) TestAllDependencies(t sweet.T) {
	config := &Config{
		Scenarios: map[string]*Scenario{
			"s1": &Scenario{Dependencies: []string{}},
			"s2": &Scenario{Dependencies: []string{"s1"}},
			"s3": &Scenario{Dependencies: []string{"s1"}},
			"s4": &Scenario{Dependencies: []string{"s3"}},
		},
	}

	dependencies := config.Scenarios["s4"].AllDependencies(config)

	Expect(dependencies).To(Equal(map[string]struct{}{
		"s1": struct{}{},
		"s3": struct{}{},
	}))
}

func (s *ScenarioSuite) TestContainsTest(t sweet.T) {
	scenario := &Scenario{
		Tests: []*Test{
			&Test{Name: "foo"},
			&Test{Name: "bar"},
			&Test{Name: "baz"},
		},
	}

	Expect(scenario.ContainsTest("foo")).To(BeTrue())
}

func (s *ScenarioSuite) TestContainsTestMissing(t sweet.T) {
	scenario := &Scenario{
		Tests: []*Test{
			&Test{Name: "foo"},
			&Test{Name: "bar"},
			&Test{Name: "baz"},
		},
	}

	Expect(scenario.ContainsTest("bonk")).To(BeFalse())
}

func (s *ScenarioSuite) TestEnableTestsSequential(t sweet.T) {
	scenario := &Scenario{
		Tests: []*Test{
			&Test{Name: "t1", Enabled: true},
			&Test{Name: "t2", Enabled: false},
			&Test{Name: "t3", Enabled: false},
			&Test{Name: "t4", Enabled: true},
			&Test{Name: "t5", Enabled: true},
		},
	}

	scenario.EnableTests([]string{"t1", "t3"})
	Expect(scenario.Tests[0].Enabled).To(BeTrue())
	Expect(scenario.Tests[1].Enabled).To(BeFalse())
	Expect(scenario.Tests[2].Enabled).To(BeTrue())
	Expect(scenario.Tests[0].Disabled).To(BeFalse())
	Expect(scenario.Tests[1].Disabled).To(BeFalse())
	Expect(scenario.Tests[2].Disabled).To(BeFalse())
	Expect(scenario.Tests[3].Disabled).To(BeTrue())
	Expect(scenario.Tests[4].Disabled).To(BeTrue())
}

func (s *ScenarioSuite) TestEnableTestsParallel(t sweet.T) {
	scenario := &Scenario{
		Parallel: true,
		Tests: []*Test{
			&Test{Name: "t1", Enabled: true},
			&Test{Name: "t2", Enabled: false},
			&Test{Name: "t3", Enabled: false},
			&Test{Name: "t4", Enabled: true},
			&Test{Name: "t5", Enabled: true},
		},
	}

	scenario.EnableTests([]string{"t1", "t3"})
	Expect(scenario.Tests[0].Enabled).To(BeTrue())
	Expect(scenario.Tests[1].Enabled).To(BeFalse())
	Expect(scenario.Tests[2].Enabled).To(BeTrue())
	Expect(scenario.Tests[3].Enabled).To(BeTrue())
	Expect(scenario.Tests[4].Enabled).To(BeTrue())
	Expect(scenario.Tests[0].Disabled).To(BeFalse())
	Expect(scenario.Tests[1].Disabled).To(BeTrue())
	Expect(scenario.Tests[2].Disabled).To(BeFalse())
	Expect(scenario.Tests[3].Disabled).To(BeTrue())
	Expect(scenario.Tests[4].Disabled).To(BeTrue())
}
