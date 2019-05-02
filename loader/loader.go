package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/loader/jsonconfig"
	"github.com/efritz/api-test/loader/schema"
	"github.com/efritz/api-test/loader/util"
	"github.com/ghodss/yaml"
)

type Loader struct {
	loadedConfigs map[string]map[string]config.Scenario
}

func NewLoader() *Loader {
	return &Loader{
		loadedConfigs: map[string]map[string]config.Scenario{},
	}
}

func (l *Loader) Load(path string) (*config.Config, error) {
	path = normalizePath(path, "")

	data, err := readPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config %s: %s", path, err.Error())
	}

	payload := &jsonconfig.MainConfig{
		Options: &jsonconfig.Options{},
	}

	if err := unmarshal(path, data, "schema/config.yaml", &payload); err != nil {
		return nil, err
	}

	scenarios, err := payload.Translate(payload.GlobalRequest)
	if err != nil {
		return nil, err
	}

	scenarioMap := map[string][]*config.Scenario{
		path: scenarios,
	}

	err = l.loadIncludes(
		path,
		payload.Includes,
		payload.GlobalRequest,
		scenarioMap,
	)

	if err != nil {
		return nil, err
	}

	flattenedMap, err := validateScenarios(scenarioMap)
	if err != nil {
		return nil, err
	}

	options, err := payload.Options.Translate()
	if err != nil {
		return nil, err
	}

	config := &config.Config{
		Options:   options,
		Scenarios: flattenedMap,
	}

	return config, nil
}

func (l *Loader) loadIncludes(
	parent string,
	rawIncludes json.RawMessage,
	globalRequest *jsonconfig.GlobalRequest,
	scenarioMap map[string][]*config.Scenario,
) error {
	paths, err := util.UnmarshalStringList(rawIncludes)
	if err != nil {
		return err
	}

	for _, path := range paths {
		if _, ok := scenarioMap[path]; ok {
			continue
		}

		if err := l.loadInclude(
			normalizePath(path, parent),
			globalRequest,
			scenarioMap,
		); err != nil {
			return err
		}
	}

	return nil
}

func (l *Loader) loadInclude(
	path string,
	globalRequest *jsonconfig.GlobalRequest,
	scenarioMap map[string][]*config.Scenario,
) error {
	data, err := readPath(path)
	if err != nil {
		return fmt.Errorf("failed to load config %s: %s", path, err.Error())
	}

	payload := &jsonconfig.BaseConfig{}
	if err := unmarshal(path, data, "schema/include.yaml", &payload); err != nil {
		return err
	}

	scenarios, err := payload.Translate(globalRequest)
	if err != nil {
		return err
	}

	scenarioMap[path] = scenarios

	return l.loadIncludes(
		path,
		payload.Includes,
		globalRequest,
		scenarioMap,
	)
}

//
// Helpers

func normalizePath(path, source string) string {
	if source == "" || filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(filepath.Dir(source), path)
}

func readPath(path string) ([]byte, error) {
	rawData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return yaml.YAMLToJSON(rawData)
}

func unmarshal(path string, data []byte, schemaName string, payload interface{}) error {
	if err := schema.Validate(schemaName, data); err != nil {
		return fmt.Errorf(
			"failed to validate config %s: %s",
			path,
			err.Error(),
		)
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	return nil
}
