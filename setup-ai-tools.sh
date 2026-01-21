#!/bin/bash

# AI Tools Setup Script for fabric-lite
# This script installs all required AI tools locally

set -e

echo "🤖 Setting up AI tools for fabric-lite..."
echo "    GitHub = code sharing only"
echo "    Local computer = where AI work happens"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

check_tool() {
    if command -v $1 &> /dev/null; then
        print_success "$1 is already installed"
        return 0
    else
        print_warning "$1 not found - will install"
        return 1
    fi
}

# Check system
check_system() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "Detected: Linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "Detected: macOS"
    elif [[ "$OSTYPE" == "msys" ]]; then
        echo "Detected: Windows (Git Bash)"
    else
        print_error "Unsupported system: $OSTYPE"
        exit 1
    fi
}

# Install Go if needed
install_go() {
    if command -v go &> /dev/null; then
        print_success "Go is already installed"
        return 0
    fi

    echo "Installing Go..."
    
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
        sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
        rm go1.21.0.linux-amd64.tar.gz
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        wget https://go.dev/dl/go1.21.0.darwin-amd64.tar.gz
        sudo tar -C /usr/local -xzf go1.21.0.darwin-amd64.tar.gz
        rm go1.21.0.darwin-amd64.tar.gz
    fi
    
    # Add to PATH
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
    print_success "Go installed - restart terminal to use"
}

# Install Node.js if needed
install_node() {
    if command -v npm &> /dev/null; then
        print_success "Node.js/npm is already installed"
        return 0
    fi

    echo "Installing Node.js..."
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
        sudo apt-get install -y nodejs
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        brew install node
    fi
    
    print_success "Node.js installed"
}

# Install tools
install_codex() {
    if check_tool "codex"; then
        return 0
    fi
    
    echo "Installing Codex..."
    go install github.com/rice0649/codex/cmd/codex@latest
    print_success "Codex installed"
}

install_claude() {
    if check_tool "claude"; then
        return 0
    fi
    
    echo "Installing Claude CLI..."
    npm install -g @anthropic-ai/claude-cli
    print_success "Claude CLI installed"
}

install_gemini() {
    if check_tool "gemini"; then
        return 0
    fi
    
    echo "Installing Gemini CLI..."
    go install github.com/google/gemini-cli/cmd/gemini@latest
    print_success "Gemini CLI installed"
}

install_fabric() {
    if check_tool "fabric"; then
        return 0
    fi
    
    echo "Installing Fabric (original)..."
    
    FABRIC_DIR="$HOME/.fabric"
    git clone https://github.com/danielmiessler/fabric.git "$FABRIC_DIR"
    cd "$FABRIC_DIR"
    npm install
    npm run build
    
    # Add to PATH
    echo "export PATH=\$PATH:$HOME/.fabric" >> ~/.bashrc
    print_success "Fabric installed"
}

# Main installation flow
main() {
    echo "🚀 fabric-lite AI Tools Setup"
    echo "=================================="
    check_system
    echo ""
    
    # Install prerequisites
    install_go
    install_node
    echo ""
    
    # Install AI tools
    install_codex
    install_claude
    install_gemini
    install_fabric
    echo ""
    
    # Build fabric-lite
    echo "Building fabric-lite..."
    cd "$(dirname "$0")"
    make build
    print_success "fabric-lite built"
    echo ""
    
    # Instructions
    echo "🎉 Setup Complete!"
    echo "=================="
    echo ""
    echo "IMPORTANT: Restart your terminal or run:"
    echo "source ~/.bashrc"
    echo ""
    echo "Then verify installation:"
    echo "./bin/fabric-lite list"
    echo ""
    echo "All AI work happens locally on your computer!"
    echo "GitHub is just for sharing code."
}

# Run main function
main "$@"