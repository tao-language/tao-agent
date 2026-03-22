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

type OpenAICompatibleClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewOpenAICompatibleClient(baseURL string) *OpenAICompatibleClient {
	return &OpenAICompatibleClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (c *OpenAICompatibleClient) Chat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)

	// OpenAI format mapping
	type openAIMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	type openAIRequest struct {
		Model    string          `json:"model"`
		Messages []openAIMessage `json:"messages"`
	}

	messages := make([]openAIMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = openAIMessage{Role: m.Role, Content: m.Content}
	}

	data, err := json.Marshal(openAIRequest{Model: req.Model, Messages: messages})
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
		return ChatResponse{}, fmt.Errorf("openai error (status %d): %s", resp.StatusCode, string(body))
	}

	// Simplistic OpenAI response mapping
	var openAIResp struct {
		Choices []struct {
			Message openAIMessage `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return ChatResponse{}, err
	}

	if len(openAIResp.Choices) == 0 {
		return ChatResponse{}, fmt.Errorf("empty response from openai-compatible server")
	}

	return ChatResponse{
		Message: ChatMessage{
			Role:    openAIResp.Choices[0].Message.Role,
			Content: openAIResp.Choices[0].Message.Content,
		},
	}, nil
}

func (c *OpenAICompatibleClient) IsRunning() bool {
	client := http.Client{Timeout: 500 * time.Millisecond}
	resp, err := client.Get(c.BaseURL + "/models")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (c *OpenAICompatibleClient) Start() error {
	cmd := exec.Command("lms", "server", "start")
	return cmd.Start()
}

func (c *OpenAICompatibleClient) GetModels() ([]string, error) {
	url := fmt.Sprintf("%s/models", c.BaseURL)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var modelsResp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, err
	}

	var models []string
	for _, m := range modelsResp.Data {
		models = append(models, m.ID)
	}
	return models, nil
}
