package ui

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"tao-agent/pkg/agent"
	"tao-agent/pkg/model"
	"tao-agent/pkg/workflow"
)

//go:embed templates/*
var templates embed.FS

type Server struct {
	Workflow *workflow.Workflow
	Manager  *model.ProviderManager
}

func NewServer(w *workflow.Workflow, m *model.ProviderManager) *Server {
	return &Server{Workflow: w, Manager: m}
}

func (s *Server) Start(port int) error {
	tmpl := template.Must(template.ParseFS(templates, "templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})

	http.HandleFunc("/models", func(w http.ResponseWriter, r *http.Request) {
		provider := model.ProviderType(r.URL.Query().Get("provider"))

		var models []string
		var err error
		if provider == model.ProviderOllama {
			models, err = s.Manager.Ollama.GetModels()
		} else {
			models, err = s.Manager.LMStudio.GetModels()
		}

		if err != nil {
			fmt.Fprintf(w, "<option disabled>Error loading models: %v</option>", err)
			return
		}

		for _, m := range models {
			fmt.Fprintf(w, "<option value='%s'>%s</option>", m, m)
		}
	})

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		goal := r.FormValue("goal")
		complexity := r.FormValue("complexity")
		provider := model.ProviderType(r.FormValue("provider"))
		selectedModel := r.FormValue("model")

		task := agent.Task{
			Goal:       goal,
			Complexity: complexity,
			Provider:   provider,
			Model:      selectedModel,
		}

		result, err := s.Workflow.Run(r.Context(), task)
		if err != nil {
			fmt.Fprintf(w, "<div class='error'>Error: %v</div>", err)
			return
		}

		fmt.Fprintf(w, "<div class='result'><h3>Result</h3><pre>%s</pre></div>", result)
	})

	fmt.Printf("Tao UI starting at http://localhost:%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
