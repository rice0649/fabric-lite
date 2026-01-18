---
status: complete
agent: planning_agent
timestamp: 2026-01-18T07:00:00Z
next: review_agent
---

# Fabric-Lite Project Planning Document

## 1. Executive Summary

This document outlines the comprehensive plan for building **fabric-lite**, a streamlined version of the danielmiessler/fabric AI augmentation framework. The goal is to create a personal, lightweight CLI tool that executes AI patterns against text input while maintaining extensibility for future enhancements.

**Key Decisions:**
- **Language**: Go (maintains compatibility with original Fabric, excellent CLI tooling)
- **Initial Scope**: MVP with single provider (OpenAI), pattern execution, basic CLI
- **Target**: Solo developer workflow optimized for rice0649's GitHub

---

## 2. Architecture Overview

### 2.1 Original Fabric Architecture Analysis

The original Fabric project (`/home/oak38/projects/fabric/`) is a mature Go application with:

**Core Components:**
```
fabric/
├── cmd/fabric/main.go          # Entry point - delegates to internal/cli
├── internal/
│   ├── cli/cli.go              # Flag parsing, command dispatch
│   ├── core/
│   │   ├── plugin_registry.go  # Central registry for all plugins
│   │   └── chatter.go          # Chat session management
│   ├── plugins/
│   │   ├── ai/                 # Vendor interface + 15+ implementations
│   │   │   ├── vendor.go       # Core Vendor interface
│   │   │   ├── openai/
│   │   │   ├── ollama/
│   │   │   ├── anthropic/
│   │   │   └── ...
│   │   └── db/fsdb/            # File-system database for patterns/sessions
│   ├── domain/                 # Shared types and utilities
│   └── tools/                  # YouTube, Jina, pattern loaders
├── data/
│   ├── patterns/               # 234 patterns (system.md + user.md each)
│   └── strategies/             # Meta-prompting strategies
└── web/                        # REST API server (Gin-based)
```

**Key Interfaces:**

```go
// Vendor interface - core abstraction for AI providers
type Vendor interface {
    plugins.Plugin
    ListModels() ([]string, error)
    SendStream([]*chat.ChatCompletionMessage, *domain.ChatOptions, chan domain.StreamUpdate) error
    Send(context.Context, []*chat.ChatCompletionMessage, *domain.ChatOptions) (string, error)
    NeedsRawMode(modelName string) bool
}

// Plugin interface - base for all configurable components
type Plugin interface {
    GetName() string
    Configure() error
    IsConfigured() bool
    Setup() error
    GetSetupDescription() string
    SetupFillEnvFileContent(*bytes.Buffer)
}
```

**Data Flow:**
```
User Input → CLI Flags → Pattern Loader → Session Builder → Vendor.Send() → Output
     ↓                       ↓                  ↓
   stdin              ~/.config/fabric    patterns/summarize/
                      /patterns/          system.md + user.md
```

### 2.2 Simplified Architecture for fabric-lite

```
fabric-lite/
├── cmd/fabric-lite/main.go     # Entry point
├── internal/
│   ├── cli/                    # Cobra-based CLI
│   │   ├── root.go
│   │   └── flags.go
│   ├── core/
│   │   ├── pattern.go          # Pattern loading and variable substitution
│   │   ├── executor.go         # Execute patterns against providers
│   │   └── config.go           # Configuration management
│   └── providers/
│       ├── provider.go         # Provider interface
│       ├── openai/             # OpenAI implementation
│       ├── ollama/             # Ollama implementation (Tier 2)
│       └── anthropic/          # Anthropic implementation (Tier 2)
├── patterns/                   # Built-in patterns
├── config/                     # Config templates
└── scripts/                    # Build scripts
```

---

## 3. Feature Matrix (with Priorities)

### Tier 1 - MVP (Must Have)

| Feature | Description | Complexity | Files |
|---------|-------------|------------|-------|
| Pattern Loading | Load system.md + user.md from pattern directories | Low | `internal/core/pattern.go` |
| OpenAI Provider | Basic OpenAI API integration | Medium | `internal/providers/openai/` |
| CLI Interface | Basic flags: --pattern, --list, --help, --model | Low | `internal/cli/` |
| Config Management | YAML config + env vars for API keys | Low | `internal/core/config.go` |
| Stdin Processing | Read input from stdin/pipes | Low | `cmd/fabric-lite/main.go` |
| Output to stdout | Print AI response to terminal | Low | Built-in |

### Tier 2 - Enhanced

| Feature | Description | Complexity | Files |
|---------|-------------|------------|-------|
| Ollama Provider | Local LLM support | Medium | `internal/providers/ollama/` |
| Anthropic Provider | Claude API support | Medium | `internal/providers/anthropic/` |
| Pattern Variables | `{{variable}}` substitution in patterns | Low | `internal/core/pattern.go` |
| Streaming Output | Real-time token streaming | Medium | All providers |
| Session Management | Save/load conversation history | Medium | `internal/core/session.go` |
| Custom Pattern Directory | User-defined patterns in ~/.config | Low | `internal/core/pattern.go` |

### Tier 3 - Advanced

| Feature | Description | Complexity | Files |
|---------|-------------|------------|-------|
| REST API Server | HTTP API for patterns | High | `internal/server/` |
| Web UI | Simple browser interface | High | `web/` |
| Plugin System | Dynamic provider loading | High | `internal/plugins/` |
| Context Support | Load context files | Low | `internal/core/context.go` |
| Strategies | Meta-prompting frameworks | Medium | `internal/core/strategy.go` |

---

## 4. Project Structure

```
fabric-lite/
├── cmd/
│   └── fabric-lite/
│       └── main.go                 # Entry point
├── internal/
│   ├── cli/
│   │   ├── root.go                 # Root command (Cobra)
│   │   ├── run.go                  # Pattern execution command
│   │   ├── list.go                 # List patterns command
│   │   └── setup.go                # Interactive setup
│   ├── core/
│   │   ├── pattern.go              # Pattern struct and loading
│   │   ├── pattern_loader.go       # Find and load patterns
│   │   ├── executor.go             # Execute pattern with provider
│   │   ├── config.go               # Configuration management
│   │   └── types.go                # Shared types
│   └── providers/
│       ├── provider.go             # Provider interface
│       ├── registry.go             # Provider registry
│       └── openai/
│           ├── client.go           # OpenAI client
│           └── models.go           # Model definitions
├── patterns/
│   ├── summarize/
│   │   ├── system.md
│   │   └── user.md
│   ├── analyze_code/
│   │   └── system.md
│   └── ... (10-15 core patterns)
├── config/
│   └── config.example.yaml         # Example configuration
├── scripts/
│   ├── build.sh
│   └── install.sh
├── docs/
│   └── patterns.md                 # Pattern authoring guide
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── LICENSE
```

---

## 5. Technology Stack Decisions

### Primary Language: **Go**

**Justification:**
1. **Compatibility**: Original Fabric is Go; easier to port patterns and understand design
2. **Single Binary**: No runtime dependencies, easy distribution
3. **CLI Excellence**: Excellent libraries (Cobra, Viper)
4. **Concurrency**: Built-in goroutines for streaming
5. **Cross-Platform**: Easy cross-compilation

### Dependencies

| Package | Purpose | Version |
|---------|---------|---------|
| `github.com/spf13/cobra` | CLI framework | v1.8+ |
| `github.com/spf13/viper` | Configuration | v1.18+ |
| `github.com/openai/openai-go` | OpenAI client | latest |
| `gopkg.in/yaml.v3` | YAML parsing | v3 |

### Testing Framework

- **Unit Tests**: Standard `testing` package
- **Integration Tests**: Separate `_test.go` files with build tags

### CI/CD

- **GitHub Actions** for:
  - Build validation
  - Test execution
  - Release automation (goreleaser)

---

## 6. Implementation Phases

### Phase 1: Foundation (Core MVP)
**Goal**: Basic pattern execution with OpenAI

**Tasks:**
1. Set up project structure with go.mod
2. Implement `internal/core/pattern.go` - Pattern struct and loading
3. Implement `internal/core/config.go` - YAML config + env vars
4. Implement `internal/providers/openai/client.go` - Basic OpenAI calls
5. Implement `internal/cli/root.go` - Cobra root command
6. Add 5 core patterns: summarize, analyze_code, explain, improve_writing, extract_wisdom

**Deliverable**: `echo "text" | fabric-lite -p summarize` works

### Phase 2: CLI Enhancement
**Goal**: Full CLI experience

**Tasks:**
1. Add `--list` command to show available patterns
2. Add `--model` flag for model selection
3. Add `--setup` for interactive configuration
4. Implement streaming output (`--stream`)
5. Add `--output` flag for file output
6. Add `--dry-run` for debugging

**Deliverable**: Feature-complete CLI matching basic Fabric functionality

### Phase 3: Additional Providers
**Goal**: Multi-provider support

**Tasks:**
1. Create provider interface abstraction
2. Implement Ollama provider
3. Implement Anthropic provider
4. Add provider selection via `--vendor` flag
5. Implement model listing per provider

**Deliverable**: Works with OpenAI, Ollama, and Anthropic

### Phase 4: Polish & Distribution
**Goal**: Production-ready release

**Tasks:**
1. Comprehensive testing (unit + integration)
2. Documentation (README, pattern guide)
3. goreleaser configuration
4. GitHub Actions for CI/CD
5. Homebrew formula (optional)

**Deliverable**: v1.0.0 release on GitHub

---

## 7. Patterns to Port Initially

Based on analysis of most valuable patterns for a developer workflow:

### Core Analysis Patterns
| Pattern | Purpose | Location in Fabric |
|---------|---------|-------------------|
| `summarize` | Content summarization | `data/patterns/summarize/` |
| `analyze_code` | Code review and analysis | `data/patterns/analyze_code/` |
| `analyze_paper` | Academic paper analysis | `data/patterns/analyze_paper/` |
| `analyze_logs` | Log file analysis | `data/patterns/analyze_logs/` |

### Content Creation Patterns
| Pattern | Purpose | Location in Fabric |
|---------|---------|-------------------|
| `improve_writing` | Text enhancement | `data/patterns/improve_writing/` |
| `create_5_sentence_summary` | Quick summaries | `data/patterns/create_5_sentence_summary/` |
| `explain_code` | Code explanation | `data/patterns/explain_code/` |

### Developer-Focused Patterns
| Pattern | Purpose | Location in Fabric |
|---------|---------|-------------------|
| `coding_master` | General coding assistance | `data/patterns/coding_master/` |
| `create_coding_feature` | Feature implementation | `data/patterns/create_coding_feature/` |
| `extract_wisdom` | Extract insights | `data/patterns/extract_wisdom/` |

---

## 8. Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| API changes in openai-go | Medium | Medium | Pin versions, monitor releases |
| Pattern format incompatibility | Low | Low | Follow Fabric's format exactly |
| Scope creep | High | Medium | Strict MVP focus, defer features |
| Limited testing time | Medium | Medium | Focus on critical paths |
| Config complexity | Low | Low | Start simple, add features later |

---

## 9. Success Criteria

### MVP (Phase 1)
- [ ] Can execute `summarize` pattern on stdin text
- [ ] OpenAI API key configured via env or config
- [ ] Outputs AI response to stdout
- [ ] Handles errors gracefully

### Full Release (Phase 4)
- [ ] 10+ working patterns ported from Fabric
- [ ] OpenAI, Ollama, Anthropic providers functional
- [ ] Streaming output works
- [ ] `fabric-lite --setup` configures API keys
- [ ] `fabric-lite --list` shows all patterns
- [ ] CI/CD pipeline with releases
- [ ] README with usage examples
- [ ] Pattern authoring documentation

---

## 10. File Mapping: Fabric → fabric-lite

| Fabric File | fabric-lite Equivalent | Notes |
|-------------|------------------------|-------|
| `cmd/fabric/main.go` | `cmd/fabric-lite/main.go` | Simplified |
| `internal/cli/cli.go` | `internal/cli/root.go` | Cobra-based |
| `internal/core/chatter.go` | `internal/core/executor.go` | Simplified |
| `internal/plugins/ai/vendor.go` | `internal/providers/provider.go` | Simplified interface |
| `internal/plugins/ai/openai/` | `internal/providers/openai/` | Direct port |
| `internal/plugins/db/fsdb/patterns.go` | `internal/core/pattern_loader.go` | Simplified |
| `data/patterns/` | `patterns/` | Selective port |

---

## Next Steps

1. **Review Agent**: Validate this plan against project constraints
2. **Code Analysis Agent**: Deep dive into specific Fabric files for porting
3. **Implementation Agent**: Begin Phase 1 implementation

---

*Generated by Planning Agent for fabric-lite project*
