package config

type Scenario struct {
	Name         string
	Enabled      bool
	Disabled     bool
	Dependencies []string
	Parallel     bool
	Tests        []*Test
}

func (s *Scenario) AllDependencies(config *Config) map[string]struct{} {
	dependencies := map[string]struct{}{}
	for _, dependency := range s.Dependencies {
		dependencies[dependency] = struct{}{}

		for dependency := range config.Scenarios[dependency].AllDependencies(config) {
			dependencies[dependency] = struct{}{}
		}
	}

	return dependencies
}

func (s *Scenario) ContainsTest(testName string) bool {
	for _, test := range s.Tests {
		if test.Name == testName {
			return true
		}
	}

	return false
}

func (s *Scenario) EnableTests(enabled []string) {
	if s.Parallel {
		s.enableParallel(enabled)
	} else {
		s.enableSequential(enabled)
	}
}

func (s *Scenario) enableSequential(enabled []string) {
	found := false
	for i := len(s.Tests) - 1; i >= 0; i-- {
		if isEnabled(s.Tests[i].Name, enabled) {
			found = true
			s.Tests[i].Enabled = true
			continue
		}

		if !found {
			s.Tests[i].Disabled = true
		}
	}
}

func (s *Scenario) enableParallel(enabled []string) {
	for _, test := range s.Tests {
		if isEnabled(test.Name, enabled) {
			test.Enabled = true
		} else {
			test.Disabled = true
		}
	}
}

func isEnabled(testName string, enabled []string) bool {
	for _, name := range enabled {
		if testName == name {
			return true
		}
	}

	return false
}
