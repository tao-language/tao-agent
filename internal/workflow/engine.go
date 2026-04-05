package workflow

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"tao-agent/internal/agent"
	"tao-agent/internal/eval"
	mcpmanager "tao-agent/internal/mcp"
	"tao-agent/internal/provider"
)

type Engine struct {
	Context    eval.Context
	Agents     map[string]*agent.Agent
	MCPManager *mcpmanager.Manager
}

func NewEngine() *Engine {
	return &Engine{
		Context:    make(eval.Context),
		Agents:     make(map[string]*agent.Agent),
		MCPManager: mcpmanager.NewManager(),
	}
}

func (e *Engine) Execute(w *Workflow) error {
	for _, step := range w.Steps {
		if err := e.ExecuteStep(step); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) resolveAgent(name string) (*agent.Agent, error) {
	if name == "" {
		name = "default"
	}

	if a, ok := e.Agents[name]; ok {
		return a, nil
	}

	// Try to load from agents/name.yaml
	path := filepath.Join("agents", name+".yaml")
	a, err := agent.Load(path)
	if err != nil {
		if name == "default" {
			return agent.GetDefault(), nil
		}
		return nil, fmt.Errorf("failed to load agent %s: %w", name, err)
	}

	e.Agents[name] = a
	return a, nil
}

func (e *Engine) ExecuteStep(s *Step) error {
	if s.If != "" {
		val, ok := e.Context[s.If].(bool)
		if ok && !val {
			return nil
		}
	}

	var result interface{}
	var err error

	switch {
	case s.Print != "":
		msg, _ := eval.Evaluate(s.Print, e.Context)
		fmt.Println(msg)
		result = msg

	case s.Ask != "":
		fmt.Printf("%s ", s.Ask)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		result = input

	case s.Tool != "":
		// For simplicity, we assume server:tool format or just tool
		// If just tool, we might need a default server
		server := "fs" // Default for now
		if parts := strings.Split(s.Tool, ":"); len(parts) == 2 {
			server = parts[0]
			s.Tool = parts[1]
		}

		fmt.Printf("Tool: calling %s:%s...\n", server, s.Tool)
		res, err := e.MCPManager.CallTool(server, s.Tool, s.Inputs)
		if err != nil {
			return err
		}
		fmt.Printf("Tool result: %s\n", res)
		result = res

	case s.Loop != nil:
		for {
			for _, step := range s.Loop {
				if err := e.ExecuteStep(step); err != nil {
					return err
				}
			}

			// Evaluate until condition
			if s.Until != "" {
				untilVal, _ := eval.Evaluate(s.Until, e.Context)

				// Get the value to check from 'for' field or last result
				checkVal := result
				if s.For != "" {
					checkVal = e.Context[s.For]
				}

				if fmt.Sprintf("%v", checkVal) == untilVal {
					break
				}
			} else {
				// If no until, it's an infinite loop or handled elsewhere
				// For safety in MVP, let's break if no until is provided for now
				// unless it's specifically meant to be handled by an internal step
				break
			}
		}

	case s.Match != "":
		val, _ := eval.Evaluate(s.Match, e.Context)
		if steps, ok := s.Cases[val]; ok {
			for _, step := range steps {
				if err := e.ExecuteStep(step); err != nil {
					return err
				}
			}
		}

	case s.Use != "":
		fmt.Printf("Use: calling workflow %s...\n", s.Use)
		path := s.Use
		if !filepath.IsAbs(path) {
			// This might need better path resolution relative to current workflow
			path = filepath.Join("workflows", s.Use)
		}

		subW, err := Load(path)
		if err != nil {
			return err
		}

		// Create a sub-engine with inherited context for inputs
		subEngine := NewEngine()
		for k, v := range e.Context {
			subEngine.Context[k] = v
		}
		// Also add explicit inputs from the 'use' step
		for k, v := range s.Inputs {
			evalV, _ := eval.Evaluate(fmt.Sprintf("%v", v), e.Context)
			subEngine.Context[k] = evalV
		}

		if err := subEngine.Execute(subW); err != nil {
			return err
		}

		// Map back context if needed or store sub-workflow final result
		// Spec doesn't specify how 'use' returns values yet,
		// but let's assume it stores the sub-context for now or a specific 'result'
		result = subEngine.Context["result"]

	case s.Prompt != "":
		a, err := e.resolveAgent(s.Agent)
		if err != nil {
			return err
		}

		prompt, _ := eval.Evaluate(s.Prompt, e.Context)

		// Build history from agent initial messages
		history := []map[string]string{}
		if a.Messages != nil {
			// Note: The spec says user/agent keys, but Ollama uses user/assistant roles
			for role, content := range a.Messages {
				evalContent, _ := eval.Evaluate(content, e.Context)
				if role == "agent" {
					role = "assistant"
				}
				history = append(history, map[string]string{"role": role, "content": evalContent})
			}
		}

		// Map provider
		var p provider.Provider
		switch a.Model.Provider {
		case "ollama":
			p = provider.NewOllamaProvider(a.Model.URL)
		default:
			return fmt.Errorf("unsupported provider: %s", a.Model.Provider)
		}

		fmt.Printf("Agent (%s): thinking...\n", a.Name)
		if s.Outputs != nil {
			res, err := p.Structure(prompt, a.SystemPrompt, history, a.Model.Name, s.Outputs.ToJSONSchema())
			if err != nil {
				return err
			}
			fmt.Printf("Agent (%s): [Structured Output]\n", a.Name)
			result = res
		} else {
			res, err := p.Prompt(prompt, a.SystemPrompt, history, a.Model.Name)
			if err != nil {
				return err
			}
			fmt.Printf("Agent (%s): %s\n", a.Name, res)
			result = res
		}
	default:
		return fmt.Errorf("unknown step type in step: %s", s.Name)
	}

	if s.ID != "" {
		e.Context[s.ID] = result
	}

	return err
}
