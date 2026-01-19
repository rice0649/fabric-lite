# Getting Started with AI Project Forge

AI Project Forge orchestrates multiple AI coding assistants through structured development phases to help you build any software project.

## Installation

### Prerequisites

- Go 1.21 or later
- At least one AI CLI tool installed:
  - [Gemini CLI](https://github.com/google/gemini-cli) (recommended for discovery/testing)
  - [Codex CLI](https://github.com/openai/codex-cli) (recommended for implementation)
  - [OpenCode](https://github.com/opencode/opencode) (recommended for planning/design)
  - [fabric-lite](https://github.com/rice0649/fabric-lite) (recommended for documentation)

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

# Or specify a tool
forge run --tool gemini

# Use a fabric-lite pattern
forge run --pattern research_topic
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

| Phase | Tool | Purpose |
|-------|------|---------|
| Discovery | Gemini CLI | Research, requirements gathering |
| Planning | OpenCode | Architecture design |
| Design | OpenCode | API and data modeling |
| Implementation | Codex CLI | Code development |
| Testing | Gemini CLI | Test creation, coverage analysis |
| Deployment | fabric-lite | Documentation, release notes |

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
