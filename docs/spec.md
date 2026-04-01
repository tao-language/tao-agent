# Tao Agent

A lightweight configurable agent.

Define workflows with a mixture of agent prompts, structured outputs, and deterministic control flow in a declarative type safe YAML file.

Built from the ground up to support local model providers like Ollama and LM Studio, as well as Cloud models.

## Core Components

All these components can be provided and configured by the user.

- **Tools**: (MCP server) Ways for an agent to interact with the environment.
- **Skills**: (Markdown) Instructions or knowledge given to the agent as part of its context.
- **Agents**: (YAML) A configured model with skills, tools, system prompt, and parameters.
- **Workflows**: (YAML) Deterministic control flow and orchestration of steps.

## User Interface

There should be two ways of interacting with the agent:

- Terminal UI (TUI): Minimal, lightweight, fast, minimal resources used, text only, provides basic functionality like Gemini CLI or Claude Code.
- Graphical UI (GUI): Browser server, provides a richer chat interface, accessible remotely, supports images for multi-modal agents and image generation.
- Integrations: Can integrate with messaging apps like Telegram, should allow for multi-modal agents and image generation.

# Use Cases

Must be flexible enough to allow for any kind of use case.

Some examples, but not limited to:
- Coding workflow: Research, plan, implement, review, with human feedback loop
- Asset creation: Generate and edit images, audio, 3D models, etc.
- Role playing: Virtual D&D with character and inventory management

# Types

- Null
- Boolean
- Number
- String
- Literal
- List
- Tuple
- Record
- Union
- Result

The type is always required, but can have optional descriptions and optional default values:

```yaml
# Required, no description.
- type: Number

# Required, with description.
# Descriptions are used as hints for the model.
- type: Number
  description: A score between 0 (worst) and 100 (best)
  
# Optional, because there's a default value.
- type: Number
  default: 0
```

Here are some examples of how types are represented:

```yaml
- type: Null

- type: Boolean

- type: Number

- type: String

- type: Literal
  value: 42 # required
  
# A bare value is interpreted as a Literal type.
# String clashes with types like "Number" must use the long {type: "Literal", value: "Number"}.
- 42 
  
- type: List
  items:
    type: String

- type: List
  items: String # syntax sugar for type, no description, no default
  
- type: List # if no items field, defaults to List(String)

- type: Tuple
  items: # required
  - type: Number
  - String # syntax sugar for type, no description, no default
  
- type: Record
  fields: # required
    x:
      type: Number
    y: Number # syntax sugar for type, no description, no default

- type: Union
  alternatives:
  - type: Number
  - String # syntax sugar for type, no description, no default

# For error handling
- type: Result
  ok: # defaults to String
    type: Number
  error: # defaults to String
    type: String

- type: Result
  ok: Number # syntax sugar for type, no description, no default
  # error: String
```

# Examples

## Tools Example

Tools are defined in MCP servers, they can be local or remote.

File: `tools/fs.yaml`

```yaml
# Local MCP server
version: v1.0
name: FileSystem
description: Access to local files and directories.
command: python
arguments: ["-m", "fs-mcp"]
tools:
  list-files:
    description: Lists the files for a given directory
    inputs:
      directory:
        type: String
        default: "."
    outputs:
      type: List
      items: String
  read-file:
    description: Reads a file contents
    inputs:
      path: String
    # outputs: String # defaults to String
```

```yaml
# Remote MCP server
version: v1.0
name: Calculator
description: Solves a math expression.
url: http://localhost:1234
tools:
  evaluate:
    inputs:
      expression: String
    outputs:
      type: Result
      ok: Number
```

## Skill Example

File: `skills/pdf-processing.md`

```md
---
name: pdf-processing
description: Extract text and tables from PDF files, fill forms, and merge documents. Use when handling PDF files.
license: MIT
---

# PDF Processing

## When to Use This Skill
Use this skill when the user explicitly needs to work with PDF files, such as extracting content, merging documents, or filling out form fields.

## How to Extract Text
1. Use the `pdfplumber` library for text extraction.
2. Iterate through each page of the document.
3. Extract text using `page.extract_text()`.

## How to Fill Forms
* Identify form fields and their names.
* Use the appropriate function from the bundled script in the `/scripts` directory to fill in form data.
* Save the modified PDF.

## Bundled Resources
This skill includes a `/scripts` directory with a Python script (`pdf_utils.py`) that contains helper functions for form filling and merging. The agent should leverage this script when performing these actions.
```

## Agent Examples

File: `agents/assistant.yaml`

```yaml
# A simple agent
model:
  provider: ollama
  name: qwen3.5:4b
tools: # mcp-server@tool-name: permissions
  - fs@list-files: allow # from tools/fs.yaml
  - fs@read-file # permission defaults to ask
system-prompt: |
  You are a helpful assistant.
  You can list and read files, but not modify any file.
```


File: `agents/security-reviewer.yaml`

```yaml
# A more customized agent
version: v1.0
name: Security Reviewer
description: Analyzes code and provides a security review.
model:
  provider: lmstudio
  name: qwen3.5:4b
parameters:
  context-length: 8k
  temperature: 0.2
tools:
  - fs@list-files: allow
  - fs@read-file:
      ask: "~/**"
      allow: "./**"
      deny: "/**"
inputs:
  name:
    type: String
    description: User's name
    default: user
system-prompt: |
  You are a code security expert, your job is to review code for security issues.
messages:
  # Used to "initialize" a conversation.
  user: Hi, my name is {{name}}
  agent: What code should I review for you?
```

For Cloud models:

```yaml
# A Cloud API key agent
model:
  provider: http://my-model-url.com/
  name: my-model-id
  api-key: my-provider-api-key
```

```yaml
# A Cloud OAuth agent
model:
  provider: http://my-model-url.com/
  name: my-model-id
  oauth: TODO
```

## Workflows Example

Used to orchestrate agents with deterministic control flow.

The first user message is stored in the automatic variable `{{input}}`.

```yaml
version: v1.0
name: My workflow
description: A simple workflow example
steps:
- name: Sentiment analysis
  agent: assistant
  prompt: |
    Classify the sentiment of the following:
    {{input}}
  outputs:
    type: Union
    alternatives:
    - Positive
    - Neutral
    - Negative
```
