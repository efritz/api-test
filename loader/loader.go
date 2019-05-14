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
	loadedConfigs map[string]struct{}
}

func NewLoader() *Loader {
	return &Loader{
		loadedConfigs: map[string]struct{}{},
	}
}

func (l *Loader) Load(path string) (*config.Config, error) {
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

	scenarios, err = l.loadIncludes(
		path,
		payload.Includes,
		payload.GlobalRequest,
		scenarios,
	)

	if err != nil {
		return nil, err
	}

	flattenedMap, err := validateScenarios(scenarios)
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

func (l *Loader) LoadOverride(path string) (*config.Override, error) {
	data, err := readPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load override %s: %s", path, err.Error())
	}

	payload := &jsonconfig.Override{
		Options: &jsonconfig.Options{},
	}

	if err := unmarshal(path, data, "schema/override.yaml", &payload); err != nil {
		return nil, err
	}

	return payload.Translate()
}

func (l *Loader) loadIncludes(
	parent string,
	rawIncludes json.RawMessage,
	globalRequest *jsonconfig.GlobalRequest,
	scenarios []*config.Scenario,
) ([]*config.Scenario, error) {
	paths, err := util.UnmarshalStringList(rawIncludes)
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		scenarios, err = l.loadInclude(
			normalizePath(path, parent),
			globalRequest,
			scenarios,
		)

		if err != nil {
			return nil, err
		}
	}

	return scenarios, nil
}

func (l *Loader) loadInclude(
	path string,
	globalRequest *jsonconfig.GlobalRequest,
	scenarios []*config.Scenario,
) ([]*config.Scenario, error) {
	if _, ok := l.loadedConfigs[path]; ok {
		return scenarios, nil
	}

	l.loadedConfigs[path] = struct{}{}

	data, err := readPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config %s: %s", path, err.Error())
	}

	payload := &jsonconfig.BaseConfig{}
	if err := unmarshal(path, data, "schema/include.yaml", &payload); err != nil {
		return nil, err
	}

	translated, err := payload.Translate(globalRequest)
	if err != nil {
		return nil, err
	}

	return l.loadIncludes(
		path,
		payload.Includes,
		globalRequest,
		append(scenarios, translated...),
	)
}

//
// Command Line

func Load(path string, override *config.Override) (*config.Config, error) {
	overridePath, err := GetOverridePath()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to determine override path: %s",
			err.Error(),
		)
	}

	loader := NewLoader()

	config, err := loader.Load(normalizePath(path, ""))
	if err != nil {
		return nil, err
	}

	if overridePath != "" {
		localOverride, err := loader.LoadOverride(overridePath)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to apply override: %s",
				err.Error(),
			)
		}

		config.ApplyOverride(localOverride)
	}

	if override != nil {
		config.ApplyOverride(override)
	}

	return config, nil
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
			"failed to validate input %s: %s",
			path,
			err.Error(),
		)
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	return nil
}
