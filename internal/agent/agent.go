package agent

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ModelConfig struct {
	Provider string `yaml:"provider"`
	Name     string `yaml:"name"`
	URL      string `yaml:"url,omitempty"`
}

type Agent struct {
	Name         string      `yaml:"name"`
	Model        ModelConfig `yaml:"model"`
	SystemPrompt string      `yaml:"system-prompt"`
}

// Load reads an agent definition from a YAML file.
func Load(path string) (*Agent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var a Agent
	if err := yaml.Unmarshal(data, &a); err != nil {
		return nil, err
	}

	return &a, nil
}

// GetDefault returns a sensible default agent if no file is found.
func GetDefault() *Agent {
	return &Agent{
		Name: "default",
		Model: ModelConfig{
			Provider: "ollama",
			Name:     "qwen2.5:latest",
			URL:      "http://localhost:11434",
		},
		SystemPrompt: "You are a helpful assistant.",
	}
}
