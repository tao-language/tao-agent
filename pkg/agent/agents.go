package agent

import "tao-agent/pkg/model"

type Registry struct {
	Planner    *Planner
	Researcher *Researcher
	Coder      *Coder
	Reviewer   *Reviewer
	Quality    *Quality
	Docs       *Docs
}

func NewRegistry(m *model.ProviderManager) *Registry {
	return &Registry{
		Planner:    NewPlanner(m),
		Researcher: NewResearcher(m),
		Coder:      NewCoder(m),
		Reviewer:   NewReviewer(m),
		Quality:    NewQuality(m),
		Docs:       NewDocs(m),
	}
}

// Researcher Agent
type Researcher struct{ BaseAgent }

func NewResearcher(m *model.ProviderManager) *Researcher {
	return &Researcher{BaseAgent{name: "Researcher", role: "Context Gatherer", manager: m}}
}

// Coder Agent
type Coder struct{ BaseAgent }

func NewCoder(m *model.ProviderManager) *Coder {
	return &Coder{BaseAgent{name: "Coder", role: "Software Developer", manager: m}}
}

// Reviewer Agent
type Reviewer struct{ BaseAgent }

func NewReviewer(m *model.ProviderManager) *Reviewer {
	return &Reviewer{BaseAgent{name: "Reviewer", role: "Code Reviewer", manager: m}}
}

// Quality Agent
type Quality struct{ BaseAgent }

func NewQuality(m *model.ProviderManager) *Quality {
	return &Quality{BaseAgent{name: "Quality", role: "Quality Assurance", manager: m}}
}

// Docs Agent
type Docs struct{ BaseAgent }

func NewDocs(m *model.ProviderManager) *Docs {
	return &Docs{BaseAgent{name: "Docs", role: "Technical Writer", manager: m}}
}
