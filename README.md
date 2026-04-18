# Tao Agent

A lightweight configurable agent framework built in Go.

## Features

- **Declarative Workflows**: Define agent logic using simple YAML files.
- **Type-Safe**: Built-in type system for predictable outputs and inputs.
- **Local-First**: Native support for Ollama, LM Studio, and other local model providers.
- **MCP Support**: Interact with your environment using the Model Context Protocol.

## Installation

```bash
go build -o tao ./cmd/tao
```

## Usage

To run a simple hello world workflow:
```bash
./tao run workflows/hello.yaml
```

To run a workflow that uses an agent (requires Ollama):
```bash
./tao run workflows/prompt-test.yaml
```

## Configuration

### Agents
Agents are defined in `agents/<name>.yaml`. Example:

```yaml
name: assistant
model:
  provider: ollama
  name: qwen3.5:4b
system-prompt: |
  You are a helpful assistant.
```

## Testing

Run all unit tests:
```bash
go test ./...
```

## Status

- [x] Phase 1: Foundation (Types, Workflows, Basic CLI)
- [x] Phase 2: Agent & Basic Steps (Ollama Provider, Print/Ask/Prompt)
- [x] Phase 3: Tools & Structured Output (MCP Lite, Structured Prompting)
- [x] Phase 4: Control Flow (If, Match, Loop, Use)
- [ ] Phase 5: UI (TUI with BubbleTea, GUI/Browser)


## Documentation

- [Specification](docs/spec.md)
- [Architecture](docs/architecture/README.md)
- [Guides](docs/guides/README.md)
