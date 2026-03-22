# Tao Coding Agents

Tao is a simple, performant, local-first AI orchestration system for coding tasks. 
It leverages local LLMs (via Ollama) and a Go-based backend with a minimalist HTMX UI.

## Features
- **Local-First:** Designed for 24GB+ VRAM machines.
- **Agent Registry:** Researcher, Planner, Coder, Reviewer, Quality, and Docs.
- **Workflow Orchestration:** Automatically chooses the path based on task complexity.
- **Learning:** Agents log lessons to directory-specific `README.md` files.
- **Low Overhead:** Single binary, no heavy dependencies, HTMX for interactivity.

## Prerequisites
- [Ollama](https://ollama.ai/) running locally.
- Recommended models:
  - `llama3.2:3b` (Trivial)
  - `llama3.1:8b` (Simple)
  - `phi4:14b` (Moderate)
  - `gemma2:27b` (Complex)

## Getting Started
1. Start Ollama.
2. Run the Tao server:
   ```bash
   go run cmd/tao/main.go
   ```
3. Open http://localhost:8080.
4. Enter a coding goal and select a complexity level.

## Architecture
- `pkg/agent`: Defines agent behaviors and registry.
- `pkg/model`: Local LLM integration.
- `pkg/workflow`: Orchestration logic.
- `pkg/ui`: Web-based management interface.
- `pkg/storage`: Persistent state and lesson logging.
