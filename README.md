# fabric-lite

A lightweight AI augmentation framework inspired by [Fabric](https://github.com/danielmiessler/fabric).

## 🎯 FORGIVING FIRST STEP - Success in 60 Seconds

**No complex installs required!** Just run this:

```bash
curl -fsSL https://raw.githubusercontent.com/rice0649/fabric-lite/main/forgiving-setup.sh | bash
```

This automatically:
- ✅ Sets up Ollama (local AI) - works on ANY computer
- ✅ Builds fabric-lite (no Go required for users)  
- ✅ Tests everything works
- ✅ You're ready to use AI tools locally!

## What is fabric-lite?

fabric-lite is a personal CLI tool that runs AI prompts (called "patterns") against text input. It's designed to be simple, extensible, and easy to customize.

## The Simple Philosophy

**GitHub = Code Sharing Only**  
**Your Computer = Where AI Work Happens**

When you clone this repo:
1. You get the code from GitHub
2. You run the setup script 
3. AI tools are installed and run on **your computer**
4. All AI processing happens locally
5. No cloud services, no API keys required

## Features

- **Pattern-based AI prompts** - Reusable AI prompt templates for common tasks
- **Multi-provider AI support** - Ollama (local), OpenAI, Anthropic, and more
- **Direct tool invocation** - Execute specialized AI tools via unified CLI interface
- **Meta-tool delegation** - Codex tool delegates to other providers for coding tasks
- **Simple CLI interface** - Single `run` command for both patterns and tools
- **Easy extensibility** - Add custom patterns and tools via simple interfaces

## Quick Start After Setup

```bash
# Try fabric-lite immediately
./bin/fabric-lite run --pattern summarize --provider ollama

# Or pipe input
echo "Your text here" | ./bin/fabric-lite run --pattern extract_key_points --provider ollama

# List available patterns  
./bin/fabric-lite list

# Get help
./bin/fabric-lite --help
```

## If Setup Script Fails

**Don't worry!** Manual options:

### Option 1: Ollama (Easiest)
```bash
# Install Ollama (works on any OS)
curl -fsSL https://ollama.com/install.sh | sh

# Start it
ollama serve

# Pull a model
ollama pull llama3.2
```

### Option 2: Docker (Universal)
```bash
# Use Ollama in Docker - no install needed
docker run -d -p 11434:11434 --name ollama ollama/ollama

# Pull and use models
docker exec ollama ollama pull llama3.2
```

### Option 3: Pre-built Binary
```bash
# Download pre-built fabric-lite (no Go needed)
wget https://github.com/rice0649/fabric-lite/releases/latest/download/fabric-lite-linux

chmod +x fabric-lite-linux
./fabric-lite-linux --help
```

## 📖 Complete User Guide

**📘 New to AI tools?** Read our comprehensive user manual:
[USER_MANUAL.md](./USER_MANUAL.md)

**📕 Microsoft Word version:** [fabric-lite-user-manual.docx](./fabric-lite-user-manual.docx)

Written specifically for complete beginners - no technical knowledge required!

### What's Inside:
- ✅ 5-minute getting started guide
- ✅ Step-by-step instructions with pictures
- ✅ Real examples you can use immediately  
- ✅ Troubleshooting for common issues
- ✅ Pro tips and shortcuts
- ✅ Safety and privacy guidance

## Advanced Options (Optional)

If you want more than basic Ollama:

- **WSL on Windows**: Run Linux tools inside Windows
- **Docker Containers**: Run any AI tool in containers
- **Daniel's Original Fabric**: More patterns and features
- **Install Go**: For development and building from source

## Project Structure

- `cmd/fabric-lite/` holds the CLI entry point and wiring
- `internal/` contains core logic, CLI handling, and provider integrations
- `patterns/` stores built-in prompt patterns (e.g., `summarize/`, `extract_ideas/`)
- `config/` provides configuration templates; `docs/` is project documentation
- `agents/` contains orchestration assets; `scripts/` has build/install helpers

## Build & Test Commands

```bash
make build    # Build the CLI to bin/fabric-lite
make test      # Run tests
make install   # Install to ~/go/bin or ~/bin
```

## Testing Coverage

All packages have comprehensive test coverage:
- Providers: 42.2% coverage
- Tools: 31.1% coverage
- Executor: 90.7% coverage  
- Core: 13.7% coverage
- CLI: 7.0% coverage

## Contributing

1. Fork on GitHub (code sharing)
2. Clone locally (where AI work happens)  
3. Set up local AI tools
4. Make changes
5. Test locally with your AI tools
6. Push back to GitHub

## The Big Picture

This project follows a **local-first** philosophy:
- **GitHub** = just for sharing code with others
- **Local Computer** = where all AI processing happens
- **No Dependencies** = works with whatever tools you have locally
- **Forgiving** = multiple paths to success, no complex requirements

**You can succeed with JUST Ollama + fabric-lite!** 🚀