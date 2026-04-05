package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Provider interface {
	Prompt(prompt string, system string, messages []map[string]string, model string) (string, error)
	Structure(prompt string, system string, messages []map[string]string, model string, schema interface{}) (interface{}, error)
}

type OllamaProvider struct {
	BaseURL string
}

func NewOllamaProvider(baseURL string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaProvider{BaseURL: baseURL}
}

func (p *OllamaProvider) buildMessages(prompt string, system string, history []map[string]string) []map[string]string {
	messages := []map[string]string{}
	if system != "" {
		messages = append(messages, map[string]string{"role": "system", "content": system})
	}
	for _, m := range history {
		messages = append(messages, m)
	}
	messages = append(messages, map[string]string{"role": "user", "content": prompt})
	return messages
}

func (p *OllamaProvider) Prompt(prompt string, system string, history []map[string]string, model string) (string, error) {
	messages := p.buildMessages(prompt, system, history)

	reqBody, err := json.Marshal(map[string]interface{}{
		"model":    model,
		"messages": messages,
		"stream":   false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(p.BaseURL+"/api/chat", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to call ollama at %s: %w", p.BaseURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read ollama response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errRes struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(body, &errRes); err == nil && errRes.Error != "" {
			return "", fmt.Errorf("ollama error: %s", errRes.Error)
		}
		if resp.StatusCode == http.StatusNotFound {
			return "", fmt.Errorf("ollama model '%s' not found (404). Check 'ollama list'", model)
		}
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Message struct {
			Content  string `json:"content"`
			Thinking string `json:"thinking"`
		} `json:"message"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to decode ollama response: %w (body: %s)", err, string(body))
	}

	res := result.Message.Content
	if res == "" && result.Message.Thinking != "" {
		res = result.Message.Thinking
	}

	return res, nil
}

func (p *OllamaProvider) Structure(prompt string, system string, history []map[string]string, model string, schema interface{}) (interface{}, error) {
	messages := p.buildMessages(prompt, system, history)

	reqBody, err := json.Marshal(map[string]interface{}{
		"model":    model,
		"messages": messages,
		"stream":   false,
		"format":   schema,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(p.BaseURL+"/api/chat", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call ollama at %s: %w", p.BaseURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read ollama response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode ollama response: %w (body: %s)", err, string(body))
	}

	var jsonRes interface{}
	if err := json.Unmarshal([]byte(result.Message.Content), &jsonRes); err != nil {
		return nil, fmt.Errorf("failed to parse structured output: %w (content: %s)", err, result.Message.Content)
	}

	return jsonRes, nil
}
