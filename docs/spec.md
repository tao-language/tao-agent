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

Even though agents are user-defined, there should be a default (fallback) agent.

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
- Role playing: Virtual D&D with character and inventory management.

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

# Workflows

They are defined by a sequence of steps.

These are the supported steps:

- **prompt**: prompts an agent, can return structured outputs
- **ask**: asks the user, useful for feedback loops
- **tool**: directly calls a tool
- **print**: prints a message directly
- **loop**: loops until a condition is met
- **match**: branches to different cases depending on the input
- **use**: calls another workflow file, like a function 

You can interpolate variables like `{{var_name}}`.

The first user message is stored in the automatic variable `{{input}}`.

Every step can have an optional `id` field to reference.

Every step can have an optional `if` field to only run if a `Bool` condition is met.

Runtime session state should be saved to a JSON file every time a step finishes to allow resuming a session.

## Prompt

Triggers if a step has a `prompt` field.

```yaml
steps:
# Minimal case
- prompt: Prompt to send to the agent

# Custom case
- name: Display name # defaults to "<agent>: <prompt-first-line>"
  id: my-prompt-id # defaults to None
  if: {{condition}} # defaults to True
  prompt: Prompt to send to the agent
  agent: agent-name # defaults to the default/fallback agent
  parameters: # defaults to what is defined in the agent yaml file
    temperature: 0.2
  tools: # defaults to what is define din the agent yaml file
  - write-file # this is added on top of what is defined by the agent
  outputs: # defaults to String
    type: Bool
```

## Ask

Triggers if a step has a `ask` field.

```yaml
steps:
# Minimal case
- ask: What's your name?

# Custom case
- name: Human approval
  id: my-ask-id # defaults to None
  if: {{condition}} # defaults to True
  ask: Do you approve?
  outputs: Bool # defaults to String
  outputs-agent: default # model to interpret the output type if not String, defaults to default agent
  choices: # message-display: output-value
  - Yes: True
  - No: False
  allow-open-answer: False # defaults to True, if True it's interpreted by the agent on the expected type
  default: Yes # defaults to None, default selected value to "ghost" input in TUI
```

## Tool

Triggers if a step has a `tool` field.

```yaml
steps:
# Minimal case
- tool: read-file
  inputs:
    path: my-file.txt

# Custom case
- name: Read file
  id: my-tool-id # defaults to None
  if: {{condition}} # defaults to True
  tool: read-file
  inputs:
    path: my-file.txt
```

## Print

Triggers if a step has a `print` field.

```yaml
steps:
# Minimal case
- print: Some message

# Custom case
- name: Print status
  id: my-print-id # defaults to None
  if: {{condition}} # defaults to True
  print: Approved
```

## Loop

Triggers if a step has a `loop` field.

The `until` checks the last output of the `loop` steps by default, unless a `for` is provided.

```yaml
steps:
# Minimal case
- until: Exit
  loop:
  - ask: How to proceed?
    outputs:
      Continue:
      Exit:

# Custom case
- name: Loop until user wants to exit
  id: my-loop-id # defaults to None
  if: {{condition}} # defaults to True
  for: status
  until: Exit
  loop:
  - id: status # used for loop.until
    ask: How to proceed?
    outputs:
      Continue:
      Exit:
```

## Match

Triggers if a step has a `match` field.

```yaml
steps:
# Minimal case

# Custom case

```

## Use

Triggers if a step has a `use` field.

```yaml
steps:
# Minimal case
- use: path/to/my-workflow.yaml

# Custom case
- use: path/to/my-custom-workflow.yaml
  inputs:
    filename: {{file}}
    retries: 3
```

## Conditional steps

All steps support an `if` field to run _only_ if the condition is met.

The condition **must** be a `Bool`.

```yaml
steps:
- id: approved
  prompt: |
    Do you approve of this?
    {{input}}
  outputs: Bool
- name: Conditional step
  if: {{approved}} # this must be a Bool
  print: ✅ Approved
```

## Invalid steps

```yaml
steps:
- name: Step cannot be identified
- name: Ambiguous step kind, 'use' or 'prompt'?
  use: my-workflow.yaml
  prompt: Hello
```

# Examples

## Tools Example

Tools are defined in MCP servers, they can be local, remote, or workflow YAML files.

MCP servers expose `tools/list` to list all the available tools in the server.

```yaml
# Local MCP server
version: v1.0
name: FileSystem
description: Access to local files and directories.
command: python
arguments: ["-m", "fs-mcp"]
```

```yaml
# Remote MCP server
version: v1.0
name: Calculator
description: Solves a math expression.
url: http://localhost:1234
```

For local workflow files, each workflow function must be defined explicitly.
Each workflow file includes typed inputs and outputs, as well as descriptions.

```yaml
# Workflow tool
version: v1.0
name: Workflow tools
description: Tools defined in workflow files.
tools:
  feedback-loop:
    use: workflows/tools/feedback-loop.yaml
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
tools: # tool-name: permissions
  - list-files: allow # from tools/fs.yaml
  - read-file # permission defaults to ask
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
  temperature: 0.2
context:
  length: 8k
  full-strategy: summarize # summarize | sliding-window | truncate-start | truncate-end | truncate-middle
tools:
  - list-files: allow
  - read-file:
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
  api-key: {{env.MY_API_KEY}}
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

```yaml
version: v1.0
name: Human feedback loop
description: Demonstrate agent reviewer with human feedback loop
steps:
- name: Draft
  id: draft
  agent: assistant
  prompt: |
    Write a draft for:
    {{input}}
- name: Feedback loop
  for: status
  until: Approved
  loop:
  - name: Agent review
    id: status
    prompt: |
      Review the following:
      {{draft}}
    outputs:
    - Approved: Reviewer agent approves
    - Feedback: Found issues that must be addressed
  - name: Get user feedback
    match: {{status}}
    cases:
      Approved:
      - id: {{status}}
        ask: Looks good to me, do you approve?
        outputs:
        - Approved: User approves
        - Feedback: User provided feedback to address
      Feedback: # Reviewer feedback, address directly without user intervention
  - name: Address feedback
    match: {{status}}
    cases:
      Approved: # Approved, do nothing, will exit loop
      Feedback:
      - name: Address feedback
        id: draft
        prompt: |
          Please address the provided feedback on the following draft:
          Feedback:
          {{status}}
          Draft:
          {{draft}}
```
