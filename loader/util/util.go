package util

import "encoding/json"

func UnmarshalStringList(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	single := ""
	if err := json.Unmarshal(raw, &single); err == nil {
		return []string{single}, nil
	}

	multiple := []string{}
	if err := json.Unmarshal(raw, &multiple); err != nil {
		return nil, err
	}

	return multiple, nil
}
