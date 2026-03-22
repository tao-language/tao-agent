package model

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

type ProviderType string

const (
	ProviderOllama   ProviderType = "ollama"
	ProviderLMStudio ProviderType = "lm-studio"
)

type ModelProvider interface {
	Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)
	IsRunning() bool
	Start() error
	GetModels() ([]string, error)
}

type ProviderManager struct {
	Ollama   *OllamaClient
	LMStudio *OpenAICompatibleClient
}

func NewProviderManager() *ProviderManager {
	return &ProviderManager{
		Ollama:   NewOllamaClient("http://localhost:11434"),
		LMStudio: NewOpenAICompatibleClient("http://localhost:1234/v1"),
	}
}

// Check if a service is listening on a port
func isPortOpen(url string) bool {
	client := http.Client{Timeout: 500 * time.Millisecond}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (pm *ProviderManager) StartProvider(p ProviderType) error {
	var cmd *exec.Cmd
	switch p {
	case ProviderOllama:
		if isPortOpen("http://localhost:11434/api/tags") {
			return nil
		}
		cmd = exec.Command("ollama", "serve")
	case ProviderLMStudio:
		if isPortOpen("http://localhost:1234/v1/models") {
			return nil
		}
		// LM Studio usually needs to be started manually or via 'lms' CLI if installed
		cmd = exec.Command("lms", "server", "start")
	default:
		return fmt.Errorf("unknown provider: %s", p)
	}

	// Run in background
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", p, err)
	}

	// Wait for service to be ready
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		url := "http://localhost:11434/api/tags"
		if p == ProviderLMStudio {
			url = "http://localhost:1234/v1/models"
		}
		if isPortOpen(url) {
			return nil
		}
	}

	return fmt.Errorf("%s failed to start within timeout", p)
}
