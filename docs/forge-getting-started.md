# Getting Started with AI Project Forge

AI Project Forge orchestrates multiple AI coding assistants through structured development phases to help you build any software project.

## Installation

### Prerequisites

- Go 1.21 or later
- At least one AI CLI tool installed:
  - [fabric-lite](https://github.com/rice0649/fabric-lite) (NEW - unified patterns & tools)
  - [Gemini CLI](https://github.com/google/gemini-cli) (research & analysis)
  - [Codex CLI](https://github.com/openai/codex-cli) (coding & implementation)
  - [Claude CLI](https://github.com/anthropic/claude-cli) (advanced reasoning)
  - [OpenCode](https://github.com/opencode/opencode) (interactive coding)
  - Local Ollama (built-in to fabric-lite) for fast, private tasks

### Install from Source

```bash
# Clone the repository
git clone https://github.com/rice0649/fabric-lite
cd fabric-lite

# Build and install forge
make build-forge
make install-forge

# Verify installation
forge --version
```

### Quick Install

```bash
# Using the install script
curl -sSL https://raw.githubusercontent.com/rice0649/fabric-lite/main/scripts/install.sh | bash
```

## Your First Project

### 1. Initialize a Project

```bash
# Create a new directory
mkdir my-app && cd my-app

# Initialize with interactive mode
forge init --interactive

# Or specify options directly
forge init --name my-app --template webapp
```

### 2. Start the Discovery Phase

```bash
# View available phases
forge phase list

# Start discovery
forge phase start discovery
```

### 3. Run AI Tools

```bash
# Run the default tool for the current phase
forge run

# Specify any available tool
forge run --tool gemini        # Research with Gemini CLI
forge run --tool claude         # Advanced reasoning with Claude CLI
forge run --tool codex         # Code implementation with Codex CLI
forge run --tool opencode       # Interactive coding with OpenCode
forge run --tool ollama         # Fast local analysis with built-in Ollama

# Use fabric-lite patterns with forge
forge run --pattern research_topic   # Uses fabric-lite patterns
forge run --tool fabric        # Direct fabric-lite pattern execution
```

### 4. Complete the Phase

```bash
# Check status
forge status

# Complete with checkpoint validation
forge phase complete
```

### 5. Continue Through Phases

```bash
# Move to planning
forge phase start planning
forge run
forge phase complete

# Continue through design, implementation, testing, deployment...
```

## Project Structure

After initialization, your project will have:

```
my-app/
├── .forge/
│   ├── config.yaml       # Project configuration
│   ├── state.yaml        # Current phase and progress
│   ├── history/          # Phase completion logs
│   └── artifacts/        # AI-generated outputs
│       ├── discovery/
│       ├── planning/
│       └── ...
├── src/                  # Source code (if template used)
├── tests/                # Tests
└── docs/                 # Documentation
```

## Development Phases

| Phase | Primary Tool | Secondary Tools | Purpose |
|-------|--------------|--------------|---------|
| Discovery | gemini | fabric-lite (research) | Requirements gathering, market analysis |
| Planning | claude | opencode, fabric-lite | Architecture design, system planning |
| Design | claude | opencode, fabric-lite | API design, data modeling |
| Implementation | codex | ollama, fabric-lite | Code development, feature implementation |
| Testing | gemini | ollama, codex | Test creation, validation, coverage analysis |
| Deployment | fabric | all tools | Documentation, release automation |

## Next Steps

- Read the [Phases Guide](phases.md) for detailed phase information
- Check the [Tutorials](tutorials/) for end-to-end examples
- Configure tools in `.forge/config.yaml`

## Getting Help

```bash
# Show help
forge --help

# Get help for a specific command
forge phase --help
forge run --help
```
