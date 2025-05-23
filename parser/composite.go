package parser

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type CompositeAction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Inputs      map[string]any `json:"inputs"`
	Outputs     map[string]any `json:"outputs"`
	Path        string         `json:"-"`
}

// ParseCompositeActions scans a directory for composite GitHub Actions
func ParseCompositeActions(root string) ([]CompositeAction, error) {
	var actions []CompositeAction

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".yaml" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var raw map[string]any
			if err := yaml.Unmarshal(data, &raw); err != nil {
				return nil // skip invalid YAML
			}

			if raw["runs"] != nil {
				if runs, ok := raw["runs"].(map[string]any); ok && runs["using"] == "composite" {
					name, nameOk := raw["name"].(string)
					description, descOk := raw["description"].(string)

					if !nameOk || !descOk {
						return nil
					}

					action := CompositeAction{
						Name:        name,
						Description: description,
						Path:        path,
					}
					if v, ok := raw["Inputs"].(map[string]any); ok {
						action.Inputs = v
					}
					if v, ok := raw["Outputs"].(map[string]any); ok {
						action.Outputs = v
					}
					actions = append(actions, action)
				}
			}
		}
		return nil
	})

	return actions, err
}
