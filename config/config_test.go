package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ConfigSuite struct{}

func (s *ConfigSuite) TestApplyOverride(t sweet.T) {
	config := &Config{
		Options: &Options{},
	}

	config.ApplyOverride(&Override{
		Options: &Options{
			ForceSequential: true,
			MaxParallelism:  10,
		},
	})

	Expect(config.Options.ForceSequential).To(BeTrue())
	Expect(config.Options.MaxParallelism).To(Equal(10))
}

func (s *ConfigSuite) TestApplyOverrideNoValues(t sweet.T) {
	config := &Config{
		Options: &Options{
			ForceSequential: true,
			MaxParallelism:  10,
		},
	}

	config.ApplyOverride(&Override{
		Options: &Options{
			ForceSequential: false,
			MaxParallelism:  0,
		},
	})

	Expect(config.Options.ForceSequential).To(BeTrue())
	Expect(config.Options.MaxParallelism).To(Equal(10))
}

func (s *ConfigSuite) TestApplyOverrideNil(t sweet.T) {
	config := &Config{
		Options: &Options{
			ForceSequential: true,
			MaxParallelism:  10,
		},
	}

	config.ApplyOverride(nil)
	Expect(config.Options.ForceSequential).To(BeTrue())
	Expect(config.Options.MaxParallelism).To(Equal(10))
}

func (s *ConfigSuite) TestEnableTests(t sweet.T) {
	config := &Config{
		Scenarios: map[string]*Scenario{
			"s1": &Scenario{
				Tests: []*Test{
					&Test{Name: "t1"},
					&Test{Name: "t2"},
					&Test{Name: "t3"},
				},
			},
			"s2": &Scenario{
				Dependencies: []string{"s1"},
				Tests: []*Test{
					&Test{Name: "t1"},
					&Test{Name: "t2"},
					&Test{Name: "t3"},
				},
			},
			"s3": &Scenario{},
		},
	}

	err := config.EnableTests([]string{"s2/t2"})
	Expect(err).To(BeNil())
	Expect(config.Scenarios["s1"].Disabled).To(BeFalse())
	Expect(config.Scenarios["s2"].Disabled).To(BeFalse())
	Expect(config.Scenarios["s3"].Disabled).To(BeTrue())
	Expect(config.Scenarios["s1"].Tests[0].Disabled).To(BeFalse())
	Expect(config.Scenarios["s1"].Tests[1].Disabled).To(BeFalse())
	Expect(config.Scenarios["s1"].Tests[2].Disabled).To(BeFalse())
	Expect(config.Scenarios["s2"].Tests[0].Disabled).To(BeFalse())
	Expect(config.Scenarios["s2"].Tests[1].Disabled).To(BeFalse())
	Expect(config.Scenarios["s2"].Tests[2].Disabled).To(BeTrue())
}

func (s *ConfigSuite) TestEnableTestsEnablesExplicit(t sweet.T) {
	config := &Config{
		Scenarios: map[string]*Scenario{
			"s1": &Scenario{
				Enabled:  false,
				Parallel: true,
				Tests: []*Test{
					&Test{Name: "t1", Enabled: false},
					&Test{Name: "t2", Enabled: false},
					&Test{Name: "t3", Enabled: false},
				},
			},
		},
	}

	err := config.EnableTests([]string{"s1/t2"})
	Expect(err).To(BeNil())
	Expect(config.Scenarios["s1"].Enabled).To(BeTrue())
	Expect(config.Scenarios["s1"].Disabled).To(BeFalse())
	Expect(config.Scenarios["s1"].Tests[0].Enabled).To(BeFalse())
	Expect(config.Scenarios["s1"].Tests[1].Enabled).To(BeTrue())
	Expect(config.Scenarios["s1"].Tests[2].Enabled).To(BeFalse())
	Expect(config.Scenarios["s1"].Tests[0].Disabled).To(BeTrue())
	Expect(config.Scenarios["s1"].Tests[1].Disabled).To(BeFalse())
	Expect(config.Scenarios["s1"].Tests[2].Disabled).To(BeTrue())
}

func (s *ConfigSuite) TestEnableTestsNoExplicitList(t sweet.T) {
	config := &Config{
		Scenarios: map[string]*Scenario{
			"s1": &Scenario{
				Enabled: true,
				Tests: []*Test{
					&Test{Name: "t1", Enabled: true},
					&Test{Name: "t2", Enabled: false},
					&Test{Name: "t3", Enabled: true},
				},
			},
			"s2": &Scenario{
				Enabled: false,
				Tests: []*Test{
					&Test{Name: "t1", Enabled: true},
					&Test{Name: "t2", Enabled: false},
					&Test{Name: "t3", Enabled: true},
				},
			},
		},
	}

	err := config.EnableTests(nil)
	Expect(err).To(BeNil())

	for _, scenario := range config.Scenarios {
		Expect(scenario.Disabled).To(BeFalse())

		for _, test := range scenario.Tests {
			Expect(test.Disabled).To(BeFalse())
		}
	}
}

func (s *ConfigSuite) TestEnableTestsUnknownScenario(t sweet.T) {
	config := &Config{
		Scenarios: map[string]*Scenario{
			"s1": &Scenario{},
			"s2": &Scenario{},
			"s3": &Scenario{},
		},
	}

	err := config.EnableTests([]string{"s2", "s4"})
	Expect(err).To(MatchError("unknown scenario 's4'"))
}

func (s *ConfigSuite) TestEnableTestsUnknownTest(t sweet.T) {
	config := &Config{
		Scenarios: map[string]*Scenario{
			"s1": &Scenario{Tests: []*Test{
				&Test{Name: "t1"},
			}},
		},
	}

	err := config.EnableTests([]string{"s1/t1", "s1/t2"})
	Expect(err).To(MatchError("unknown test 's1/t2'"))
}

func (s *ConfigSuite) TestGetDependencies(t sweet.T) {
	config := &Config{
		Scenarios: map[string]*Scenario{
			"s1": &Scenario{Dependencies: []string{}},
			"s2": &Scenario{Dependencies: []string{"s1"}},
			"s3": &Scenario{Dependencies: []string{"s2"}},
			"s4": &Scenario{Dependencies: []string{"s3"}},
			"s5": &Scenario{Dependencies: []string{}},
			"s6": &Scenario{Dependencies: []string{"s5"}},
			"s7": &Scenario{Dependencies: []string{"s6"}},
		},
	}

	dependencies := config.getDependencies(map[string][]string{
		"s2": []string{},
		"s3": []string{},
		"s7": []string{},
	})

	Expect(dependencies).To(Equal(map[string]struct{}{
		"s1": struct{}{},
		"s2": struct{}{},
		"s5": struct{}{},
		"s6": struct{}{},
	}))
}

func (s *ConfigSuite) TestGetEnabled(t sweet.T) {
	enabled, err := getEnabled([]string{"foo", "bar", "foo/a", "baz/b", "baz/c"})
	Expect(err).To(BeNil())
	Expect(enabled).To(Equal(map[string][]string{
		"foo": []string{"a"},
		"bar": nil,
		"baz": []string{"b", "c"},
	}))
}

func (s *ConfigSuite) TestGetEnabledIllegalName(t sweet.T) {
	_, err := getEnabled([]string{"foo/bar/baz"})
	Expect(err).To(MatchError("illegal test name 'foo/bar/baz'"))
}
