package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/efritz/api-test/config"
	"github.com/ghodss/yaml"
)

func Load(path string) (*config.Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to load config %s: %s",
			path,
			err.Error(),
		)
	}

	data, err = yaml.YAMLToJSON(data)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to load config %s: %s",
			path,
			err.Error(),
		)
	}

	if err := validate("schema/config.yaml", data); err != nil {
		return nil, fmt.Errorf(
			"failed to validate config %s: %s",
			path,
			err.Error(),
		)
	}

	payload := &jsonConfig{
		Tests: []string{},
	}

	if err := json.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	return payload.Translate()
}
