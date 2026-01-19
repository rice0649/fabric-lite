#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Installing AI Project Forge...${NC}"

# Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed.${NC}"
    echo "Please install Go 1.21 or later from https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
GO_MAJOR=$(echo "$GO_VERSION" | cut -d. -f1)
GO_MINOR=$(echo "$GO_VERSION" | cut -d. -f2)

if [[ "$GO_MAJOR" -lt 1 ]] || [[ "$GO_MAJOR" -eq 1 && "$GO_MINOR" -lt 21 ]]; then
    echo -e "${RED}Error: Go 1.21 or later is required (found $GO_VERSION)${NC}"
    exit 1
fi

echo -e "${GREEN}Found Go $GO_VERSION${NC}"

# Determine install directory
INSTALL_DIR="${GOPATH:-$HOME/go}/bin"
mkdir -p "$INSTALL_DIR"

# Check if we're in the repo or need to clone
if [[ -f "go.mod" ]] && grep -q "github.com/rice0649/fabric-lite" go.mod; then
    echo "Building from local source..."
    REPO_DIR="."
else
    echo "Cloning repository..."
    TEMP_DIR=$(mktemp -d)
    git clone https://github.com/rice0649/fabric-lite.git "$TEMP_DIR"
    REPO_DIR="$TEMP_DIR"
fi

cd "$REPO_DIR"

# Build fabric-lite
echo -e "${YELLOW}Building fabric-lite...${NC}"
go build -o "$INSTALL_DIR/fabric-lite" ./cmd/fabric-lite

# Build forge
echo -e "${YELLOW}Building forge...${NC}"
go build -o "$INSTALL_DIR/forge" ./cmd/forge

# Install patterns
PATTERNS_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/fabric-lite/patterns"
echo -e "${YELLOW}Installing patterns to $PATTERNS_DIR...${NC}"
mkdir -p "$PATTERNS_DIR"
cp -r patterns/* "$PATTERNS_DIR/"

# Cleanup temp directory if we cloned
if [[ "$REPO_DIR" != "." ]]; then
    rm -rf "$TEMP_DIR"
fi

# Verify installation
echo ""
echo -e "${GREEN}Installation complete!${NC}"
echo ""

if command -v fabric-lite &> /dev/null; then
    echo -e "  ${GREEN}✓${NC} fabric-lite installed"
else
    echo -e "  ${YELLOW}!${NC} fabric-lite installed but not in PATH"
    echo "    Add $INSTALL_DIR to your PATH"
fi

if command -v forge &> /dev/null; then
    echo -e "  ${GREEN}✓${NC} forge installed"
else
    echo -e "  ${YELLOW}!${NC} forge installed but not in PATH"
    echo "    Add $INSTALL_DIR to your PATH"
fi

echo ""
echo "Next steps:"
echo "  1. Set up your API keys:"
echo "     export OPENAI_API_KEY='your-key'"
echo "     export ANTHROPIC_API_KEY='your-key'"
echo ""
echo "  2. Create a new project:"
echo "     mkdir my-project && cd my-project"
echo "     forge init --interactive"
echo ""
echo "  3. Start building:"
echo "     forge phase start discovery"
echo "     forge run"
echo ""
echo "For more information, run: forge --help"
