# Tao Agent Implementation Plan - Phase 1: Foundation

This plan covers the initialization of the Tao Agent project, implementing the core type system, and setting up the basic CLI structure in Go.

## 1. Project Structure

```text
tao-agent/
├── cmd/
│   └── tao/                # Main entry point
│       └── main.go
├── internal/
│   ├── types/              # Core Tao type system (Null, Bool, Number, etc.)
│   │   ├── types.go
│   │   └── types_test.go
│   ├── workflow/           # Workflow engine and YAML parsing
│   │   ├── workflow.go
│   │   └── step.go
│   ├── eval/               # Expression evaluator for {{vars}}
│   │   └── eval.go
│   └── provider/           # LLM Providers (Ollama, LM Studio, etc.)
│       └── provider.go
├── pkg/
│   └── mcp/                # MCP protocol implementation (Lite)
├── docs/                   # Organized documentation
│   ├── architecture/
│   ├── guides/
│   └── reference/
├── go.mod
└── README.md
```

## 2. Core Type System (`internal/types`)

Implement the types defined in `docs/spec.md`:
- `Null`, `Boolean`, `Number`, `String`, `Literal`
- `List`, `Tuple`, `Record`, `Union`, `Result`
- Each type will have:
    - `Kind` (enum)
    - `Description` (optional)
    - `Default` (optional)
    - `Value` (interface{})

## 3. Workflow Engine (`internal/workflow`)

- Define `Workflow` and `Step` structs that map to the YAML specification.
- Implement a `Loader` that reads YAML files and validates them against the Tao type system.

## 4. Expression Evaluator (`internal/eval`)

- Use Go's `text/template` for basic `{{var}}` interpolation in strings.
- Implement a simple boolean evaluator for `if` conditions.

## 5. Implementation Steps

1.  **Initialize Directory Structure:** Create the directories as defined in Section 1.
2.  **Define Core Types:** Create `internal/types/types.go` with the base type definitions and JSON/YAML marshaling.
3.  **Implement YAML Parsing:** Create `internal/workflow/workflow.go` to parse the `steps` structure.
4.  **Basic CLI:** Set up `cmd/tao/main.go` using `spf13/cobra` to handle commands like `run`.
5.  **Documentation:** Create the initial `README.md` and organize the `docs/` directory.
6.  **Testing:** Add unit tests for type validation and YAML loading.

## 6. Verification Strategy

- **Unit Tests:** `go test ./internal/...` to verify type logic and parsing.
- **Integration Test:** A minimal YAML workflow that prints a message and performs a basic calculation.
- **Benchmarks:** Monitor memory usage during parsing and execution to ensure "low bloat".
