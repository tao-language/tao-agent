package mcpmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"tao-agent/pkg/mcp"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Version     string   `yaml:"version"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Command     string   `yaml:"command,omitempty"`
	Arguments   []string `yaml:"arguments,omitempty"`
	URL         string   `yaml:"url,omitempty"`
}

type Manager struct {
	Servers map[string]*mcp.Client
}

func NewManager() *Manager {
	return &Manager{
		Servers: make(map[string]*mcp.Client),
	}
}

func (m *Manager) LoadServer(name string) (*mcp.Client, error) {
	if client, ok := m.Servers[name]; ok {
		return client, nil
	}

	path := filepath.Join("tools", name+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ServerConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	client, err := mcp.NewClient(config.Command, config.Arguments)
	if err != nil {
		return nil, err
	}

	// Initialize MCP session (send 'initialize' request)
	// For simplicity in MCP Lite, we'll skip the full lifecycle for now

	m.Servers[name] = client
	return client, nil
}

func (m *Manager) CallTool(serverName, toolName string, inputs interface{}) (string, error) {
	client, err := m.LoadServer(serverName)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"name":      toolName,
		"arguments": inputs,
	}

	resp, err := client.Call("tools/call", params)
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", fmt.Errorf("MCP error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	// Parse CallToolResult
	var res mcp.CallToolResult
	data, _ := json.Marshal(resp.Result)
	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}

	if len(res.Content) > 0 {
		return res.Content[0].Text, nil
	}

	return "", nil
}
