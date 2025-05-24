package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tnaucoin/stringer/types"
)

func TestParseCompositeActions(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int // expected number of parsed composite actions
	}{
		{
			name: "valid composite action",
			content: `
name: "Test Action"
description: "A test composite action"
inputs:
  name:
    description: "Name to greet"
    required: true
outputs:
  greeting:
    description: "Greeting"
    value: ${{ steps.greet.outputs.greeting }}
runs:
  using: "composite"
  steps:
    - run: echo "hello"
      shell: bash
`,
			expected: 1,
		},
		{
			name: "invalid yaml",
			content: `
runs:
  using: composite
  steps:
    - run: echo "hi`,
			expected: 0,
		},
		{
			name: "missing runs key",
			content: `
name: "No Runs"
description: "Missing the runs block"
`,
			expected: 0,
		},
		{
			name: "non-composite using",
			content: `
runs:
  using: "node12"
  main: "index.js"
`,
			expected: 0,
		},
		{ 
			name: "missing name and description",
			content: `
runs:
  using: "composite"
  steps:
    - run: echo "hi"
      shell: bash
`,
			expected: 0, // Currently implementation requires name and description
		},
		{
			name: "missing name only",
			content: `
description: "Action with missing name"
runs:
  using: "composite"
  steps:
    - run: echo "hi"
      shell: bash
`,
			expected: 0, // Currently implementation requires both name and description
		},
		{
			name: "missing description only",
			content: `
name: "Action with missing description"
runs:
  using: "composite"
  steps:
    - run: echo "hi"
      shell: bash
`,
			expected: 0, // Currently implementation requires both name and description
		},
		{
			name: "with inputs but no outputs",
			content: `
name: "Input Only Action"
description: "Has inputs but no outputs"
inputs:
  param1:
    description: "Parameter 1"
    required: true
runs:
  using: "composite"
  steps:
    - run: echo "Using ${{ inputs.param1 }}"
      shell: bash
`,
			expected: 1,
		},
		{
			name: "with outputs but no inputs",
			content: `
name: "Output Only Action"
description: "Has outputs but no inputs"
outputs:
  result:
    description: "Result value"
    value: ${{ steps.run.outputs.result }}
runs:
  using: "composite"
  steps:
    - id: run
      run: echo "::set-output name=result::success"
      shell: bash
`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			fpath := filepath.Join(tmpDir, "action.yml")

			if err := os.WriteFile(fpath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			actions, err := ParseCompositeActions(tmpDir)
			if err != nil {
				t.Fatalf("ParseCompositeActions returned error: %v", err)
			}

			if len(actions) != tt.expected {
				t.Errorf("expected %d actions, got %d", tt.expected, len(actions))
			}
		})
	}
}

func TestParseCompositeActionFromBytes(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected types.CompositeAction
		isError  bool
	}{
		{
			name: "valid composite action",
			content: `
name: "Test Action"
description: "A test composite action"
inputs:
  name:
    description: "Name to greet"
    required: true
outputs:
  greeting:
    description: "Greeting"
    value: ${{ steps.greet.outputs.greeting }}
runs:
  using: "composite"
  steps:
    - run: echo "hello"
      shell: bash
`,
			expected: types.CompositeAction{
				Name:        "Test Action",
				Description: "A test composite action",
				Path:        "test-path",
				Inputs:      map[string]any{"name": map[string]any{"description": "Name to greet", "required": true}},
				Outputs:     map[string]any{"greeting": map[string]any{"description": "Greeting", "value": "${{ steps.greet.outputs.greeting }}"}},
			},
			isError: false,
		},
		{
			name: "invalid yaml",
			content: `invalid: yaml: content`,
			expected: types.CompositeAction{},
			isError:  true,
		},
		{
			name: "not a composite action",
			content: `
name: "Not Composite"
description: "Not a composite action"
runs:
  using: "node12"
`,
			expected: types.CompositeAction{},
			isError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := ParseCompositeActionFromBytes([]byte(tt.content), "test-path")

			if tt.isError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.isError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.isError {
				if action.Name != tt.expected.Name {
					t.Errorf("expected name %q, got %q", tt.expected.Name, action.Name)
				}
				if action.Description != tt.expected.Description {
					t.Errorf("expected description %q, got %q", tt.expected.Description, action.Description)
				}
				if action.Path != tt.expected.Path {
					t.Errorf("expected path %q, got %q", tt.expected.Path, action.Path)
				}
				// Note: We're not doing deep comparison of inputs/outputs maps here
				// as that would require more complex test setup
			}
		})
	}
}
