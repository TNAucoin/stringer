package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tnaucoin/stringer/types"
	"gopkg.in/yaml.v3"
)

func ParseCompositeActions(root string) ([]types.CompositeAction, error) {
	var actions []types.CompositeAction
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".yaml" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			action, err := ParseCompositeActionFromBytes(data, path)
			if err != nil {
				// TODO: log errors
			} else {
				actions = append(actions, action)
			}
		}
		return nil
	})
	return actions, err
}

// ParseCompositeActions scans a directory for composite GitHub Actions
func ParseCompositeActionFromBytes(data []byte, path string) (types.CompositeAction, error) {
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return types.CompositeAction{}, fmt.Errorf("invalid yaml file") // skip invalid YAML
	}

	runs, ok := raw["runs"].(map[string]any)
	if !ok || runs["using"] != "composite" {
		return types.CompositeAction{}, fmt.Errorf("not a composite action")
	}

	name := getString(raw["name"])
	description := getString(raw["description"])
	// TODO: name and desc, are optional on valid composite actions
	// should handle it, but for now treat it as invalid
	if name == "" || description == "" {
		return types.CompositeAction{}, fmt.Errorf("the composite action must have a name, and description")
	}

	action := types.CompositeAction{
		Name:        name,
		Description: description,
		Path:        path,
	}

	if v, ok := raw["inputs"].(map[string]any); ok {
		fmt.Println("Found input...")
		action.Inputs = v
	}

	if v, ok := raw["outputs"].(map[string]any); ok {
		fmt.Println("found output...")
		action.Outputs = v
	}

	return action, nil
}

func getString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
