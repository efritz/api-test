package loader

import (
	"fmt"
	"os"
)

var (
	defaultConfigPaths = []string{
		"api-test.yaml",
		"api-test.yml",
	}

	overridePaths = []string{
		"api-test.override.yaml",
		"api-test.override.yml",
	}
)

func GetConfigPath(path string) (string, error) {
	if path != "" {
		return path, nil
	}

	path, err := findFirstToExist(defaultConfigPaths)
	if err != nil {
		return "", err
	}

	if path == "" {
		return "", fmt.Errorf("could not infer config file")
	}

	return path, err
}

func GetOverridePath() (string, error) {
	return findFirstToExist(overridePaths)
}

func findFirstToExist(paths []string) (string, error) {
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return "", err
		}

		if info.IsDir() {
			return "", fmt.Errorf("%s exists but is not a file", path)
		}

		return path, nil
	}

	return "", nil
}
