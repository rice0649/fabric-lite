See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# IDENTITY and PURPOSE

You are an expert software architect and implementation engineer. Your role is to design and implement a comprehensive solution for completing the fabric-lite AI augmentation framework and optimizing the forge orchestration system based on extensive analysis from Codex and Gemini expert systems.

## STEPS

1. Analyze the current state and gaps in fabric-lite and forge systems
2. Design optimized architecture using config-driven provider pattern
3. Plan phased implementation approach with immediate priorities
4. Define implementation tasks with specific code deliverables
5. Document integration strategy and success metrics

# OUTPUT INSTRUCTIONS

Produce output in following sections:

## Architecture Overview
High-level description of system architecture with a text-based diagram.

## Components
For each component:
- **Name**: Component name
- **Responsibility**: What it does
- **Technology**: Suggested technology stack
- **Interfaces**: How it communicates with other components

## Technology Decisions
| Decision | Choice | Rationale | Alternatives Considered |
|----------|--------|-----------|------------------------|

## Implementation Phases
For each phase:
- **Phase X: [Name]** (Priority: X)
  - Duration: X days
  - Key Deliverables: [list]
  - Implementation Tasks: [detailed tasks]
  - Validation Criteria: [how to verify completion]

## Data Flow
Describe how data flows through the system during pattern execution and tool orchestration.

## Risks and Mitigations
| Risk | Impact | Mitigation Strategy |
|------|--------|---------------------|

## Dependencies
External services and libraries required for successful implementation.

## Immediate Action Plan
Detailed tasks to execute in the next 48 hours with specific file locations and code snippets.

## Integration Strategy
How fabric-lite and forge will integrate and share components.

## Success Metrics
Specific, measurable criteria for v1.0 completion.

# INPUT

## CURRENT STATE ANALYSIS

Based on comprehensive analysis from multiple AI systems (Codex, Gemini, and expert review), the following critical gaps have been identified:

### Fabric-Lite Critical Issues:
1. **Missing CLI Implementation**: `cmd/fabric-lite/main.go:20` contains TODO comment, no actual functionality
2. **No Provider Layer**: `internal/providers/` directory completely missing, no AI integrations
3. **Incomplete Pattern System**: Discovery exists but execution engine missing
4. **Unused Configuration**: `config/config.example.yaml` exists but not connected to runtime

### Forge Current State:
1. **Well-Structured**: Complete CLI with Cobra, comprehensive tool wrappers
2. **95% Complete**: All phases implemented except final deployment completion
3. **Sophisticated Architecture**: Phase management, state tracking, tool orchestration working
4. **Missing Tool Chaining**: Single tool per phase, no pipeline capability

### Key Architecture Decisions Validated:
1. **Config-Driven Provider System**: Preferred over Go plugins (Gemini recommendation)
2. **Generic HTTP Provider**: Single implementation for OpenAI, Anthropic, Ollama APIs
3. **Tool Chaining**: Sequential pipeline execution for forge phases
4. **Streaming Support**: Server-sent events for real-time pattern execution
5. **Fail-Fast Error Handling**: Stop on first error with configurable policies

## OPTIMIZED ARCHITECTURE OVERVIEW

```
┌─────────────────────────────────────────────────────────────┐
│                    FABRIC-LITE v1.0                   │
│              (Pattern Execution Engine)                 │
├─────────────────────────────────────────────────────────────┤
│ Components:                                         │
│ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐│
│ │   CLI       │ │  Providers   │ │  Executor   ││
│ │ - Cobra     │ │ - HTTP      │ │ - Patterns  ││
│ │ - Subcmds   │ │ - Executable │ │ - Streaming  ││
│ └─────────────┘ │ - Config    │ │ - Variables  ││
│                └─────────────┘ └─────────────┘│
├─────────────────────────────────────────────────────────────┤
│ Configuration:                                        │
│ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐│
│ │  Config     │ │  Sessions    │ │  Cache      ││
│ │ - YAML      │ │ - History    │ │ - RateLimit  ││
│ │ - Env Vars  │ │ - Context    │ │ - Pool       ││
│ └─────────────┘ └─────────────┘ └─────────────┘│
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                    FORGE v1.0                          │
│              (Orchestration Layer)                     │
├─────────────────────────────────────────────────────────────┤
│ Components:                                         │
│ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐│
│ │   Phases    │ │  Tools       │ │   Chains    ││
│ │ - 6 phases  │ │ - Gemini     │ │ - Pipeline   ││
│ │ - Checkpts  │ │ - Codex      │ │ - Flow      ││
│ │ - State     │ │ - OpenCode   │ │ - Context   ││
│ │ - Auto      │ │ - Fabric     │ │ - Error     ││
│ └─────────────┘ └─────────────┘ └─────────────┘│
├─────────────────────────────────────────────────────────────┤
│ Integration Points:                                   │
│ ✅ Fabric-Lite used as pattern engine              │
│ ✅ Shared configuration system                     │
│ ✅ Unified tool interface                       │
│ ✅ Cross-project session management              │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                 EXTERNAL SERVICES                        │
├─────────────────────────────────────────────────────────────┤
│ Providers:                                         │
│ • OpenAI API (gpt-4o, gpt-4o-mini)          │
│ • Anthropic API (claude-3-5-sonnet)             │
│ • Ollama Local (llama3.2, mistral)             │
│ • Google API (gemini-pro)                        │
│ • Custom Executables (user-defined scripts)          │
└─────────────────────────────────────────────────────────────┘
```

## COMPONENTS

### **Fabric-Lite Core Components**

#### **CLI Interface** 
- **Responsibility**: Command-line interface with Cobra framework, subcommands (run, list, config), argument parsing, help system
- **Technology**: Go with Cobra/Viper, configuration management
- **Interfaces**: Connects to Executor for pattern execution, loads Providers from config

#### **Provider System**
- **Responsibility**: Pluggable AI provider interface with HTTP and executable implementations, configuration-driven loading, API key management
- **Technology**: Go interfaces, YAML configuration, HTTP clients, process execution
- **Interfaces**: Provider interface, factory pattern, config struct binding

#### **Pattern Executor**
- **Responsibility**: Pattern loading, variable substitution, streaming execution, session management, response formatting
- **Technology**: Go file system, template processing, SSE streaming, context management
- **Interfaces**: Accepts providers, handles patterns, manages conversation state

### **Forge Orchestration Components**

#### **Phase Manager**
- **Responsibility**: Six development phases, checkpoint validation, state transitions, progress tracking
- **Technology**: Go state machine, YAML configuration, file system checkpoints
- **Interfaces**: Coordinates with Tools, manages Phase state, validates completion

#### **Tool Registry**
- **Responsibility**: Dynamic tool loading, execution wrapper, result aggregation, error handling
- **Technology**: Go reflection, factory pattern, configuration-driven discovery
- **Interfaces**: Implements Tool interface, passes execution context, returns results

#### **Chain Engine**
- **Responsibility**: Sequential tool execution, output-to-input passing, error policy enforcement
- **Technology**: Go pipelines, context propagation, fail-fast strategy
- **Interfaces**: Connects multiple Tools, manages data flow, handles rollbacks

## TECHNOLOGY DECISIONS

| Decision | Choice | Rationale | Alternatives Considered |
|----------|--------|-----------|------------------------|
| Provider Architecture | Config-driven factory pattern | Simple, secure, no Go plugin complexity | Go plugins (overkill), Hardcoded providers (inflexible) |
| HTTP Provider Implementation | Single generic HTTP client with configurable endpoints | Reduces code duplication, easy extensibility | Separate implementations per provider (more code) |
| Pattern Variable Syntax | `{{variable}}` substitution with validation | Familiar template syntax, easy parsing | Custom DSL (complex), No variables (limiting) |
| Tool Chaining | Sequential pipeline with fail-fast | Simple to implement, predictable behavior | Parallel execution (complex), Branching flows (overkill) |
| Configuration Format | YAML with environment variable substitution | Human-readable, standard, secure env vars | JSON (less readable), TOML (less common) |
| Error Handling | Fail-fast with configurable policies | Predictable, easy to debug, good UX | Continue-on-error (complex), Rollback (overkill) |
| Session Storage | File-based with JSON | Simple, persistent, easy debugging | Database (overkill), Memory-only (lost on restart) |

## IMPLEMENTATION PHASES

### **Phase 1: Fabric-Lite Foundation** (Priority: CRITICAL)
- **Duration**: 2 days
- **Key Deliverables**: Working CLI, provider system, basic pattern execution
- **Implementation Tasks**:
  1. Create `internal/providers/provider.go` interface and factory
  2. Implement `internal/providers/http_provider.go` with OpenAI integration
  3. Replace `cmd/fabric-lite/main.go` with Cobra implementation
  4. Create `internal/executor/pattern.go` for basic pattern execution
  5. Implement configuration loading from `~/.config/fabric-lite/providers.yaml`
  6. Add pattern discovery and listing functionality
- **Validation Criteria**: 
  - CLI accepts basic commands (run, list, config)
  - OpenAI provider can execute simple patterns
  - Configuration system loads providers from YAML
  - Pattern directory scanning works correctly

### **Phase 2: Enhanced Fabric-Lite** (Priority: HIGH)  
- **Duration**: 3 days
- **Key Deliverables**: Streaming, sessions, variable substitution, multiple providers
- **Implementation Tasks**:
  1. Add Anthropic and Ollama provider configurations
  2. Implement streaming with SSE in `internal/executor/streaming.go`
  3. Create session management in `internal/executor/session.go`
  4. Add variable substitution engine in `internal/executor/variables.go`
  5. Implement executable provider for custom scripts
  6. Add comprehensive error handling and retry logic
- **Validation Criteria**:
  - Streaming responses work for all providers
  - Session persistence across CLI runs
  - Variable substitution in patterns functions correctly
  - All major providers (OpenAI, Anthropic, Ollama) operational

### **Phase 3: Forge Tool Chaining** (Priority: HIGH)
- **Duration**: 2 days  
- **Key Deliverables**: Multi-tool phases, pipeline execution, output passing
- **Implementation Tasks**:
  1. Enhance `internal/tools/tool.go` interface for input/output passing
  2. Create `internal/core/chains.go` for pipeline execution
  3. Update phase configuration to support tool arrays
  4. Implement error policies (continue/stop/skip)
  5. Add chain execution logging and diagnostics
- **Validation Criteria**:
  - Phases can specify multiple tools in sequence
  - Output of tool A becomes input of tool B
  - Error handling respects configured policy
  - Chain execution is properly logged

### **Phase 4: Production Features** (Priority: MEDIUM)
- **Duration**: 3 days
- **Key Deliverables**: Security, performance, monitoring, testing
- **Implementation Tasks**:
  1. Implement API key encryption with system keyring
  2. Add rate limiting and connection pooling
  3. Create comprehensive test suite (>80% coverage)
  4. Add metrics collection and health checks
  5. Implement response caching with TTL
- **Validation Criteria**:
  - API keys encrypted at rest
  - Rate limits enforced per provider
  - Test suite passes with >80% coverage
  - Performance metrics collected
  - Caching reduces API calls appropriately

## DATA FLOW

### **Pattern Execution Flow (Fabric-Lite)**:
```
User Input (CLI) → Parse Command → Load Provider Config → Execute Pattern → Stream Response → Save Session
     │                │                    │                  │                   │
     ▼                ▼                    ▼                  ▼                   ▼
fabric-lite run     provider.yaml        pattern.md         HTTP POST         ~/.config/fabric-lite/
--pattern test      → HTTP Provider      → Variables       → OpenAI API       sessions/{timestamp}.json
--model gpt-4o      → Auth Headers      → Substitution     → SSE Stream         → CLI Output
--input "..."       → Request           → Compiled Prompt   → Real-time Display   → Save History
```

### **Tool Orchestration Flow (Forge)**:
```
forge phase implementation → Load Phase Config → Initialize Tools → Execute Chain → Validate Checkpoints → Update State
          │                        │                  │             │                  │
          ▼                        ▼                  ▼             ▼                  ▼
     Command Parser        tools.yaml        Tool A → Tool B → Result C    Check Files      .forge/state.yaml
     Phase Definition      → Tool Registry      → Pass Output   → Pass Output   → Write Artifacts   → Progress Update
     → tools: [codex,     → Execute()         → Execute()     → Execute()     → phase: completed
        fabric]            → Context            → Collect      → Aggregate    → artifacts/*   → next: testing
```

## RISKS AND MITIGATIONS

| Risk | Impact | Mitigation Strategy |
|------|--------|---------------------|
| Provider API Changes | Medium | Version pinning, adapter pattern for breaking changes |
| Configuration Errors | High | Validation, clear error messages, environment variable fallbacks |
| Tool Chain Failures | Medium | Fail-fast default, configurable error policies, partial result preservation |
| Session Data Loss | Low | Regular backups, file-based persistence with recovery |
| Performance Issues | Medium | Connection pooling, request caching, rate limiting |
| Security Vulnerabilities | High | API key encryption, input validation, no eval() usage |
| Provider Unavailability | Low | Graceful degradation, multiple provider fallbacks |

## DEPENDENCIES

### **External Services**
- **OpenAI API**: gpt-4o, gpt-4o-mini models
- **Anthropic API**: claude-3-5-sonnet-20241022 model  
- **Ollama Local**: llama3.2, mistral models
- **Custom Scripts**: User-defined executable tools

### **Go Libraries**
- github.com/spf13/cobra (CLI framework)
- github.com/spf13/viper (configuration management)  
- gopkg.in/yaml.v3 (YAML parsing)
- golang.org/x/net (HTTP requests, SSE streaming)

### **System Requirements**
- Go 1.21+ (already specified in go.mod)
- POSIX-compliant system (for shell script providers)
- ~/.config/fabric-lite/ and ~/.config/forge/ directories
- Internet access for cloud providers, localhost for Ollama

## IMMEDIATE ACTION PLAN (Next 48 Hours)

### **Priority 1: Fabric-Lite CLI Foundation**
```bash
# File: cmd/fabric-lite/main.go - COMPLETE REPLACEMENT
cat > cmd/fabric-lite/main.go << 'EOF'
package main

import (
    "fmt"
    "os"
    "github.com/rice0649/fabric-lite/internal/cli"
    "github.com/spf13/cobra"
)

var Version = "0.1.0"

func main() {
    rootCmd := cli.NewRootCmd(Version)
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
EOF
```

### **Priority 2: Provider Interface System**
```bash
# Directory: internal/providers/
mkdir -p internal/providers

# File: internal/providers/provider.go
cat > internal/providers/provider.go << 'EOF'
package providers

import (
    "context"
    "time"
)

type CompletionRequest struct {
    System    string            `json:"system"`
    Prompt    string            `json:"prompt"`
    Model     string            `json:"model"`
    MaxTokens int              `json:"max_tokens,omitempty"`
    Stream    bool              `json:"stream,omitempty"`
    Options   map[string]any    `json:"options,omitempty"`
}

type CompletionResponse struct {
    Content   string    `json:"content"`
    Model     string    `json:"model"`
    Tokens    int       `json:"tokens"`
    Duration  time.Duration `json:"duration"`
    Error     error     `json:"-"`
}

type Provider interface {
    Name() string
    Execute(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)
    IsAvailable() bool
    GetModels() []string
}

type ProviderConfig struct {
    Name   string                 `yaml:"name"`
    Type   string                 `yaml:"type"`
    Config map[string]interface{}   `yaml:"config"`
}
EOF
```

### **Priority 3: HTTP Provider Implementation**
```bash
# File: internal/providers/http_provider.go (First 100 lines)
cat > internal/providers/http_provider.go << 'EOF'
package providers

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strconv"
    "time"
)

type HTTPProvider struct {
    name     string
    endpoint string
    apiKey   string
    headers  map[string]string
    client   *http.Client
}

func NewHTTPProvider(name string, config map[string]interface{}) (*HTTPProvider, error) {
    endpoint, _ := config["endpoint"].(string)
    apiKey := os.Getenv(config["api_key_env"].(string))
    
    headers := map[string]string{
        "Content-Type": "application/json",
        "Authorization": "Bearer " + apiKey,
    }
    
    return &HTTPProvider{
        name:     name,
        endpoint: endpoint,
        apiKey:   apiKey,
        headers:  headers,
        client:   &http.Client{Timeout: 30 * time.Second},
    }, nil
}

func (p *HTTPProvider) Name() string {
    return p.name
}

func (p *HTTPProvider) IsAvailable() bool {
    return p.apiKey != "" && p.endpoint != ""
}

func (p *HTTPProvider) GetModels() []string {
    // TODO: Implement model discovery per provider
    return []string{"gpt-4o", "gpt-4o-mini"}
}
EOF
```

### **Priority 4: Pattern Executor Basic**
```bash
# File: internal/executor/pattern.go (Core structure)
mkdir -p internal/executor
cat > internal/executor/pattern.go << 'EOF'
package executor

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    
    "github.com/rice0649/fabric-lite/internal/providers"
)

type PatternExecutor struct {
    providers map[string]providers.Provider
    patternsDir string
}

func NewPatternExecutor() *PatternExecutor {
    return &PatternExecutor{
        patternsDir: filepath.Join(os.Getenv("HOME"), ".config", "fabric-lite", "patterns"),
        providers: make(map[string]providers.Provider),
    }
}

func (e *PatternExecutor) Execute(ctx context.Context, patternName, input string, providerName string) (*providers.CompletionResponse, error) {
    provider, ok := e.providers[providerName]
    if !ok {
        return nil, fmt.Errorf("provider not found: %s", providerName)
    }
    
    pattern, err := e.loadPattern(patternName)
    if err != nil {
        return nil, fmt.Errorf("failed to load pattern: %w", err)
    }
    
    request := providers.CompletionRequest{
        System: pattern.System,
        Prompt: pattern.User + "\n\nInput:\n" + input,
        Model:  providerName, // TODO: Get from provider
    }
    
    return provider.Execute(ctx, request)
}
EOF
```

### **Priority 5: Configuration Integration**
```bash
# File: internal/providers/config.go
cat > internal/providers/config.go << 'EOF'
package providers

import (
    "fmt"
    "os"
    "path/filepath"
    
    "gopkg.in/yaml.v3"
)

type Config struct {
    Providers []ProviderConfig `yaml:"providers"`
}

func LoadConfig(configPath string) (*Config, error) {
    if configPath == "" {
        configPath = filepath.Join(os.Getenv("HOME"), ".config", "fabric-lite", "providers.yaml")
    }
    
    data, err := os.ReadFile(configPath)
    if err != nil {
        if os.IsNotExist(err) {
            return &Config{}, nil // Return empty config if file doesn't exist
        }
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    var config Config
    err = yaml.Unmarshal(data, &config)
    return &config, err
}
EOF
```

## INTEGRATION STRATEGY

### **Component Sharing Architecture**
```
┌─────────────────────────────────────────────────────────────┐
│              SHARED INFRASTRUCTURE                    │
├─────────────────────────────────────────────────────────────┤
│ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐│
│ │  Config     │ │  Logging    │ │  Errors     ││
│ │ - YAML      │ │ - Structured │ │ - Types     ││
│ │ - Env Vars  │ │ - Context    │ │ - Recovery   ││
│ │ - Validation│ │ - Levels     │ │ - Policies   ││
│ └─────────────┘ └─────────────┘ └─────────────┘│
└─────────────────────────────────────────────────────────────┘
          ▼                ▼                ▼
┌───────────────┬───────────────┬───────────────┐
│  FABRIC-LITE  │     FORGE     │  EXTERNAL     │
│ Uses Config   │ Uses Config   │  API PROVIDERS │
│ Uses Logging  │ Uses Logging  │  TOOLS         │
│ Uses Errors   │ Uses Errors   │  SERVICES       │
└───────────────┴───────────────┴───────────────┘
```

### **Integration Points**
1. **Configuration System**: Shared YAML structure with environment variable substitution
2. **Error Handling**: Unified error types and recovery strategies  
3. **Logging Framework**: Structured logging with configurable levels
4. **Provider Registry**: Fabric-Lite providers used by Forge tool wrapper
5. **Session Management**: Cross-project conversation and state persistence
6. **Tool Interface**: Common abstraction for both systems

## SUCCESS METRICS

### **Fabric-Lite v1.0 Completion Criteria**
- [ ] **CLI Functionality**: All commands (run, list, config, version) working
- [ ] **Provider Coverage**: OpenAI, Anthropic, Ollama, and executable providers operational  
- [ ] **Pattern Execution**: Variable substitution, streaming, session management working
- [ ] **Configuration**: YAML configuration loading with validation complete
- [ ] **Test Coverage**: Unit tests covering >80% of codebase
- [ ] **Performance**: Response time <5 seconds for simple patterns
- [ ] **Documentation**: User guide and API reference complete
- [ ] **Integration**: Forge successfully uses Fabric-Lite as pattern engine

### **Forge v1.0 Completion Criteria**
- [ ] **Phase Completion**: All 6 phases implemented and tested
- [ ] **Tool Chaining**: Multi-tool pipeline execution working
- [ ] **Dynamic Loading**: Provider and tool configuration from YAML
- [ ] **Session Management**: Cross-project state persistence functional
- [ ] **Error Handling**: Fail-fast with configurable policies
- [ ] **Deployment Ready**: Production configuration and monitoring
- [ ] **Documentation**: Complete user and developer documentation
- [ ] **Integration Tests**: End-to-end workflow validation

## FINAL RECOMMENDATION

**Immediate Priority**: Implement Phase 1 (Fabric-Lite Foundation) within 48 hours. This establishes the core architecture that both projects depend on. The config-driven provider pattern recommended by Gemini provides optimal extensibility without complexity, while the tool chaining capability will significantly enhance Forge's automation capabilities.

**Strategic Approach**: Complete Fabric-Lite core first (Phase 1), then immediately use it to finish Forge deployment phase. This validates the integration and provides immediate value to users.