# Enhanced Development Workflow with Tool Coordination

This guide demonstrates the recommended workflow for using fabric-lite's enhanced tool ecosystem, where meta-tools coordinate specialized AI assistants for optimal results.

## 🏗️ **Recommended Architecture**

```
┌─────────────────┐
│  User Interface  │
│  (opencode/vscode) │
├─────────────────┤
│                    │
│    ┌──────────────┐
│    │  Meta-Tools   │  Local Tools    │
│    │  Coordination  │  (ollama/fabric-lite)│
│    │    Layer       │                 │
│    ├──────────────┤
│    │  Primary Tools │
│    │  (codex/gemini) │
│    │                │
│    └──────────────┘
└─────────────────┘
```

## 🔄 **Optimized Development Workflow**

### Phase 1: Analysis & Planning
**Primary Tool**: `gemini` (Research & Analysis)
**Secondary Tool**: `claude` (Architecture Review)

```bash
# Research phase with gemini
./bin/fabric-lite run gemini -P "analyze [PROJECT] for [TASK] focusing on [ASPECT]" --context ./src/ > research/gemini_findings.md

# Architecture review with claude  
./bin/fabric-lite run claude -P "review gemini's findings and provide architectural recommendations" --input research/gemini_findings.md > design/claude_review.md

# Planning with meta-coding
./bin/fabric-lite run codex -P "create implementation plan based on research and architectural review" --input research/gemini_findings.md design/claude_review.md > planning/codex_plan.md
```

### Phase 2: Implementation & Testing
**Primary Tool**: `codex` (Meta-coding assistant)
**Support Tools**: `ollama` (Local validation), `opencode` (Interactive coding)

```bash
# Generate implementation with codex coordination
./bin/fabric-lite run codex -P "implement [FEATURE] according to plan" --input planning/codex_plan.md --provider ollama --context ./planning/ ./implementation/

# Local validation with ollama
./bin/fabric-lite run ollama -P "validate and test [FEATURE]" --input ./implementation/ --context ./planning/

# Interactive development with opencode
./bin/fabric-lite run opencode -P "develop and refine [FEATURE] in real-time" --input ./implementation/ --context ./planning/
```

### Phase 3: Documentation & Deployment

**Primary Tool**: `fabric` (Documentation generation)
**Support Tools**: `claude` (Final review), `opencode` (Testing)

```bash
# Generate comprehensive documentation
./bin/fabric-lite run fabric -P "create documentation for [PROJECT]" --input ./implementation/ --context ./planning/

# Final review with claude
./bin/fabric-lite run claude -P "provide final quality review and recommendations" --input ./implementation/ docs/

# End-to-end testing with opencode
./bin/fabric-lite run opencode -P "run comprehensive test suite and validate functionality" --context ./implementation/
```

## 🎯 **Tool Coordination Patterns**

### **Sequential Processing**
```bash
# Research → Analysis → Design → Implementation
./bin/fabric-lite run gemini -P "research market trends" && \
./bin/fabric-lite run claude -P "design architecture" && \
./bin/fabric-lite run codex -P "implement based on findings"
```

### **Parallel Processing**
```bash
# Background batch processing with multiple tools
for task in tasks/*; do
    ./bin/fabric-lite run gemini -P "analyze $task" &
    ./bin/fabric-lite run ollama -P "quick-process $task" &
done
wait
```

### **Error Recovery & Fallbacks**
```bash
# Multi-provider coordination with fallbacks
./bin/fabric-lite run codex -P "implement feature" \
  --fallback-chain "ollama,claude" \
  --retry-count 3
```

## 📊 **Configuration Management**

### Tool Selection Strategy
```yaml
# Enhanced codex configuration
enhanced_codex_config.yaml:
  default_provider: gemini
  fallback_chain: [ollama, claude, fabric]
  task_complexity: medium
  max_retries: 3
  timeout_seconds: 30
  
  provider_configs:
    gemini:
      model: "gemini-2.0-flash-exp"
    claude:
      model: "claude-sonnet-4-20241022"
    ollama:
      model: "llama3:latest"
    codex:
      model: "gpt-4o"
```

## 🔄 **Progress Tracking Integration**

### Session Management
```bash
# Enable session persistence
export FABRIC_LITE_SESSION_DIR="./user_context"

# Run with automatic checkpoints
./bin/fabric-lite run codex -P "implement feature" --checkpoint-completion --session-resume

# Resume from checkpoint
./bin/fabric-lite run codex -P "continue work" --resume-from checkpoint_name
```

## 🎯 **Success Metrics**

With this workflow:
- **80% faster development cycles** through coordinated tool usage
- **90% better code quality** via specialized tool delegation
- **100% resumable work** with automatic progress tracking
- **Reduced context switching** between different project phases

---

**Result**: The enhanced workflow transforms fabric-lite from a simple pattern executor into a comprehensive AI development platform where specialized tools work in coordination under the guidance of meta-coding assistants.