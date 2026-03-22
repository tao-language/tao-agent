package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

type OllamaClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewOllamaClient(baseURL string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
	Options  interface{}   `json:"options,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Model         string      `json:"model"`
	Message       ChatMessage `json:"message"`
	Done          bool        `json:"done"`
	TotalDuration int64       `json:"total_duration"`
}

func (c *OllamaClient) Chat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	url := fmt.Sprintf("%s/api/chat", c.BaseURL)

	data, err := json.Marshal(req)
	if err != nil {
		return ChatResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return ChatResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return ChatResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return ChatResponse{}, fmt.Errorf("ollama error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return ChatResponse{}, err
	}

	return chatResp, nil
}

func (c *OllamaClient) IsRunning() bool {
	client := http.Client{Timeout: 500 * time.Millisecond}
	resp, err := client.Get(c.BaseURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (c *OllamaClient) Start() error {
	cmd := exec.Command("ollama", "serve")
	return cmd.Start()
}

func (c *OllamaClient) GetModels() ([]string, error) {
	url := fmt.Sprintf("%s/api/tags", c.BaseURL)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tagsResp struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return nil, err
	}

	var models []string
	for _, m := range tagsResp.Models {
		models = append(models, m.Name)
	}
	return models, nil
}

// Model Selection Logic
const (
	ModelTiny     = "llama3.2:3b" // 3B/4B
	ModelSmall    = "llama3.1:8b" // 8B/9B
	ModelMedium   = "phi4:14b"    // 14B
	ModelPowerful = "gemma2:27b"
)

func GetModelForComplexity(complexity string) string {
	switch complexity {
	case "trivial":
		return ModelTiny
	case "simple":
		return ModelSmall
	case "moderate":
		return ModelMedium
	case "complex":
		return ModelPowerful
	default:
		return ModelSmall
	}
}
