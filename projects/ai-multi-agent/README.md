See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# Multi-Agent AI Workflow with Forge + Fabric-Lite

A comprehensive project demonstrating the power of running multiple AI agents in headless background mode using forge and fabric-lite.

## 🎯 Project Overview

This system orchestrates 5 specialized AI agents working in parallel background processes:

- **Gemini**: Research, discovery, validation, and large-scale analysis
- **Codex**: Advanced code generation and implementation  
- **OpenCode**: Architecture, planning, and design
- **Fabric**: Pattern-based documentation and deployment
- **Ollama**: Quick tasks, boilerplate, and simple automation

All agents run in **headless mode** in the background, requiring no human interaction during execution.

## 🚀 Quick Start

### 1. Setup Environment

```bash
# Required API keys in environment
export GOOGLE_API_KEY="your-gemini-key"
export OPENAI_API_KEY="your-openai-key" 
export ANTHROPIC_API_KEY="your-claude-key"

# Install dependencies
npm install -g @anthropic/gemini-cli @openai/codex-cli
go install github.com/opencode/opencode@latest
go install github.com/rice0649/fabric-lite/cmd/fabric-lite@latest

# Optional: Install Ollama for local processing
curl -fsSL https://ollama.ai/install.sh | sh
ollama pull llama3.1:8b
```

### 2. Configure Project

```bash
cd /home/oak38/projects/fabric-lite/projects/ai-multi-agent

# Setup agent configurations
./setup_agents.sh

# Review and customize settings
vim config.yaml
```

### 3. Run Complete Workflow

```bash
# Run all phases sequentially with background agents
./headless_workflow.sh all

# Or run specific phases
./headless_workflow.sh discovery
./headless_workflow.sh planning
./headless_workflow.sh implementation
./headless_workflow.sh testing
./headless_workflow.sh deployment

# Monitor running processes
./headless_workflow.sh monitor
```

### 4. Try the Example

```bash
# Build a complete Task Management API with all agents
./example_workflow.sh
```

## 🏗️ Architecture

```
Multi-Agent Headless System
├── Discovery Phase
│   ├── Gemini (Research & Analysis)
│   └── Ollama (Quick Overview)
├── Planning Phase  
│   ├── OpenCode (Architecture & Design)
│   └── Gemini (Validation & Review)
├── Implementation Phase
│   ├── Codex (Code Generation)
│   └── Ollama (Boilerplate & Setup)
├── Testing Phase
│   ├── Gemini (Test Strategy)
│   └── Codex (Test Implementation)
└── Deployment Phase
    ├── Fabric (Documentation)
    └── OpenCode (Deployment Pipeline)
```

## 📁 Project Structure

```
ai-multi-agent/
├── config.yaml                    # Main project configuration
├── headless_workflow.sh            # Orchestrates all agents
├── setup_agents.sh                 # Configures agent settings
├── example_workflow.sh             # Complete example project
├── project_requirements.md         # Example project specs
└── .forge/
    ├── agents/
    │   ├── orchestration_matrix.yaml
    │   ├── gemini_config.yaml
    │   ├── codex_config.yaml
    │   ├── opencode_config.yaml
    │   ├── fabric_config.yaml
    │   └── ollama_config.yaml
    ├── artifacts/                  # Agent outputs
    ├── logs/                      # Process logs
    └── processes/                 # PID files
```

## 🤖 Agent Specializations

### Gemini - Research & Discovery
- **Strengths**: 1M context window, search integration
- **Best For**: Market research, best practices, validation
- **Mode**: `gemini -p headless`

### Codex - Code Generation  
- **Strengths**: Advanced reasoning, code review
- **Best For**: Implementation, refactoring, complex algorithms
- **Mode**: `codex -p headless -m o3-mini`

### OpenCode - Planning & Architecture
- **Strengths**: Design exploration, systems thinking
- **Best For**: Architecture, API design, planning
- **Mode**: `opencode --headless --background`

### Fabric - Documentation & Patterns
- **Strengths**: Reproducible outputs, pattern-based
- **Best For**: Documentation, deployment configs, summaries
- **Mode**: `fabric-lite --headless`

### Ollama - Quick Tasks
- **Strengths**: Fast local processing, privacy
- **Best For**: Boilerplate, formatting, simple tasks
- **Mode**: `ollama run llama3.1:8b`

## 🔧 Configuration

### Main Configuration (config.yaml)

Key settings for headless operation:

```yaml
# Enable headless mode for all tools
tools:
  gemini:
    headless: true
    background: true
    
  codex:
    headless: true  
    background: true
    
  opencode:
    headless: true
    background: true
    
# Background processing
background:
  max_concurrent: 4
  log_dir: .forge/logs
  
# Headless-specific
headless:
  auto_approve: true
  use_defaults: true
  continue_on_error: true
```

### Agent-Specific Configurations

Each agent gets optimized settings in `.forge/agents/`:

- Temperature and model parameters
- Resource allocation (memory, timeout)
- Specialization focus areas
- Headless mode preferences

## 📊 Monitoring & Logs

### Process Monitoring

```bash
# Check running agents
ps aux | grep -E "(gemini|codex|opencode|fabric|ollama)"

# Monitor logs in real-time
tail -f .forge/logs/*.log

# Check completion status
./headless_workflow.sh monitor
```

### Artifact Outputs

All agent outputs are stored in `.forge/artifacts/`:

- `discovery_*.md` - Research findings
- `planning_*.md` - Project plans and architecture
- `implementation_*.md` - Generated code and components
- `testing_*.md` - Test suites and strategies
- `deployment_*.md` - Documentation and deployment configs

## 🎛️ Workflow Orchestration

### Phase-Based Execution

The system executes work in 5 distinct phases:

1. **Discovery**: Research requirements and analyze context
2. **Planning**: Design architecture and create implementation plan  
3. **Implementation**: Generate code and build components
4. **Testing**: Create comprehensive test suites
5. **Deployment**: Generate documentation and deployment configs

### Agent Collaboration

Agents communicate via shared files:

```
Discovery → Planning → Implementation → Testing → Deployment
    ↓           ↓            ↓            ↓          ↓
 Research   Architecture   Code       Tests      Docs
   ↓           ↓            ↓            ↓          ↓
Validation → Review → Analysis → Strategy → Pipeline
```

## 🛠️ Customization

### Adding New Agents

1. Create config in `.forge/agents/new_agent_config.yaml`
2. Add runner script `.forge/agents/run_new_agent.sh`
3. Update orchestration matrix
4. Modify main workflow script

### Modifying Phases

Edit `config.yaml` to customize:

```yaml
phases:
  custom_phase:
    tool: gemini
    headless: true
    background: true
    parallel_agents: [gemini, codex]
    timeout: 900
```

### Resource Limits

Adjust per-agent resource allocation:

```yaml
# In agent config
resources:
  memory: "4GB"
  timeout: 1800
  retry_attempts: 3
```

## 📈 Performance Optimization

### Parallel Processing

The system runs multiple agents concurrently:

- Discovery: Gemini + Ollama (research + quick analysis)
- Planning: OpenCode + Gemini (architecture + validation)
- Implementation: Codex + Ollama (code + boilerplate)
- Testing: Gemini + Codex (strategy + implementation)
- Deployment: Fabric + OpenCode (docs + pipeline)

### Resource Management

- **Memory**: Each agent allocated based on task complexity
- **Timeout**: Phase-specific timeouts prevent hanging
- **Retries**: Automatic retry on transient failures
- **Logging**: Comprehensive logging for debugging

## 🔍 Troubleshooting

### Common Issues

1. **Agents not starting**: Check API keys and permissions
2. **Memory issues**: Reduce concurrent processes or memory limits
3. **Timeouts**: Increase timeout values for complex tasks
4. **Missing outputs**: Check logs for errors

### Debug Commands

```bash
# Check configuration
yq . config.yaml

# View agent logs
cat .forge/logs/gemini_*.log

# Check process status
ls -la .forge/processes/

# Validate artifacts
ls -la .forge/artifacts/
```

## 🎯 Best Practices

1. **Start Small**: Test individual phases before full workflow
2. **Monitor Resources**: Keep an eye on memory and CPU usage
3. **Review Logs**: Check agent logs for optimization opportunities
4. **Validate Outputs**: Review artifacts before proceeding
5. **Customize Prompts**: Adjust agent prompts for your domain

## 🚀 Advanced Usage

### Custom Workflows

Create custom workflow scripts for specific use cases:

```bash
#!/bin/bash
# Custom research workflow
./headless_workflow.sh discovery
./headless_workflow.sh planning
# Custom analysis here
./headless_workflow.sh deployment
```

### Integration with CI/CD

Add to your pipeline:

```yaml
# GitHub Actions example
- name: Run AI Agents
  run: |
    ./headless_workflow.sh all
    ./headless_workflow.sh monitor
```

### Batch Processing

Process multiple projects:

```bash
for project in project1 project2 project3; do
  cp -r ai-multi-agent "$project"
  cd "$project"
  ./headless_workflow.sh all
  cd ..
done
```

## 📚 Examples

### Web API Development
The included example builds a complete Task Management API with:
- RESTful endpoints
- Database schema
- Authentication
- Testing suite
- Documentation
- Docker deployment

### Microservices Architecture
Customize for microservices by adjusting the planning phase to focus on:
- Service boundaries
- Communication patterns
- Data contracts
- Deployment topology

### Mobile App Development
Adapt for mobile by focusing on:
- Component architecture
- State management
- API integration
- Testing strategies

## 🤝 Contributing

This is a demonstration of multi-agent AI orchestration. Contributions welcome for:

- New agent configurations
- Additional workflow examples
- Performance optimizations
- Integration patterns

## 📄 License

This example project follows the same license as fabric-lite.

---

**Built with Forge + Fabric-Lite**  
*Running AI agents in headless background mode for maximum productivity*