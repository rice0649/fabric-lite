# fabric-lite

A lightweight AI augmentation framework inspired by [Fabric](https://github.com/danielmiessler/fabric).

## What is fabric-lite?

fabric-lite is a personal CLI tool that runs AI prompts (called "patterns") against text input. It's designed to be simple, extensible, and easy to customize.

## Features

- Pattern-based AI prompts
- Multiple AI provider support (OpenAI, Ollama, Anthropic)
- Simple CLI interface
- Easy to add custom patterns

## Installation

```bash
# Clone the repository
git clone https://github.com/rice0649/fabric-lite.git
cd fabric-lite

# Build
make build

# Install to your PATH
make install
```

## Quick Start

```bash
# Set your API key
export OPENAI_API_KEY="your-key-here"

# Run a pattern
echo "Your text here" | fabric-lite --pattern summarize

# List available patterns
fabric-lite --list

# Get help
fabric-lite --help
```

## Patterns

Patterns are reusable AI prompt templates stored in `~/.config/fabric-lite/patterns/` or the local `patterns/` directory.

Each pattern is a folder containing:
- `system.md` - The system prompt (AI identity and instructions)
- `user.md` - The user prompt template (optional)

### Example Pattern Structure

```
patterns/
└── summarize/
    ├── system.md    # "You are an expert summarizer..."
    └── user.md      # Optional user prompt
```

### Creating Your Own Pattern

```bash
mkdir -p ~/.config/fabric-lite/patterns/my-pattern
echo "# Your system prompt here" > ~/.config/fabric-lite/patterns/my-pattern/system.md
```

## Configuration

Create a config file at `~/.config/fabric-lite/config.yaml`:

```yaml
default_provider: openai
default_model: gpt-4o-mini

providers:
  openai:
    api_key: ${OPENAI_API_KEY}
  ollama:
    base_url: http://localhost:11434
  anthropic:
    api_key: ${ANTHROPIC_API_KEY}
```

## Project Structure

```
fabric-lite/
├── cmd/fabric-lite/     # CLI entry point
├── internal/
│   ├── core/            # Core logic (pattern loading, execution)
│   ├── cli/             # Command-line interface
│   └── providers/       # AI provider integrations
├── patterns/            # Built-in patterns
├── config/              # Configuration templates
├── scripts/             # Build and install scripts
└── docs/                # Documentation
```

## Development

```bash
# Run tests
make test

# Build for development
make build

# Run linter
make lint
```

## Contributing

Contributions are welcome! Feel free to:
- Add new patterns
- Improve existing patterns
- Add new AI providers
- Fix bugs or improve documentation

## License

MIT License - See [LICENSE](LICENSE) for details.

## Acknowledgments

Inspired by [Fabric](https://github.com/danielmiessler/fabric) by Daniel Miessler.
