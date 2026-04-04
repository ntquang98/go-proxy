package rules

import (
	"encoding/json"
	"os"
)

func LoadRules(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rules []Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}
