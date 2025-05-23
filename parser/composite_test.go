package parser

import (
	"os"
	"path/filepath"
	"testing"
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
		{ // TODO: This is actually a valid case, but currently we wont handle it, name and desc are optional on composite actions
			name: "missing name and description",
			content: `
runs:
  using: "composite"
  steps:
    - run: echo "hi"
      shell: bash
`,
			expected: 0,
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
