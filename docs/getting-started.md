# Getting Started with fabric-lite

This guide will help you set up and start using fabric-lite.

## Prerequisites

- Go 1.21 or later
- An API key for at least one AI provider (OpenAI, Anthropic, or local Ollama)

## Installation

### Option 1: Build from Source

```bash
git clone https://github.com/rice0649/fabric-lite.git
cd fabric-lite
make build
make install
```

### Option 2: Go Install

```bash
go install github.com/rice0649/fabric-lite/cmd/fabric-lite@latest
```

## Configuration

### 1. Set Up Your API Key

The easiest way is to set an environment variable:

```bash
# For OpenAI
export OPENAI_API_KEY="sk-your-key-here"

# For Anthropic
export ANTHROPIC_API_KEY="your-key-here"
```

Add this to your `~/.bashrc` or `~/.zshrc` to make it permanent.

### 2. Create a Config File (Optional)

```bash
mkdir -p ~/.config/fabric-lite
cp config/config.example.yaml ~/.config/fabric-lite/config.yaml
```

Edit the config file to customize your settings.

### 3. Install Patterns

```bash
make install-patterns
```

Or manually copy patterns:

```bash
cp -r patterns/* ~/.config/fabric-lite/patterns/
```

## Basic Usage

### Summarize Content

```bash
# From a file
fabric-lite --pattern summarize < article.txt

# From clipboard (Linux with xclip)
xclip -selection clipboard -o | fabric-lite --pattern summarize

# From a URL (with curl)
curl -s https://example.com/article | fabric-lite --pattern summarize
```

### Explain Code

```bash
cat mycode.py | fabric-lite --pattern explain_code
```

### Extract Ideas

```bash
fabric-lite --pattern extract_ideas < book_chapter.txt
```

### List Available Patterns

```bash
fabric-lite --list
```

## Creating Custom Patterns

1. Create a new pattern directory:

```bash
mkdir -p ~/.config/fabric-lite/patterns/my-pattern
```

2. Create the system prompt:

```bash
cat > ~/.config/fabric-lite/patterns/my-pattern/system.md << 'EOF'
# IDENTITY and PURPOSE

You are an expert at [your task here].

# OUTPUT SECTIONS

- Section 1: [description]
- Section 2: [description]

# OUTPUT INSTRUCTIONS

- Instruction 1
- Instruction 2

# INPUT:

INPUT:
EOF
```

3. Use your pattern:

```bash
echo "test input" | fabric-lite --pattern my-pattern
```

## Using with Ollama (Local Models)

1. Install and run Ollama:

```bash
curl -fsSL https://ollama.com/install.sh | sh
ollama pull llama3.2
```

2. Use with fabric-lite:

```bash
fabric-lite --pattern summarize --provider ollama --model llama3.2 < input.txt
```

## Tips

- **Pipe everything**: fabric-lite works great with Unix pipes
- **Combine with other tools**: Use with `curl`, `jq`, `xclip`, etc.
- **Version your patterns**: Keep your custom patterns in a git repo

## Troubleshooting

### "API key not found"

Make sure your API key is set:

```bash
echo $OPENAI_API_KEY
```

### "Pattern not found"

Check available patterns:

```bash
fabric-lite --list
ls ~/.config/fabric-lite/patterns/
```

## Next Steps

- Explore the [patterns](../patterns/) directory for more examples
- Read about [creating advanced patterns](./advanced-patterns.md)
- Join the community and share your patterns!
