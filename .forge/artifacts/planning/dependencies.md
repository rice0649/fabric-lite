See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# Dependencies Identification - Fabric-Lite Multi-Agent System

## External Service Dependencies

### 1. AI Provider Services

#### Cloud AI Providers (Required)
**OpenAI Platform**
- **API Endpoint**: https://api.openai.com/v1
- **Authentication**: Bearer token (API key)
- **Models Available**:
  - `gpt-4o`: High-quality multimodal ($0.015/1K input, $0.06/1K output)
  - `gpt-4o-mini`: Fast, cost-effective ($0.00015/1K input, $0.0006/1K output)
  - `o3-mini`: Advanced reasoning ($0.0015/1K input, $0.006/1K output)
- **Rate Limits**: 10K requests/minute, 500K tokens/minute
- **Regional Availability**: Global
- **Dependencies**: HTTP client, TLS 1.3+

**Anthropic Claude**
- **API Endpoint**: https://api.anthropic.com/v1
- **Authentication**: Bearer token (API key)
- **Models Available**:
  - `claude-3-5-sonnet-20241022`: High-performance ($0.015/1K input, $0.075/1K output)
  - `claude-3-haiku-20240307`: Fast, cost-effective ($0.00025/1K input, $0.00125/1K output)
- **Rate Limits**: 5K requests/minute, 400K tokens/minute
- **Regional Availability**: US, EU
- **Dependencies**: HTTP client, TLS 1.3+

**Google AI Platform**
- **API Endpoint**: https://generativelanguage.googleapis.com/v1
- **Authentication**: OAuth 2.0 or API key
- **Models Available**:
  - `gemini-2.0-flash-exp`: 1M context, search integration ($0.000125/1K tokens)
  - `gemini-1.5-pro`: High-quality reasoning ($0.0025/1K tokens)
- **Rate Limits**: 15K requests/minute, 2M tokens/minute
- **Regional Availability**: Global
- **Dependencies**: HTTP client, OAuth 2.0 library

#### Local AI Providers (Optional)
**Ollama**
- **Installation**: curl -fsSL https://ollama.ai/install.sh | sh
- **Models Available**:
  - `llama3.1:8b`: 8B parameters (4.7GB download)
  - `mistral:7b`: 7B parameters (4.1GB download)
  - `codellama:7b`: Code specialized (3.8GB download)
- **Resource Requirements**: 8GB+ RAM for 7B models
- **Dependencies**: Local server, HTTP client
- **Cost**: Free (hardware required)

### 2. CLI Tool Dependencies

#### Required AI CLI Tools
**Gemini CLI**
- **Installation**: npm install -g @anthropic/gemini-cli
- **Purpose**: Google Gemini integration with web search
- **Authentication**: Google account
- **Dependencies**: Node.js 18+, npm
- **Cost**: Free tier available

**Codex CLI**
- **Installation**: npm install -g @openai/codex-cli
- **Purpose**: OpenAI o3-mini code generation
- **Authentication**: OpenAI API key
- **Dependencies**: Node.js 18+, npm
- **Cost**: Based on OpenAI usage

**OpenCode CLI**
- **Installation**: go install github.com/opencode/opencode@latest
- **Purpose**: Anthropic Claude integration for architecture
- **Authentication**: Anthropic API key
- **Dependencies**: Go 1.19+
- **Cost**: Based on Anthropic usage

**Cursor CLI**
- **Installation**: npm install cursor-cli
- **Purpose**: Terminal-based AI coding
- **Authentication**: Cursor account
- **Dependencies**: Node.js 18+, npm
- **Cost**: Subscription based

#### Browser Automation
**agent-browser**
- **Installation**: npm install agent-browser
- **Purpose**: AI browser automation (Vercel Labs)
- **Features**: 93% less context than Playwright MCP
- **Dependencies**: Node.js 18+, Chromium
- **Cost**: Free

### 3. Development & Build Dependencies

#### Core Development Tools
**Go Runtime**
- **Version**: 1.21+ (required)
- **Purpose**: Core language runtime
- **Installation**: 
  - Linux: `sudo apt install golang-go` or download from golang.org
  - macOS: `brew install go`
  - Windows: Download installer from golang.org
- **Dependencies**: Git, C compiler (for some packages)

**Node.js Runtime**
- **Version**: 18+ (required for CLI tools)
- **Purpose**: Run AI CLI tools
- **Installation**: 
  - Linux: `sudo apt install nodejs npm`
  - macOS: `brew install node`
  - Windows: Download from nodejs.org
- **Dependencies**: npm

#### Development Utilities
**Go Development Tools**
```bash
# Linting and formatting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install mvdan.cc/gofumpt@latest

# Testing and coverage
go install github.com/gotestyourself/gotest/tools/gotestsum@latest
go install github.com/fzipp/gocycmd/cmd/gocycmd@latest

# Live reload
go install github.com/air-verse/air@latest

# Dependency management
go install github.com/psampaz/go-mod-outdated@latest
```

**Docker & Containerization**
```bash
# Docker (for development and deployment)
# Linux: sudo apt install docker.io docker-compose
# macOS: Download Docker Desktop
# Windows: Download Docker Desktop

# Verification
docker --version  # Should be 20.10+
docker-compose --version  # Should be 2.0+
```

## Internal Dependencies

### 1. Go Module Dependencies

#### Core Framework Dependencies
```go
// CLI Framework
github.com/spf13/cobra v1.8.0
github.com/spf13/pflag v1.0.5
github.com/spf13/viper v1.17.0

// Configuration
gopkg.in/yaml.v3 v3.0.1
github.com/kelseyhightower/envconfig v1.4.0

// Logging
go.uber.org/zap v1.26.0
go.uber.org/multierr v1.11.0

// Utilities
github.com/google/uuid v1.4.0
github.com/pkg/errors v0.9.1
github.com/imdario/mergo v0.3.16
```

#### AI Provider Clients
```go
// OpenAI
github.com/sashabaranov/go-openai v1.17.9

// Anthropic
github.com/anthropics/anthropic-go v0.4.0

// Google AI
cloud.google.com/go/ai v0.24.0
google.golang.org/api v0.154.0

// Local models
github.com/ollama/ollama-go v0.1.0
```

#### Database & Storage
```go
// ORM
gorm.io/gorm v1.25.5
gorm.io/driver/sqlite v1.5.4

// Caching
github.com/dgraph-io/badger/v3 v3.2103.5
github.com/go-redis/redis/v8 v8.11.5
github.com/patrickmn/go-cache v2.1.0+incompatible
```

#### HTTP & Networking
```go
// HTTP clients
github.com/valyala/fasthttp v1.51.0
github.com/hashicorp/go-retryablehttp v0.7.4

// WebSockets
github.com/gorilla/websocket v1.5.1
nhooyr.io/websocket v1.8.10

// HTTP utilities
github.com/stretchr/testify v1.8.4
github.com/jarcoal/httpmock v1.3.1
```

#### Testing & Quality
```go
// Testing framework
github.com/stretchr/testify v1.8.4
github.com/golang/mock v1.6.0
github.com/onsi/ginkgo/v2 v2.13.0
github.com/onsi/gomega v1.30.0

// Coverage
github.com/wadey/gocovmerge v0.1.0
github.com/matm/gocov-html v1.3.0
```

### 2. System Dependencies

#### Operating System Support
**Linux**
- **Distributions**: Ubuntu 20.04+, Debian 11+, CentOS 8+, Fedora 36+
- **Libraries**: glibc 2.17+, libssl 1.1.1+
- **Architecture**: amd64, arm64
- **Package Manager**: apt, yum, dnf, pacman

**macOS**
- **Versions**: macOS 11.0+ (Big Sur)
- **Architecture**: amd64, arm64 (Apple Silicon)
- **Package Manager**: Homebrew recommended
- **Dependencies**: Xcode Command Line Tools

**Windows**
- **Versions**: Windows 10+ (build 19042+)
- **Architecture**: amd64
- **Dependencies**: Windows 10 SDK
- **Package Manager**: Chocolatey, Scoop optional

#### Hardware Requirements
**Minimum Requirements**
- **CPU**: 2 cores, x86_64 or arm64
- **Memory**: 4GB RAM (8GB recommended for local models)
- **Storage**: 500MB free space (2GB+ for local models)
- **Network**: Internet connection for cloud providers

**Recommended Requirements**
- **CPU**: 4+ cores, x86_64 or arm64
- **Memory**: 16GB RAM (32GB for multiple local models)
- **Storage**: 10GB+ SSD storage
- **Network**: Broadband connection for optimal performance

## Security Dependencies

### 1. Encryption & Security Libraries
```go
// Cryptography
golang.org/x/crypto v0.15.0
github.com/golang-jwt/jwt/v5 v5.0.0
golang.org/x/oauth2 v0.12.0

// TLS/SSL
crypto/tls (standard library)
github.com/fullsailor/pkcs7 v0.0.1
```

### 2. Authentication & Authorization
```go
// API key management
github.com/aws/aws-sdk-go-v2 v1.21.0
github.com/hashicorp/vault/api v1.9.0

// Environment variable security
github.com/joho/godotenv v1.5.1
```

## Monitoring & Observability Dependencies

### 1. Metrics & Telemetry
```go
// Prometheus metrics
github.com/prometheus/client_golang v1.17.0
github.com/prometheus/common v0.44.0

// OpenTelemetry
go.opentelemetry.io/otel v1.21.0
go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0
go.opentelemetry.io/otel/sdk v1.21.0
```

### 2. Health Checks
```go
// Health monitoring
github.com/heptiolabs/healthcheck v4.3.0
github.com/alexlod/health-go v1.0.0
```

## Integration Dependencies

### 1. CI/CD Platforms

**GitHub Actions** (Primary)
- **YAML Configuration**: `.github/workflows/`
- **Runners**: ubuntu-latest, macos-latest, windows-latest
- **Dependencies**: Pre-built Go toolchain cache
- **Cost**: Free for public repos, $200/month private

**GitLab CI** (Alternative)
- **YAML Configuration**: `.gitlab-ci.yml`
- **Runners**: Docker-based
- **Dependencies**: GitLab Runner installation
- **Cost**: Self-hosted or GitLab.com minutes

### 2. Container Orchestration

**Kubernetes** (Optional)
- **Version**: 1.25+ recommended
- **Dependencies**: kubectl, cluster access
- **Manifests**: Helm charts or Kubernetes YAML
- **Use Cases**: Multi-agent orchestration at scale

**Docker Compose** (Development)
- **Version**: 2.0+ recommended
- **File**: `docker-compose.yml`
- **Dependencies**: Docker daemon
- **Use Cases**: Local development, testing

## Optional Dependencies

### 1. Enhanced Features

**Redis Cache** (Optional)
- **Purpose**: Distributed caching for multi-user scenarios
- **Version**: Redis 7.0+
- **Installation**: Docker or native packages
- **Benefits**: Session sharing, real-time updates

**PostgreSQL** (Optional)
- **Purpose**: Enterprise-grade database for large deployments
- **Version**: PostgreSQL 14+
- **Installation**: Native packages or Docker
- **Benefits**: Better performance, advanced features

### 2. Development Enhancements

**Air Live Reload**
- **Purpose**: Automatic rebuild during development
- **Installation**: `go install github.com/air-verse/air@latest`
- **Configuration**: `.air.toml`
- **Benefits**: Faster development cycle

**Gotestsum**
- **Purpose**: Enhanced test runner with formatting
- **Installation**: `go install github.com/gotestyourself/gotest/tools/gotestsum@latest`
- **Benefits**: Better test output, JUnit XML

## Dependency Management Strategy

### 1. Version Pinning
```go
// go.mod will pin exact versions
module github.com/rice0649/fabric-lite

go 1.21

require (
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.17.0
    // ... other dependencies
)

// Use replace directives for development if needed
replace github.com/sashabaranov/go-openai => ../go-openai
```

### 2. Dependency Updates
```bash
# Update all dependencies
go get -u ./...

# Update specific dependency
go get -u github.com/spf13/cobra

# Check for outdated dependencies
go list -u -m all
```

### 3. Vulnerability Management
```bash
# Check for known vulnerabilities
go list -json -m all | nancy sleuth

# Or use govulncheck (Go 1.22+)
govulncheck ./...
```

## Risk Assessment for Dependencies

### 1. Critical Dependencies
- **Go Runtime**: Single point of failure, must maintain compatibility
- **Cloud AI APIs**: External dependencies with availability risks
- **Node.js CLI Tools**: Third-party tools with separate update cycles

### 2. Mitigation Strategies
- **Multiple Provider Support**: Failover between AI providers
- **Local Model Support**: Ollama as backup for cloud outages
- **Graceful Degradation**: Continue operation with reduced functionality
- **Caching**: Reduce dependency on external services

### 3. Monitoring Dependencies
- **Health Checks**: Regular ping of external services
- **Performance Metrics**: Track API response times and success rates
- **Cost Monitoring**: Track usage and costs across providers
- **Security Scanning**: Regular dependency vulnerability scanning

---

**Dependency Status**: All identified and documented  
**Implementation Ready**: Dependencies have clear acquisition paths  
**Risk Level**: Medium (managed through redundancy and monitoring)  
**Review Schedule**: Quarterly dependency audit and update