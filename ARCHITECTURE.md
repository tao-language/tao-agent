# Tao Coding Agents Architecture

## Principles
- **Simplicity:** Minimize complexity, favor Go's idiomatic patterns.
- **Low Overhead:** Use Go's standard library for the UI (web server + templates/HTMX).
- **Performance:** Asynchronous workflows using Go routines.
- **Local-First:** Optimized for 24GB VRAM machines (Ollama integration).
- **Learning:** Agents document lessons in `README.md` files within directories.

## Core Components

### 1. `pkg/agent`
- **Agent Interface:**
    ```go
    type Agent interface {
        Name() string
        Role() string
        Execute(ctx context.Context, task Task) (Result, error)
    }
    ```
- **Implementations:** Researcher, Planner, Coder, Reviewer, etc.

### 2. `pkg/model`
- **Ollama Client:** Wrapper around Ollama API for local model execution.
- **Model Selection:** Mapping task complexity to model sizes (4B, 9B, 27B).

### 3. `pkg/workflow`
- **Trivial:** Direct Coder execution.
- **Simple:** Coder -> Reviewer.
- **Moderate:** Researcher -> Planner -> Coder -> Reviewer.
- **Complex:** Full chain (Researcher -> Planner -> Plan Reviewer -> Coder -> Code Quality -> Code Reviewer -> Docs Writer).

### 4. `pkg/ui`
- **Web UI:** Go `net/http` server + `html/template` + HTMX for interactivity.
- **Terminal UI:** Basic CLI output for progress tracking.

### 5. `pkg/storage`
- **Tao Home:** `~/.tao/` for global state.
- **Project State:** `.tao/` within the project directory for logs and lessons.

## Directory Structure
- `cmd/tao/`: Entry point.
- `pkg/agent/`: Agent implementations.
- `pkg/model/`: LLM integration.
- `pkg/workflow/`: Orchestration logic.
- `pkg/ui/`: UI components.
- `pkg/storage/`: State and memory management.
