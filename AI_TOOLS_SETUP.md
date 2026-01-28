See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# AI Tool Setup Guide

## Overview
This project requires local AI tools to be installed on your computer. GitHub only stores the code - the actual AI work happens locally.

## Required Tools

### 1. Codex
**Purpose**: Code generation and review meta-tool  
**Installation**:
```bash
# Clone and install
git clone https://github.com/your-org/codex.git
cd codex
go install ./cmd/codex

# Or using package manager
npm install -g codex-cli
```

### 2. Claude CLI
**Purpose**: Large-scale architecture and refactoring  
**Installation**:
```bash
# Using npm
npm install -g @anthropic-ai/claude-cli

# Or download binary
curl -L https://github.com/anthropics/claude-cli/releases/latest/download/claude-linux-x64 -o claude
chmod +x claude
sudo mv claude /usr/local/bin/
```

### 3. Gemini CLI
**Purpose**: Research and large context analysis  
**Installation**:
```bash
# Using go install
go install github.com/google/gemini-cli/cmd/gemini@latest

# Or using cargo
cargo install gemini-cli
```

### 4. Fabric (Original)
**Purpose**: Pattern-based prompt execution  
**Installation**:
```bash
# Clone and install
git clone https://github.com/danielmiessler/fabric.git
cd fabric
npm install
npm run build

# Add to PATH
echo 'export PATH="$PATH:$HOME/fabric"' >> ~/.bashrc
```

## Automated Setup Script

Save this as `setup-ai-tools.sh` and run it:

```bash
#!/bin/bash

echo "🤖 Setting up AI tools for fabric-lite..."

# Check for existing tools
check_tool() {
    if command -v $1 &> /dev/null; then
        echo "✅ $1 is already installed"
        return 0
    else
        echo "❌ $1 not found - installing..."
        return 1
    fi
}

# Install Go if not present
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    curl -L https://go.dev/dl/go1.21.0.linux-amd64.tar.gz | tar xz
    sudo mv go /usr/local
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
fi

# Install tools
install_codex() {
    if ! check_tool "codex"; then
        echo "Installing Codex..."
        go install github.com/your-org/codex/cmd/codex@latest
    fi
}

install_claude() {
    if ! check_tool "claude"; then
        echo "Installing Claude CLI..."
        npm install -g @anthropic-ai/claude-cli
    fi
}

install_gemini() {
    if ! check_tool "gemini"; then
        echo "Installing Gemini CLI..."
        go install github.com/google/gemini-cli/cmd/gemini@latest
    fi
}

install_fabric() {
    if ! check_tool "fabric"; then
        echo "Installing Fabric..."
        git clone https://github.com/danielmiessler/fabric.git ~/.fabric
        cd ~/.fabric
        npm install
        npm run build
        echo 'export PATH="$PATH:$HOME/.fabric"' >> ~/.bashrc
    fi
}

# Run installations
install_codex
install_claude  
install_gemini
install_fabric

echo ""
echo "🎉 Setup complete! Please restart your terminal or run: source ~/.bashrc"
echo "Run 'fabric-lite list' to verify all tools are working."
```

## Quick Verify Installation

After setup, verify tools work:

```bash
# Check each tool
codex --version
claude --version  
gemini --version
fabric --list patterns

# Test fabric-lite
make build
./bin/fabric-lite list
```

## Troubleshooting

### Tool Not Found
```bash
# Add Go tools to PATH
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc

# Or manually
export PATH="$PATH:$HOME/go/bin"
```

### Permissions Error
```bash
# Fix npm global permissions
npm config set prefix ~/.npm-global
echo 'export PATH="$PATH:~/.npm-global/bin"' >> ~/.bashrc
```

### Claude CLI Issues
```bash
# Set API key
echo "export ANTHROPIC_API_KEY='your-key-here'" >> ~/.bashrc

# For local account-based auth
claude auth login
```

## Next Steps

1. Run setup script: `bash setup-ai-tools.sh`
2. Restart terminal
3. Build fabric-lite: `make build`  
4. Test: `./bin/fabric-lite --help`

**All work happens locally - GitHub is just for code sharing! 🚀**