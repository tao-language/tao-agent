package main

import (
	"log"
	"tao-agent/pkg/agent"
	"tao-agent/pkg/model"
	"tao-agent/pkg/ui"
	"tao-agent/pkg/workflow"
)

func main() {
	// Initialize provider manager (Ollama + LM Studio)
	manager := model.NewProviderManager()

	// Initialize agents with the manager
	registry := agent.NewRegistry(manager)

	// Initialize workflow
	wf := workflow.NewWorkflow(registry)

	// Start UI server with the manager for model fetching
	server := ui.NewServer(wf, manager)
	if err := server.Start(8080); err != nil {
		log.Fatal(err)
	}
}
