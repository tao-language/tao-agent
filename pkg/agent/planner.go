package agent

import (
	"context"
	"fmt"
	"tao-agent/pkg/model"
)

type Planner struct {
	BaseAgent
}

func NewPlanner(m *model.ProviderManager) *Planner {
	return &Planner{
		BaseAgent: BaseAgent{
			name:    "Planner",
			role:    "Architect and Plan Writer",
			manager: m,
		},
	}
}

func (p *Planner) Execute(ctx context.Context, task Task) (Result, error) {
	var mp model.ModelProvider
	if task.Provider == model.ProviderOllama {
		mp = p.manager.Ollama
	} else {
		mp = p.manager.LMStudio
	}

	// Auto-start if needed
	if err := p.manager.StartProvider(task.Provider); err != nil {
		return Result{}, fmt.Errorf("failed to ensure provider is running: %w", err)
	}

	selectedModel := task.Model
	if selectedModel == "" {
		models, err := mp.GetModels()
		if err == nil && len(models) > 0 {
			selectedModel = models[0]
		} else {
			selectedModel = model.GetModelForComplexity(task.Complexity)
		}
	}

	systemPrompt := "You are an expert software architect. Your goal is to write detailed plans for tasks."
	userPrompt := fmt.Sprintf("Goal: %s\nContext: %s\nComplexity: %s\nPlease write a plan using the standard template.", task.Goal, task.Context, task.Complexity)

	req := model.ChatRequest{
		Model: selectedModel,
		Messages: []model.ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: false,
	}

	resp, err := mp.Chat(ctx, req)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Content: resp.Message.Content,
		Success: true,
	}, nil
}
