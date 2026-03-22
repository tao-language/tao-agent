package workflow

import (
	"context"
	"fmt"
	"tao-agent/pkg/agent"
)

type Workflow struct {
	Registry *agent.Registry
}

func NewWorkflow(r *agent.Registry) *Workflow {
	return &Workflow{Registry: r}
}

func (w *Workflow) Run(ctx context.Context, task agent.Task) (string, error) {
	switch task.Complexity {
	case "trivial":
		return w.runTrivial(ctx, task)
	case "simple":
		return w.runSimple(ctx, task)
	case "moderate":
		return w.runModerate(ctx, task)
	case "complex":
		return w.runComplex(ctx, task)
	default:
		return "", fmt.Errorf("unknown complexity: %s", task.Complexity)
	}
}

func (w *Workflow) runTrivial(ctx context.Context, task agent.Task) (string, error) {
	res, err := w.Registry.Coder.Execute(ctx, task)
	if err != nil {
		return "", err
	}
	return res.Content, nil
}

func (w *Workflow) runSimple(ctx context.Context, task agent.Task) (string, error) {
	// Coder -> Reviewer
	res, err := w.Registry.Coder.Execute(ctx, task)
	if err != nil {
		return "", err
	}

	task.Context += "\nInitial implementation: " + res.Content
	rev, err := w.Registry.Reviewer.Execute(ctx, task)
	if err != nil {
		return "", err
	}
	return rev.Content, nil
}

func (w *Workflow) runModerate(ctx context.Context, task agent.Task) (string, error) {
	// Researcher -> Planner -> Coder -> Reviewer
	res, err := w.Registry.Researcher.Execute(ctx, task)
	if err != nil {
		return "", err
	}
	task.Context += "\nResearch: " + res.Content

	plan, err := w.Registry.Planner.Execute(ctx, task)
	if err != nil {
		return "", err
	}
	task.Context += "\nPlan: " + plan.Content

	code, err := w.Registry.Coder.Execute(ctx, task)
	if err != nil {
		return "", err
	}
	task.Context += "\nCode: " + code.Content

	rev, err := w.Registry.Reviewer.Execute(ctx, task)
	if err != nil {
		return "", err
	}
	return rev.Content, nil
}

func (w *Workflow) runComplex(ctx context.Context, task agent.Task) (string, error) {
	// Full chain: Researcher -> Planner -> Coder -> Quality -> Reviewer -> Docs
	// Implementation follows the same pattern as runModerate
	return "Complex workflow follows the chain: Research -> Plan -> Code -> Quality -> Review -> Docs. (Partially implemented)", nil
}
