package agent

import (
	"context"
	"fmt"
	"tao-agent/pkg/model"
)

type Task struct {
	Goal       string
	Context    string
	Complexity string
	Provider   model.ProviderType
	Model      string // User selected model
}

type Result struct {
	Content string
	Success bool
	Lessons string
}

type Agent interface {
	Name() string
	Role() string
	Execute(ctx context.Context, task Task) (Result, error)
}

type BaseAgent struct {
	name    string
	role    string
	manager *model.ProviderManager
}

func (b *BaseAgent) Name() string { return b.name }
func (b *BaseAgent) Role() string { return b.role }

func (b *BaseAgent) Execute(ctx context.Context, task Task) (Result, error) {
	var p model.ModelProvider
	if task.Provider == model.ProviderOllama {
		p = b.manager.Ollama
	} else {
		p = b.manager.LMStudio
	}

	// Auto-start if needed
	if err := b.manager.StartProvider(task.Provider); err != nil {
		return Result{}, fmt.Errorf("failed to ensure provider is running: %w", err)
	}

	selectedModel := task.Model
	if selectedModel == "" {
		// Fallback: try to find any installed model
		models, err := p.GetModels()
		if err == nil && len(models) > 0 {
			selectedModel = models[0]
		} else {
			// Absolute fallback to hardcoded if everything fails,
			// but we should have caught it by now.
			selectedModel = model.GetModelForComplexity(task.Complexity)
		}
	}

	systemPrompt := fmt.Sprintf("You are the %s agent. Your role is: %s.", b.name, b.role)
	userPrompt := fmt.Sprintf("Goal: %s\nContext: %s\nComplexity: %s", task.Goal, task.Context, task.Complexity)

	req := model.ChatRequest{
		Model: selectedModel,
		Messages: []model.ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: false,
	}

	resp, err := p.Chat(ctx, req)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Content: resp.Message.Content,
		Success: true,
	}, nil
}
