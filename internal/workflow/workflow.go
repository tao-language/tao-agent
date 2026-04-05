package workflow

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Workflow struct {
	Version     string  `yaml:"version"`
	Name        string  `yaml:"name"`
	Description string  `yaml:"description,omitempty"`
	Steps       []*Step `yaml:"steps"`
}

// Load reads a workflow from a YAML file.
func Load(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var w Workflow
	if err := yaml.Unmarshal(data, &w); err != nil {
		return nil, err
	}

	return &w, nil
}
