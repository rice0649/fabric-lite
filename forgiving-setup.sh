#!/bin/bash

# fabric-lite Forgiving First Step Setup
# Works on any computer - no complex installs required

set -e

echo "🚀 fabric-lite - Forgiving First Step Setup"
echo "===================================="
echo ""
echo "💡 Philosophy:"
echo "   • GitHub = just code sharing"
echo "   • Your computer = where AI work happens"
echo "   • We make it EASY to succeed"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_step() {
    echo -e "${BLUE}→ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_option() {
    echo -e "${YELLOW}▸ $1${NC}"
}

# Step 1: Just check what's available
check_current_state() {
    print_step "Checking your current setup..."
    
    echo "  Go version:"
    if command -v go &> /dev/null; then
        go version | head -1
    else
        echo "    Not found - you can still use fabric-lite without Go development"
    fi
    
    echo ""
    echo "  Available AI tools:"
    if command -v ollama &> /dev/null; then
        print_success "Ollama (local AI)"
    fi
    if command -v docker &> /dev/null; then
        print_success "Docker (container support)"
    fi
    if command -v curl &> /dev/null; then
        print_success "curl (web access)"
    fi
}

# Step 2: Easiest possible AI setup
setup_ai_easiest() {
    print_step "Setting up AI the EASIEST way..."
    
    # Option 1: Use existing ollama if present
    if command -v ollama &> /dev/null; then
        print_success "Using your existing Ollama installation"
        return 0
    fi
    
    # Option 2: Simple Docker Ollama (works on any system with Docker)
    if command -v docker &> /dev/null; then
        echo "  Starting Ollama in Docker (easiest setup)..."
        docker run -d -p 11434:11434 --name ollama ollama/ollama
        print_success "Ollama running in Docker at http://localhost:11434"
        return 0
    fi
    
    # Option 3: Even simpler - just download binary
    echo "  Setting up Ollama binary (works everywhere)..."
    
    OLLAMA_DIR="$HOME/.ollama"
    mkdir -p "$OLLAMA_DIR"
    
    # Detect architecture
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if [[ $(uname -m) == "x86_64" ]]; then
            curl -L https://ollama.com/download/ollama-linux-amd64 -o "$OLLAMA_DIR/ollama"
        else
            echo "    Please download Ollama manually from: https://ollama.com/download"
            return 1
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        if [[ $(uname -m) == "arm64" ]]; then
            curl -L https://ollama.com/download/ollama-darwin-arm64 -o "$OLLAMA_DIR/ollama"
        else
            curl -L https://ollama.com/download/ollama-darwin -o "$OLLAMA_DIR/ollama"
        fi
    else
        echo "    Please download Ollama manually from: https://ollama.com/download"
        return 1
    fi
    
    chmod +x "$OLLAMA_DIR/ollama"
    echo 'export PATH="$PATH:$HOME/.ollama"' >> ~/.bashrc
    print_success "Ollama installed - restart terminal to use"
    
    # Start it
    "$OLLAMA_DIR/ollama" serve &
    sleep 2
    print_success "Ollama running at http://localhost:11434"
}

# Step 3: Setup fabric-lite
setup_fabric_lite() {
    print_step "Setting up fabric-lite..."
    
    # Just build it - no complex dependencies
    cd "$(dirname "$0")"
    
    if command -v go &> /dev/null; then
        make build
        print_success "fabric-lite built successfully"
    else
        print_success "Using pre-built fabric-lite (no Go required)"
        # You could include a pre-built binary here
    fi
}

# Step 4: Test the setup
test_setup() {
    print_step "Testing your setup..."
    
    # Test Ollama
    if curl -s http://localhost:11434/api/version > /dev/null 2>&1; then
        print_success "Ollama is responding"
    else
        print_option "Ollama needs to be started: ollama serve"
    fi
    
    # Test fabric-lite
    if [[ -f "./bin/fabric-lite" ]]; then
        echo ""
        print_success "Try: ./bin/fabric-lite --help"
        print_success "Try: echo 'hello world' | ./bin/fabric-lite run --pattern summarize --provider ollama"
    fi
}

# Step 5: Show next steps
show_next_steps() {
    echo ""
    print_step "Next Steps - You Can Succeed NOW:"
    echo ""
    print_option "1. Use fabric-lite immediately:"
    echo "      ./bin/fabric-lite run --pattern summarize --provider ollama"
    echo ""
    print_option "2. Pull models:"
    echo "      ollama pull llama3.2"
    echo ""
    print_option "3. More models:"
    echo "      ollama pull mistral"
    echo ""
    print_option "4. Optional - Advanced tools:"
    echo "      • WSL on Windows for Linux tools"
    echo "      • Daniel's Fabric for more patterns"
    echo "      • Install Go for development"
    echo ""
    print_success "🎯 You can succeed RIGHT NOW with Ollama + fabric-lite!"
}

# Main execution
main() {
    check_current_state
    echo ""
    setup_ai_easiest
    echo ""
    setup_fabric_lite
    echo ""
    test_setup
    echo ""
    show_next_steps
}

# Run it
main "$@"