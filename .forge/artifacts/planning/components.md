# Component Breakdown - Fabric-Lite Multi-Agent System

## 1. Core Orchestration Components

### 1.1 Agent Manager (`internal/agent/manager.go`)
**Responsibility**: Lifecycle management of AI agents

```go
type AgentManager struct {
    agents      map[string]*Agent
    resources   *ResourcePool
    logger      *zap.Logger
    status      chan AgentStatus
    workspace   string
}

type Agent struct {
    ID          string
    Name        string
    Provider    string
    Model       string
    Config      AgentConfig
    Process     *os.Process
    Status      AgentStatus
    StartTime   time.Time
    Metrics     *AgentMetrics
}

type AgentConfig struct {
    Headless       bool          `yaml:"headless"`
    Background     bool          `yaml:"background"`
    MemoryLimit    string        `yaml:"memory_limit"`
    Timeout        time.Duration `yaml:"timeout"`
    RetryCount     int           `yaml:"retry_count"`
    Environment    []string      `yaml:"environment"`
    AllowedTools   []string      `yaml:"allowed_tools"`
}

type AgentMetrics struct {
    TasksCompleted    int64         `json:"tasks_completed"`
    AverageLatency   time.Duration `json:"average_latency"`
    ErrorCount       int64         `json:"error_count"`
    LastActivity     time.Time     `json:"last_activity"`
    ResourceUsage    ResourceUsage  `json:"resource_usage"`
}
```

**Key Functions**:
- `StartAgent(config AgentConfig) (*Agent, error)`
- `StopAgent(agentID string) error`
- `MonitorAgent(agentID string) (*AgentStatus, error)`
- `ListAgents() ([]*Agent, error)`
- `ScaleAgentPool(desiredCount int) error`

### 1.2 Communication Layer (`internal/communication/layer.go`)
**Responsibility**: File-based message passing and coordination

```go
type CommunicationLayer struct {
    workspace   string
    queueDir    string
    statusDir   string
    artifactDir string
    logger      *zap.Logger
}

type Message struct {
    ID          string            `json:"id"`
    From        string            `json:"from"`
    To          string            `json:"to"`
    Type        MessageType       `json:"type"`
    Payload     interface{}       `json:"payload"`
    Timestamp   time.Time         `json:"timestamp"`
    Priority    Priority         `json:"priority"`
    Metadata    map[string]string `json:"metadata"`
}

type MessageType string
const (
    MessageTypeTask     MessageType = "task"
    MessageTypeStatus   MessageType = "status"
    MessageTypeResult   MessageType = "result"
    MessageTypeError   MessageType = "error"
    MessageTypeControl  MessageType = "control"
)

type Status struct {
    AgentID     string                 `json:"agent_id"`
    State       AgentState             `json:"state"`
    Progress    float64                `json:"progress"`
    CurrentTask string                 `json:"current_task"`
    Health      HealthStatus           `json:"health"`
    LastUpdate  time.Time              `json:"last_update"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type Artifact struct {
    ID          string            `json:"id"`
    Type        ArtifactType      `json:"type"`
    Phase       string            `json:"phase"`
    Agent       string            `json:"agent"`
    Content     []byte            `json:"content"`
    Format      string            `json:"format"`
    CreatedAt   time.Time         `json:"created_at"`
    Metadata    map[string]string `json:"metadata"`
}
```

**Key Functions**:
- `SendMessage(msg Message) error`
- `ReceiveMessages(agentID string) ([]Message, error)`
- `UpdateStatus(status Status) error`
- `GetStatus(agentID string) (*Status, error)`
- `ShareArtifact(artifact Artifact) error`
- `ListArtifacts(phase string) ([]Artifact, error)`

### 1.3 Workflow Engine (`internal/workflow/engine.go`)
**Responsibility**: Phase-based execution with checkpoint validation

```go
type WorkflowEngine struct {
    phases      map[string]*Phase
    currentPhase *Phase
    agentMgr    *AgentManager
    commLayer   *CommunicationLayer
    logger      *zap.Logger
    status      WorkflowStatus
}

type Phase struct {
    Name         string        `yaml:"name"`
    Description  string        `yaml:"description"`
    Agents       []AgentSpec   `yaml:"agents"`
    Checkpoints  []Checkpoint  `yaml:"checkpoints"`
    Timeout      time.Duration `yaml:"timeout"`
    Dependencies []string      `yaml:"dependencies"`
    Artifacts    []string      `yaml:"expected_artifacts"`
}

type AgentSpec struct {
    Name     string `yaml:"name"`
    Provider string `yaml:"provider"`
    Model    string `yaml:"model"`
    Role     string `yaml:"role"`
    Priority int    `yaml:"priority"`
}

type Checkpoint struct {
    Name        string   `yaml:"name"`
    Description string   `yaml:"description"`
    Required    []string `yaml:"required_files"`
    Validation string   `yaml:"validation_command"`
    Timeout    time.Duration `yaml:"timeout"`
}
```

**Key Functions**:
- `ExecutePhase(phaseName string) error`
- `ValidateCheckpoints(phaseName string) (*ValidationResult, error)`
- `GetCurrentPhase() (*Phase, error)`
- `GetProgress() *WorkflowProgress`
- `PauseWorkflow() error`
- `ResumeWorkflow() error`

## 2. Provider System Components

### 2.1 Unified Provider Interface (`internal/provider/interface.go`)
**Responsibility**: Abstraction for all AI providers

```go
type Provider interface {
    Initialize(config ProviderConfig) error
    Generate(ctx context.Context, request GenerationRequest) (*GenerationResponse, error)
    Stream(ctx context.Context, request GenerationRequest) (<-chan StreamChunk, error)
    GetModels() []Model
    IsHealthy() bool
    GetUsage() *Usage
    Close() error
}

type GenerationRequest struct {
    Prompt      string            `json:"prompt"`
    Model       string            `json:"model"`
    MaxTokens   int               `json:"max_tokens"`
    Temperature float32           `json:"temperature"`
    Context     []Message         `json:"context"`
    Tools       []Tool            `json:"tools"`
    SystemPrompt string            `json:"system_prompt"`
    Metadata    map[string]string `json:"metadata"`
}

type GenerationResponse struct {
    Content      string            `json:"content"`
    TokenUsage   *TokenUsage       `json:"token_usage"`
    Model        string            `json:"model"`
    FinishReason string            `json:"finish_reason"`
    Metadata    map[string]string `json:"metadata"`
}

type Tool interface {
    Name        string `json:"name"`
    Description string `json:"description"`
    Execute(ctx context.Context, input interface{}) (interface{}, error)
}
```

### 2.2 Provider Implementations

#### OpenAI Provider (`internal/provider/openai/provider.go`)
```go
type OpenAIProvider struct {
    client   *openai.Client
    config   OpenAIConfig
    models   []Model
    logger   *zap.Logger
}

type OpenAIConfig struct {
    APIKey      string `yaml:"api_key"`
    BaseURL     string `yaml:"base_url"`
    Organization string `yaml:"organization"`
    Models      []string `yaml:"models"`
}
```

**Models Supported**:
- `gpt-4o`: High-quality general tasks
- `gpt-4o-mini`: Fast, cost-effective
- `o3-mini`: Advanced reasoning
- `codex`: Code generation specialist

#### Anthropic Provider (`internal/provider/anthropic/provider.go`)
```go
type AnthropicProvider struct {
    client   *anthropic.Client
    config   AnthropicConfig
    models   []Model
    logger   *zap.Logger
}

type AnthropicConfig struct {
    APIKey  string   `yaml:"api_key"`
    BaseURL string   `yaml:"base_url"`
    Models  []string `yaml:"models"`
}
```

**Models Supported**:
- `claude-3-5-sonnet-20241022`: High-performance
- `claude-3-haiku-20240307`: Fast, cost-effective

#### Google Provider (`internal/provider/google/provider.go`)
```go
type GoogleProvider struct {
    client   *genai.Client
    config   GoogleConfig
    models   []Model
    logger   *zap.Logger
}

type GoogleConfig struct {
    APIKey      string   `yaml:"api_key"`
    ProjectID   string   `yaml:"project_id"`
    SearchEnabled bool    `yaml:"search_enabled"`
    Models      []string `yaml:"models"`
}
```

**Models Supported**:
- `gemini-2.0-flash-exp`: Fast with 1M context
- `gemini-1.5-pro`: High-quality reasoning

#### Ollama Provider (`internal/provider/ollama/provider.go`)
```go
type OllamaProvider struct {
    client   *http.Client
    config   OllamaConfig
    models   []Model
    logger   *zap.Logger
}

type OllamaConfig struct {
    BaseURL     string        `yaml:"base_url"`
    Models      []string      `yaml:"models"`
    Timeout     time.Duration `yaml:"timeout"`
    PullModels  bool          `yaml:"pull_models"`
}
```

**Models Supported**:
- `llama3.1:8b`: General purpose
- `mistral`: Fast reasoning
- `codellama`: Code specialized

### 2.3 CLI Provider Integration (`internal/provider/cli/`)
**Responsibility**: Integration with external AI CLI tools

```go
type CLIProvider struct {
    name     string
    command  string
    args     []string
    config   CLIConfig
    logger   *zap.Logger
}

type CLIConfig struct {
    Path          string        `yaml:"path"`
    DefaultArgs   []string      `yaml:"default_args"`
    Timeout       time.Duration `yaml:"timeout"`
    RetryCount    int           `yaml:"retry_count"`
    OutputFormat  string        `yaml:"output_format"`
}

// Specific CLI implementations
type GeminiCLI struct{ CLIProvider }
type CodexCLI struct{ CLIProvider }
type OpenCodeCLI struct{ CLIProvider }
type CursorCLI struct{ CLIProvider }
```

## 3. Pattern System Components

### 3.1 Pattern Registry (`internal/pattern/registry.go`)
**Responsibility**: Pattern management and versioning

```go
type PatternRegistry struct {
    patterns map[string]*Pattern
    storage  PatternStorage
    logger   *zap.Logger
}

type Pattern struct {
    ID            string            `yaml:"id"`
    Name          string            `yaml:"name"`
    Version       string            `yaml:"version"`
    Category      string            `yaml:"category"`
    Description   string            `yaml:"description"`
    Template      string            `yaml:"template"`
    Variables     []PatternVariable `yaml:"variables"`
    Dependencies  []string          `yaml:"dependencies"`
    Models        []string          `yaml:"recommended_models"`
    Author        string            `yaml:"author"`
    Tags          []string          `yaml:"tags"`
    Quality       float32           `yaml:"quality_score"`
    CreatedAt     time.Time         `yaml:"created_at"`
    UpdatedAt     time.Time         `yaml:"updated_at"`
}

type PatternVariable struct {
    Name         string `yaml:"name"`
    Type         string `yaml:"type"`
    Required     bool   `yaml:"required"`
    Default      string `yaml:"default"`
    Description  string `yaml:"description"`
    Validation  string `yaml:"validation"`
}
```

**Pattern Categories**:
- **Development**: `code-generation`, `code-review`, `testing`, `refactoring`
- **Documentation**: `api-docs`, `readme`, `changelog`, `technical-writing`
- **Operations**: `infrastructure`, `deployment`, `monitoring`, `security`
- **Research**: `analysis`, `summarization`, `reporting`, `market-research`

### 3.2 Pattern Engine (`internal/pattern/engine.go`)
**Responsibility**: Pattern execution and adaptation

```go
type PatternEngine struct {
    registry  *PatternRegistry
    providers map[string]Provider
    cache     *PatternCache
    logger    *zap.Logger
}

type ExecutionContext struct {
    Pattern      *Pattern
    Variables    map[string]interface{}
    Provider     string
    Model        string
    Context      []Message
    Tools        []Tool
    Artifacts    []Artifact
    Phase        string
    UserID       string
    SessionID    string
}

type PatternCache struct {
    cache      map[string]*CachedPattern
    ttl        time.Duration
    maxSize    int
    currentSize int
}
```

## 4. CLI Enhancement Components

### 4.1 Enhanced Commands (`cmd/fabric-lite/`)
**New Commands**:
- `agent`: Agent management commands
- `workflow`: Workflow management commands
- `orchestrate`: Start multi-agent orchestration
- `pattern`: Pattern management commands
- `provider`: Provider configuration commands

```go
// Agent commands
var agentCmd = &cobra.Command{
    Use:   "agent",
    Short: "Manage AI agents",
    Long:  "Start, stop, and monitor AI agents",
}

// Workflow commands
var workflowCmd = &cobra.Command{
    Use:   "workflow",
    Short: "Manage workflows",
    Long:  "Execute and monitor development workflows",
}

// Orchestrate command
var orchestrateCmd = &cobra.Command{
    Use:   "orchestrate",
    Short: "Start multi-agent orchestration",
    Long:  "Execute multi-agent workflows in headless mode",
}
```

### 4.2 Configuration System (`internal/config/`)
**Enhanced Configuration Structure**:

```yaml
# fabric-lite configuration
version: "2.0"

# Multi-agent orchestration
orchestration:
  max_concurrent_agents: 4
  default_timeout: 300s
  resource_limits:
    memory: "4GB"
    cpu: "50%"
  logging:
    level: "info"
    format: "json"

# Provider configurations
providers:
  openai:
    api_key: "${OPENAI_API_KEY}"
    models: ["gpt-4o", "gpt-4o-mini", "o3-mini"]
    rate_limit: 100  # requests per minute
  
  anthropic:
    api_key: "${ANTHROPIC_API_KEY}"
    models: ["claude-3-5-sonnet-20241022", "claude-3-haiku-20240307"]
    rate_limit: 50
  
  google:
    api_key: "${GOOGLE_API_KEY}"
    search_enabled: true
    models: ["gemini-2.0-flash-exp"]
    rate_limit: 75
  
  ollama:
    base_url: "http://localhost:11434"
    models: ["llama3.1:8b", "mistral", "codellama"]
    pull_missing: true

# CLI provider integrations
cli_providers:
  gemini:
    path: "/usr/local/bin/gemini"
    args: ["--approval-mode", "auto_edit"]
  
  codex:
    path: "/usr/local/bin/codex"
    args: ["-m", "o3-mini", "--headless"]
  
  opencode:
    path: "/usr/local/bin/opencode"
    args: ["--headless", "--background"]

# Workflow definitions
workflows:
  development:
    phases:
      - name: "discovery"
        agents: ["gemini", "ollama"]
        timeout: 600s
      - name: "planning"
        agents: ["opencode", "gemini"]
        timeout: 900s
      - name: "implementation"
        agents: ["codex", "ollama"]
        timeout: 1800s
      - name: "testing"
        agents: ["gemini", "codex"]
        timeout: 900s
      - name: "deployment"
        agents: ["fabric", "opencode"]
        timeout: 600s

# Pattern configuration
patterns:
  paths:
    - "~/.config/fabric-lite/patterns"
    - "./patterns"
    - "./community-patterns"
  auto_update: true
  quality_threshold: 0.7

# Artifact management
artifacts:
  storage_path: ".forge/artifacts"
  compression: true
  retention_days: 90
  max_versions: 10
```

## 5. Integration Components

### 5.1 Browser Automation (`internal/browser/`)
**Agent-Browser Integration**:

```go
type BrowserAutomation struct {
    client   *http.Client
    config   BrowserConfig
    logger   *zap.Logger
}

type BrowserConfig struct {
    AgentBrowserPath string        `yaml:"agent_browser_path"`
    DefaultTimeout time.Duration `yaml:"default_timeout"`
    Headless       bool          `yaml:"headless"`
}

type BrowserTask struct {
    ID       string                 `json:"id"`
    URL      string                 `json:"url"`
    Actions  []BrowserAction       `json:"actions"`
    Context  map[string]interface{} `json:"context"`
}

type BrowserAction struct {
    Type    string      `json:"type"`    // click, fill, scroll, wait, screenshot
    Target  string      `json:"target"`  // selector or xpath
    Value   interface{} `json:"value"`   // text to fill or coordinates
    Timeout time.Duration `json:"timeout"`
}
```

### 5.2 Webhook System (`internal/webhook/`)
**Integration Hooks**:

```go
type WebhookManager struct {
    webhooks map[string]*Webhook
    server   *http.Server
    logger   *zap.Logger
}

type Webhook struct {
    ID       string            `json:"id"`
    URL      string            `json:"url"`
    Events   []WebhookEvent    `json:"events"`
    Secret   string            `json:"secret"`
    Active   bool              `json:"active"`
}

type WebhookEvent struct {
    Type      string            `json:"type"`
    Timestamp time.Time         `json:"timestamp"`
    Data      map[string]string `json:"data"`
    Agent     string            `json:"agent"`
    Phase     string            `json:"phase"`
}
```

## 6. Testing Infrastructure Components

### 6.1 Test Framework (`internal/testing/`)
**Multi-Agent Testing**:

```go
type AgentTestSuite struct {
    agents    []*MockAgent
    workspace string
    logger    *zap.Logger
}

type MockAgent struct {
    ID       string
    Responses map[string]MockResponse
    Delay    time.Duration
    Errors   map[string]error
}

type MockResponse struct {
    Content string
    Delay   time.Duration
    Headers map[string]string
}
```

### 6.2 Performance Testing (`internal/perf/`)
**Load Testing**:

```go
type PerformanceTest struct {
    Scenario    string
    Agents     int
    Duration   time.Duration
    Metrics    *PerformanceMetrics
}

type PerformanceMetrics struct {
    Throughput     float64       `json:"throughput"`     // requests per second
    Latency       time.Duration  `json:"latency"`       // average response time
    ErrorRate     float64       `json:"error_rate"`     // percentage of errors
    ResourceUsage  ResourceUsage  `json:"resource_usage"`
}
```

---

**Component Status**: Detailed breakdown complete  
**Implementation Priority**: Core orchestration → Provider system → Pattern system → CLI enhancements  
**Dependencies**: Go 1.19+, external AI providers, agent-browser  
**Testing Strategy**: Unit tests for each component, integration tests for workflows