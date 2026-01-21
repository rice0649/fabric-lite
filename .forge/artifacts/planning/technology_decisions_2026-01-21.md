# Technology Decisions & Dependencies - Fabric-Lite Enhancement

## Technology Stack Rationale

### Core Programming Language: Go

**Why Go?**
- **Performance**: Compiled language with excellent concurrency support
- **Concurrency**: Goroutines and channels for multi-agent coordination
- **Portability**: Single binary deployment across all platforms
- **Ecosystem**: Rich ecosystem for CLI tools and HTTP servers
- **Memory Efficiency**: Low memory footprint for background processing
- **Standard Library**: Comprehensive libraries for networking, file I/O, JSON

**Alternatives Considered**:
- **Rust**: Better performance but steeper learning curve
- **Node.js**: Better AI/ML ecosystem but higher resource usage
- **Python**: Excellent AI libraries but slower for CLI operations

### Concurrency Model

**Goroutines + Channels Pattern**
```go
// Agent orchestration using goroutines
func (am *AgentManager) StartParallelAgents(agents []AgentSpec) error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    errChan := make(chan error, len(agents))
    var wg sync.WaitGroup
    
    for _, spec := range agents {
        wg.Add(1)
        go func(spec AgentSpec) {
            defer wg.Done()
            if err := am.startAgent(ctx, spec); err != nil {
                errChan <- err
            }
        }(spec)
    }
    
    // Wait for completion or errors
    go func() {
        wg.Wait()
        close(errChan)
    }()
    
    return <-errChan
}
```

**Benefits**:
- Lightweight goroutines (2KB stack vs 2MB threads)
- Built-in race condition detection
- Channel-based communication prevents shared memory issues
- Context-based cancellation and timeouts

### Storage Architecture

#### File System Design
```
.forge/
├── agents/
│   ├── configs/          # Agent configurations
│   ├── logs/           # Agent execution logs
│   └── status/         # Real-time status files
├── artifacts/
│   ├── discovery/       # Phase-specific artifacts
│   ├── planning/
│   ├── implementation/
│   └── testing/
├── cache/
│   ├── patterns/        # Compiled pattern cache
│   └── models/         # Model metadata cache
└── config/
    ├── forge.yaml       # Main configuration
    ├── providers.yaml    # Provider settings
    └── workflows.yaml   # Workflow definitions
```

#### Database Choices

**SQLite for Metadata**
- **Why SQLite?**: Embedded, zero-config, ACID compliant
- **Use Cases**: Session management, artifact indexing, metrics
- **Implementation**: `gorm.io/gorm` with SQLite driver

**Redis for Caching (Future)**
- **Why Redis?**: Fast in-memory caching, Pub/Sub
- **Use Cases**: Pattern cache, session state, real-time metrics
- **Implementation**: Optional dependency for enhanced performance

### API Integration Strategy

#### Provider Abstraction Layer
```go
type ProviderRegistry struct {
    providers map[string]ProviderFactory
    default   string
}

type ProviderFactory interface {
    Create(config ProviderConfig) (Provider, error)
    Validate(config ProviderConfig) error
    GetModels() []Model
}

// Provider implementations
var providerRegistry = map[string]ProviderFactory{
    "openai":    &OpenAIProviderFactory{},
    "anthropic": &AnthropicProviderFactory{},
    "google":    &GoogleProviderFactory{},
    "ollama":    &OllamaProviderFactory{},
    "cli":       &CLIProviderFactory{},
}
```

#### CLI Provider Integration
**CLI Wrapper Pattern**:
```go
type CLIProvider struct {
    command  string
    args     []string
    env      []string
    timeout  time.Duration
    logger   *zap.Logger
}

func (cp *CLIProvider) Execute(ctx context.Context, input string) (*GenerationResponse, error) {
    cmd := exec.CommandContext(ctx, cp.command, cp.args...)
    cmd.Env = append(os.Environ(), cp.env...)
    
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, err
    }
    
    go func() {
        defer stdin.Close()
        io.WriteString(stdin, input)
    }()
    
    output, err := cmd.CombinedOutput()
    return cp.parseResponse(output), err
}
```

## Dependencies Analysis

### Core Go Dependencies

#### Essential Libraries
```go
// CLI Framework
github.com/spf13/cobra v1.8.0          // Command-line interface
github.com/spf13/viper v1.17.0        // Configuration management

// Logging
go.uber.org/zap v1.26.0              // Structured logging
github.com/lmittmann/watermill v1.2.0   // Event streaming

// Database
gorm.io/gorm v1.25.5                  // ORM framework
gorm.io/driver/sqlite v1.5.4         // SQLite driver

// HTTP Clients
github.com/sashabaranov/go-openai v1.17.9 // OpenAI client
github.com/anthropics/anthropic-go v0.4.0 // Anthropic client
cloud.google.com/go/ai v0.24.0        // Google AI client

// Configuration
gopkg.in/yaml.v3 v3.0.1              // YAML parsing
github.com/kelseyhightower/envconfig v1.4.0 // Environment variables

// Utilities
github.com/google/uuid v1.4.0          // UUID generation
github.com/pkg/errors v0.9.1            // Error handling
github.com/stretchr/testify v1.8.4       // Testing framework
```

#### Performance & Monitoring
```go
// Metrics
github.com/prometheus/client_golang v1.17.0 // Prometheus metrics
go.opentelemetry.io/otel v1.21.0         // OpenTelemetry

// Performance
github.com/valyala/fasthttp v1.51.0     // Fast HTTP (future API)
github.com/klauspost/compress v1.17.0      // Compression
```

### External Dependencies

#### AI Provider Services

**Cloud Providers** (Required)
- **OpenAI**: GPT-4, o3-mini, Codex
  - API Cost: ~$0.03/1K tokens (GPT-4o)
  - Rate Limits: 10K requests/minute
  - Latency: 2-5 seconds average
  
- **Anthropic**: Claude 3.5 Sonnet
  - API Cost: ~$0.015/1K tokens
  - Rate Limits: 5K requests/minute
  - Latency: 3-7 seconds average

- **Google**: Gemini 2.0 Flash
  - API Cost: ~$0.0075/1K tokens
  - Rate Limits: 15K requests/minute
  - Latency: 1-3 seconds average

**Local Providers** (Optional)
- **Ollama**: Free local models
  - Models: Llama 3.1 (8B), Mistral, CodeLlama
  - Hardware: 8GB+ RAM recommended
  - Performance: 5-15 seconds per request

#### CLI Tools Integration

**Required CLI Tools**
```bash
# Core CLI providers
gemini         # Google Gemini CLI (npm install @anthropic/gemini-cli)
codex          # OpenAI Codex CLI (npm install @openai/codex-cli)
opencode       # OpenCode CLI (go install github.com/opencode/opencode)
cursor         # Cursor CLI (npm install cursor-cli)

# Browser automation
agent-browser  # Vercel agent-browser (npm install agent-browser)

# Optional local AI
ollama         # Local model server (curl -fsSL https://ollama.ai/install.sh)
```

#### Development Tools

**Build & Development**
```bash
# Go tooling
go version 1.21+                    # Go runtime
golangci-lint                        # Linting (v1.54+)
air                                  # Live reload (go install github.com/cosmtrek/air)

# Testing
gotestsum                           # Test runner
go-coverage                        # Coverage reporting
```

**Containerization & Deployment**
```bash
docker       # Container runtime
docker-compose # Multi-container orchestration
kubectl      # Kubernetes (optional)
```

## Security Considerations

### API Key Management

**Environment Variable Strategy**
```go
type SecureConfig struct {
    OpenAIKey    string `env:"OPENAI_API_KEY,required"`
    AnthropicKey string `env:"ANTHROPIC_API_KEY,required"`
    GoogleKey    string `env:"GOOGLE_API_KEY,required"`
}

func LoadSecureConfig() (*SecureConfig, error) {
    config := &SecureConfig{}
    if err := envconfig.Process("", config); err != nil {
        return nil, fmt.Errorf("failed to load secure config: %w", err)
    }
    return config, nil
}
```

**Key Rotation Strategy**
- **Rotation Frequency**: Every 30 days
- **Overlap Period**: 7 days during rotation
- **Fallback**: Multiple keys per provider
- **Monitoring**: Key usage and anomaly detection

### Data Privacy

**Local Processing Options**
```go
type PrivacyConfig struct {
    AllowCloudProviders bool     `yaml:"allow_cloud_providers"`
    SensitivePatterns []string  `yaml:"sensitive_patterns"`
    LocalModelPreference bool    `yaml:"local_model_preference"`
}

func (pc *PrivacyConfig) ShouldUseLocal(content string) bool {
    for _, pattern := range pc.SensitivePatterns {
        if matched, _ := filepath.Match(pattern, content); matched {
            return true
        }
    }
    return pc.LocalModelPreference
}
```

### Network Security

**TLS Configuration**
```go
type HTTPClient struct {
    client *http.Client
}

func NewSecureClient() *HTTPClient {
    return &HTTPClient{
        client: &http.Client{
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{
                    MinVersion: tls.VersionTLS12,
                    CipherSuites: []uint16{
                        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
                    },
                },
                Timeout: 30 * time.Second,
            },
        },
    }
}
```

## Performance Considerations

### Resource Management

**Memory Limits per Agent**
```go
type ResourcePool struct {
    maxMemory    int64         // bytes
    maxCPU       float64       // percentage
    activeAgents map[string]*ResourceUsage
}

type ResourceUsage struct {
    MemoryBytes int64         `json:"memory_bytes"`
    CPUPercent float64       `json:"cpu_percent"`
    Goroutines int            `json:"goroutines"`
}

func (rp *ResourcePool) AllocateAgentResources(agentID string) error {
    if rp.getAvailableMemory() < minAgentMemory {
        return ErrInsufficientMemory
    }
    if rp.getAvailableCPU() < minAgentCPU {
        return ErrInsufficientCPU
    }
    // Allocate resources...
    return nil
}
```

**Performance Targets**
- **Agent Startup**: <2 seconds
- **Message Latency**: <100ms between agents
- **Memory Usage**: <1GB for 4 concurrent agents
- **CPU Usage**: <50% of available cores
- **API Response Time**: <10 seconds for complex requests

### Caching Strategy

**Multi-Level Caching**
```go
type CacheManager struct {
    l1Cache *sync.Map          // In-memory for current session
    l2Cache *badger.DB         // Persistent cache across sessions
    l3Cache *redis.Client      // Optional distributed cache
}

func (cm *CacheManager) Get(key string) (interface{}, error) {
    // L1: In-memory
    if value, ok := cm.l1Cache.Load(key); ok {
        return value, nil
    }
    
    // L2: Persistent storage
    if value, err := cm.l2Cache.Get([]byte(key)); err == nil {
        cm.l1Cache.Store(key, value)
        return value, nil
    }
    
    // L3: Distributed cache (optional)
    if cm.l3Cache != nil {
        if value, err := cm.l3Cache.Get(key).Result(); err == nil {
            return value, nil
        }
    }
    
    return nil, ErrCacheMiss
}
```

## Development Workflow

### Project Structure
```
fabric-lite/
├── cmd/                    # CLI entry points
│   └── fabric-lite/
├── internal/               # Internal application code
│   ├── agent/            # Agent management
│   ├── communication/    # Inter-agent communication
│   ├── config/          # Configuration handling
│   ├── pattern/         # Pattern system
│   ├── provider/        # AI provider implementations
│   ├── workflow/        # Workflow orchestration
│   └── testing/         # Testing utilities
├── pkg/                   # Public library code
├── patterns/              # Built-in patterns
├── configs/               # Configuration templates
├── scripts/               # Build and deployment scripts
├── docs/                  # Documentation
├── tests/                 # Integration tests
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── docker-compose.yml
```

### Build & Release Process

**Makefile Targets**
```makefile
.PHONY: build test lint clean release docker

# Build for current platform
build:
	go build -ldflags "-X main.Version=$(VERSION)" -o bin/fabric-lite ./cmd/fabric-lite

# Cross-platform builds
release:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o bin/fabric-lite-linux-amd64 ./cmd/fabric-lite
	GOOS=darwin GOARCH=amd64 go build -o bin/fabric-lite-darwin-amd64 ./cmd/fabric-lite
	GOOS=windows GOARCH=amd64 go build -o bin/fabric-lite-windows-amd64.exe ./cmd/fabric-lite

# Testing
test:
	go test -v ./...
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Linting
lint:
	golangci-lint run

# Docker
docker:
	docker build -t fabric-lite:$(VERSION) .
	docker tag fabric-lite:$(VERSION) fabric-lite:latest
```

---

**Technology Status**: Decisions finalized  
**Implementation Ready**: All dependencies identified and documented  
**Next Step**: Begin core orchestration component development  
**Review Cadence**: Technology decisions reviewed quarterly