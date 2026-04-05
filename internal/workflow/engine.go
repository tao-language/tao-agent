package workflow

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"tao-agent/internal/agent"
	"tao-agent/internal/eval"
	"tao-agent/internal/provider"
)

type Engine struct {
	Context eval.Context
	Agents  map[string]*agent.Agent
}

func NewEngine() *Engine {
	return &Engine{
		Context: make(eval.Context),
		Agents:  make(map[string]*agent.Agent),
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

	case s.Prompt != "":
		a, err := e.resolveAgent(s.Agent)
		if err != nil {
			return err
		}

		prompt, _ := eval.Evaluate(s.Prompt, e.Context)

		// Map provider
		var p provider.Provider
		switch a.Model.Provider {
		case "ollama":
			p = provider.NewOllamaProvider(a.Model.URL)
		default:
			return fmt.Errorf("unsupported provider: %s", a.Model.Provider)
		}

		fmt.Printf("Agent (%s): thinking...\n", a.Name)
		res, err := p.Prompt(prompt, a.Model.Name)
		if err != nil {
			return err
		}
		fmt.Printf("Agent (%s): %s\n", a.Name, res)
		result = res

	default:
		return fmt.Errorf("unknown step type in step: %s", s.Name)
	}

	if s.ID != "" {
		e.Context[s.ID] = result
	}

	return err
}
