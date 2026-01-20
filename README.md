# fabric-lite

A lightweight AI augmentation framework inspired by [Fabric](https://github.com/danielmiessler/fabric).

## What is fabric-lite?

fabric-lite is a personal CLI tool that runs AI prompts (called "patterns") against text input. It's designed to be simple, extensible, and easy to customize.

## Features

- **Pattern-based AI prompts** - Reusable AI prompt templates for common tasks
- **Multi-provider AI support** - OpenAI, Anthropic, Ollama, and executable providers
- **Direct tool invocation** - Execute specialized AI tools via unified CLI interface
- **Meta-tool delegation** - Codex tool delegates to other providers for coding tasks
- **Simple CLI interface** - Single `run` command for both patterns and tools
- **Easy extensibility** - Add custom patterns and tools via simple interfaces

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
# Clone the repository
git clone https://github.com/rice0649/fabric-lite.git
cd fabric-lite

# Build
make build

# Install to your PATH
make install

# Set your API key
export OPENAI_API_KEY="your-key-here"

# Run a pattern (original capability)
echo "Your text here" | fabric-lite --pattern summarize

# Execute a tool directly (NEW capability)
fabric-lite run codex -P "implement REST API in Go"
fabric-lite run ollama -P "analyze this data"
fabric-lite run gemini -P "research best practices"
fabric-lite run claude -P "design system architecture"

# List available patterns
fabric-lite list

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

## Available Tools

fabric-lite now includes a comprehensive tool ecosystem for specialized AI tasks:

| Tool | Description | Best Use Case | Provider |
|------|-------------|--------------|----------|
| **codex** | Meta-coding assistant | Delegates to any configured provider |
| **gemini** | Research & analysis | Research tasks, documentation |
| **claude** | Advanced reasoning | Architecture design, complex analysis |
| **opencode** | Interactive coding | Real-time code assistance |
| **ollama** | Local processing | Fast, private tasks |
| **fabric** | Pattern execution | Documentation, structured outputs |

### Media & Content Processing

fabric-lite supports integration with external transcription tools for video and audio analysis:

```bash
# YouTube transcription using yt-dlp (recommended)
yt-dlp --write-auto-sub --write-sub "https://www.youtube.com/watch?v=VIDEO_ID" --sub-format "srt" -o "transcript.srt"

# Analyze transcript with fabric-lite tools
./bin/fabric-lite run gemini -P "analyze this transcript for insights: $(cat transcript.srt)"
./bin/fabric-lite run codex -P "create analysis script based on: $(cat transcript.srt)"

# Privacy-focused alternative (invidious instance)
curl "http://localhost:8080/api/v1/videos/VIDEO_ID?format=json"
```

### Tool Execution Examples

```bash
# Direct tool invocation
fabric-lite run codex -P "implement REST API with validation"
fabric-lite run gemini -P "research microservices architecture patterns"
fabric-lite run claude -P "design scalable system architecture"
fabric-lite run ollama -P "quickly analyze this dataset"

# Pattern execution with tool override
fabric-lite run --pattern summarize --provider claude input.txt
fabric-lite run --pattern code_review --provider codex --model gpt-4
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
