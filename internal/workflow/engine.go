package workflow

import (
	"fmt"
	"path/filepath"
	"strings"

	"tao-agent/internal/agent"
	"tao-agent/internal/eval"
	mcpmanager "tao-agent/internal/mcp"
	"tao-agent/internal/provider"
	"tao-agent/internal/ui"
)

type Engine struct {
	Context    eval.Context
	Agents     map[string]*agent.Agent
	MCPManager *mcpmanager.Manager
	UI         ui.UI
}

func NewEngine(u ui.UI) *Engine {
	return &Engine{
		Context:    make(eval.Context),
		Agents:     make(map[string]*agent.Agent),
		MCPManager: mcpmanager.NewManager(),
		UI:         u,
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
		e.UI.Print(msg)
		result = msg

	case s.Ask != "":
		input := e.UI.Ask(s.Ask)
		result = input

	case s.Tool != "":
		server := "fs"
		if parts := strings.Split(s.Tool, ":"); len(parts) == 2 {
			server = parts[0]
			s.Tool = parts[1]
		}

		e.UI.Print(fmt.Sprintf("Calling tool %s:%s...", server, s.Tool))
		res, err := e.MCPManager.CallTool(server, s.Tool, s.Inputs)
		if err != nil {
			return err
		}
		e.UI.Print(fmt.Sprintf("Tool result: %s", res))
		result = res

	case s.Loop != nil:
		for {
			for _, step := range s.Loop {
				if err := e.ExecuteStep(step); err != nil {
					return err
				}
			}

			if s.Until != "" {
				untilVal, _ := eval.Evaluate(s.Until, e.Context)
				checkVal := result
				if s.For != "" {
					checkVal = e.Context[s.For]
				}

				strVal := fmt.Sprintf("%v", checkVal)
				if strVal == untilVal || strings.HasPrefix(strVal, untilVal+":") || strings.HasPrefix(strVal, untilVal+" ") {
					break
				}
			} else {
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
		} else {
			for caseKey, steps := range s.Cases {
				if strings.HasPrefix(val, caseKey+":") || strings.HasPrefix(val, caseKey+" ") {
					for _, step := range steps {
						if err := e.ExecuteStep(step); err != nil {
							return err
						}
					}
					break
				}
			}
		}

	case s.Use != "":
		e.UI.Print(fmt.Sprintf("Using workflow %s...", s.Use))
		path := s.Use
		if !filepath.IsAbs(path) {
			path = filepath.Join("workflows", s.Use)
		}

		subW, err := Load(path)
		if err != nil {
			return err
		}

		subEngine := NewEngine(e.UI) // Pass UI to sub-engine
		for k, v := range e.Context {
			subEngine.Context[k] = v
		}
		for k, v := range s.Inputs {
			evalV, _ := eval.Evaluate(fmt.Sprintf("%v", v), e.Context)
			subEngine.Context[k] = evalV
		}

		if err := subEngine.Execute(subW); err != nil {
			return err
		}
		result = subEngine.Context["result"]

	case s.Prompt != "":
		a, err := e.resolveAgent(s.Agent)
		if err != nil {
			return err
		}

		prompt, _ := eval.Evaluate(s.Prompt, e.Context)

		history := []map[string]string{}
		if a.Messages != nil {
			for role, content := range a.Messages {
				evalContent, _ := eval.Evaluate(content, e.Context)
				if role == "agent" {
					role = "assistant"
				}
				history = append(history, map[string]string{"role": role, "content": evalContent})
			}
		}

		var p provider.Provider
		switch a.Model.Provider {
		case "ollama":
			p = provider.NewOllamaProvider(a.Model.URL)
		default:
			return fmt.Errorf("unsupported provider: %s", a.Model.Provider)
		}

		if s.Outputs != nil {
			e.UI.Print("Agent: generating structured output...")
			res, err := p.Structure(prompt, a.SystemPrompt, history, a.Model.Name, s.Outputs.ToJSONSchema())
			if err != nil {
				return err
			}
			result = res
		} else {
			chunks, errs := p.Stream(prompt, a.SystemPrompt, history, a.Model.Name)
			res := e.UI.PromptStream(chunks, errs)
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
