package loader

import (
	"fmt"

	"github.com/efritz/api-test/config"
	"github.com/stevenle/topsort"
)

func validateScenarios(scenarioMap map[string][]*config.Scenario) (map[string]*config.Scenario, error) {
	flattened := map[string]*config.Scenario{}
	for _, scenarios := range scenarioMap {
		for _, scenario := range scenarios {
			if _, ok := flattened[scenario.Name]; ok {
				return nil, fmt.Errorf("scenario '%s' defined more than once", scenario.Name)
			}

			flattened[scenario.Name] = scenario
		}
	}

	for name, scenario := range flattened {
		for _, dependency := range scenario.Dependencies {
			if _, ok := flattened[dependency]; !ok {
				return nil, fmt.Errorf("unknown scenario '%s' referenced in scenario '%s'", dependency, name)
			}
		}
	}

	if err := checkCycles(flattened); err != nil {
		return nil, err
	}

	return flattened, nil
}

func checkCycles(scenarios map[string]*config.Scenario) error {
	dependencyGraph := topsort.NewGraph()

	dependencyGraph.AddNode("$")

	for name := range scenarios {
		dependencyGraph.AddNode(name)
		dependencyGraph.AddEdge("$", name)
	}

	for name, scenario := range scenarios {
		for _, depencency := range scenario.Dependencies {
			dependencyGraph.AddEdge(name, depencency)
		}
	}

	if _, err := dependencyGraph.TopSort("$"); err != nil {
		// Error messages starts with "Cycle error: "
		// return nil, fmt.Errorf("failed to extend cyclic config (%s)", err.Error()[13:])
		return err
	}

	return nil
}
