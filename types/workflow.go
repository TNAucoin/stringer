package types

type Workflow struct {
	Name string         `json:"name" yaml:"name"`
	On   map[string]any `json:"on" yaml:"on"`
	Jobs map[string]any `json:"jobs" yaml:"jobs"`
}

type Job struct {
	Name   string `json:"name" yaml:"name"`
	RunsOn string `json:"runs-on" yaml:"runs-on"`
	Steps  []Step `json:"steps" yaml:"steps"`
}

type Step struct {
	Name string            `json:"name" yaml:"name"`
	Uses string            `json:"uses,omitempty" yaml:"uses,omitempty"`
	Run  string            `json:"run,omitempty" yaml:"run,omitempty"`
	With map[string]string `json:"with,omitempty" yaml:"with,omitempty"`
}
