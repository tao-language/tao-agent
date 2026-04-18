package workflow

import (
	"tao-agent/internal/types"
)

type Step struct {
	Name       string                 `yaml:"name,omitempty"`
	ID         string                 `yaml:"id,omitempty"`
	If         string                 `yaml:"if,omitempty"`
	Prompt     string                 `yaml:"prompt,omitempty"`
	Agent      string                 `yaml:"agent,omitempty"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
	Tools      []string               `yaml:"tools,omitempty"`
	Outputs    *types.Definition      `yaml:"outputs,omitempty"`
	Ask        string                 `yaml:"ask,omitempty"`
	Choices    map[string]interface{} `yaml:"choices,omitempty"`
	Tool       string                 `yaml:"tool,omitempty"`
	Inputs     map[string]interface{} `yaml:"inputs,omitempty"`
	Print      string                 `yaml:"print,omitempty"`
	While       string                 `yaml:"while,omitempty"` // continue WHILE condition is true
	MaxIterations *int                  `yaml:"max-iterations,omitempty"` // optional upper bound (nil = unbounded)
	Steps       []*Step                `yaml:"steps,omitempty"`     // body of while loop
	Match      string                 `yaml:"match,omitempty"`
	Cases      map[string][]*Step     `yaml:"cases,omitempty"`
	Use        string                 `yaml:"use,omitempty"`
}
