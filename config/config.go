package config

import (
	"fmt"
	"strings"
)

type (
	Config struct {
		Scenarios map[string]*Scenario
		Options   *Options
	}

	Options struct {
		ForceSequential bool
		MaxParallelism  int
	}

	Override struct {
		Options *Options
	}
)

func (c *Config) ApplyOverride(override *Override) {
	if override != nil {
		if override.Options.ForceSequential {
			c.Options.ForceSequential = true
		}

		if override.Options.MaxParallelism > 0 {
			c.Options.MaxParallelism = override.Options.MaxParallelism
		}
	}
}

func (c *Config) EnableTests(tests []string) error {
	if len(tests) == 0 {
		return nil
	}

	enabled, err := getEnabled(tests)
	if err != nil {
		return err
	}

	if err := c.validateTests(enabled); err != nil {
		return err
	}

	dependencies := c.getDependencies(enabled)

	for name, scenario := range c.Scenarios {
		if _, ok := enabled[name]; ok {
			scenario.Enabled = true
			continue
		}

		if _, ok := dependencies[name]; ok {
			scenario.Enabled = true
			continue
		}

		scenario.Disabled = true
	}

	for scenarioName, testNames := range enabled {
		if _, ok := dependencies[scenarioName]; !ok && len(testNames) > 0 {
			c.Scenarios[scenarioName].EnableTests(testNames)
		}
	}

	return nil
}

func (c *Config) validateTests(enabled map[string][]string) error {
	for scenarioName, testNames := range enabled {
		scenario, ok := c.Scenarios[scenarioName]
		if !ok {
			return fmt.Errorf("unknown scenario '%s'", scenarioName)
		}

		for _, testName := range testNames {
			if !scenario.ContainsTest(testName) {
				return fmt.Errorf("unknown test '%s/%s'", scenarioName, testName)
			}
		}
	}

	return nil
}

func (c *Config) getDependencies(enabled map[string][]string) map[string]struct{} {
	dependencies := map[string]struct{}{}
	for scenarioName := range enabled {
		for dependency := range c.Scenarios[scenarioName].AllDependencies(c) {
			dependencies[dependency] = struct{}{}
		}
	}

	return dependencies
}

//
// Helperss

func getEnabled(tests []string) (map[string][]string, error) {
	enabled := map[string][]string{}
	for _, name := range tests {
		parts := strings.Split(name, "/")

		if len(parts) == 1 {
			enabled[parts[0]] = append(enabled[parts[0]])
		} else if len(parts) == 2 {
			enabled[parts[0]] = append(enabled[parts[0]], parts[1])
		} else {
			return nil, fmt.Errorf("illegal test name '%s'", name)
		}
	}

	return enabled, nil
}
