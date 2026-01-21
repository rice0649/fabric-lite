#!/bin/bash
# Agent-specific configuration launcher
# Each agent runs with its optimal settings in headless background mode

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
CONFIG_FILE="$PROJECT_DIR/config.yaml"

# Agent configurations based on their strengths
setup_gemini() {
    echo "Setting up Gemini for Research & Discovery..."
    
    cat << 'EOF' > "$PROJECT_DIR/.forge/agents/gemini_config.yaml"
agent: gemini
model: gemini-2.0-flash-exp
provider: google
mode: headless
background: true

# Gemini strengths configuration
specialization:
  - research_and_discovery
  - large_codebase_analysis
  - validation_and_review
  - test_coverage_analysis
  - best_practices_review

# Optimize for research tasks
parameters:
  temperature: 0.7
  max_tokens: 8192
  search_enabled: true
  context_window: "1M"

# Headless mode settings
headless:
  auto_approve: true
  use_defaults: true
  no_interactive: true
  batch_processing: true

# Resource allocation
resources:
  memory: "2GB"
  timeout: 600
  retry_attempts: 3
EOF
}

setup_codex() {
    echo "Setting up Codex for Code Generation..."
    
    cat << 'EOF' > "$PROJECT_DIR/.forge/agents/codex_config.yaml"
agent: codex
model: o3-mini
provider: openai
mode: headless
background: true

# Codex strengths configuration
specialization:
  - code_generation
  - implementation
  - refactoring
  - advanced_reasoning
  - code_review

# Optimize for coding tasks
parameters:
  temperature: 0.1
  max_tokens: 16384
  reasoning_mode: "high"
  code_focus: true

# Headless mode settings
headless:
  auto_approve: true
  use_defaults: true
  no_interactive: true
  batch_processing: true

# Resource allocation
resources:
  memory: "4GB"
  timeout: 1800
  retry_attempts: 3
EOF
}

setup_opencode() {
    echo "Setting up OpenCode for Planning & Architecture..."
    
    cat << 'EOF' > "$PROJECT_DIR/.forge/agents/opencode_config.yaml"
agent: opencode
model: claude-sonnet-4-20250514
provider: anthropic
mode: headless
background: true

# OpenCode strengths configuration
specialization:
  - planning_and_architecture
  - design_exploration
  - api_modeling
  - large_scale_analysis
  - structural_integrity

# Optimize for planning tasks
parameters:
  temperature: 0.5
  max_tokens: 12288
  planning_mode: true
  architecture_focus: true

# Headless mode settings
headless:
  auto_approve: true
  use_defaults: true
  no_interactive: true
  batch_processing: true

# Resource allocation
resources:
  memory: "3GB"
  timeout: 1200
  retry_attempts: 3
EOF
}

setup_fabric() {
    echo "Setting up Fabric for Pattern-based Tasks..."
    
    cat << 'EOF' > "$PROJECT_DIR/.forge/agents/fabric_config.yaml"
agent: fabric
model: gpt-4o-mini
provider: openai
mode: headless
background: true

# Fabric strengths configuration
specialization:
  - documentation_generation
  - summarization
  - pattern_based_tasks
  - changelog_creation
  - release_notes

# Optimize for pattern-based tasks
parameters:
  temperature: 0.3
  max_tokens: 4096
  patterns_dir: "~/.config/fabric-lite/patterns"
  reproducible_outputs: true

# Headless mode settings
headless:
  auto_approve: true
  use_defaults: true
  no_interactive: true
  batch_processing: true

# Resource allocation
resources:
  memory: "1GB"
  timeout: 300
  retry_attempts: 3
EOF
}

setup_ollama() {
    echo "Setting up Ollama for Quick Tasks..."
    
    cat << 'EOF' > "$PROJECT_DIR/.forge/agents/ollama_config.yaml"
agent: ollama
model: llama3.1:8b
provider: local
mode: headless
background: true

# Ollama strengths configuration
specialization:
  - quick_tasks
  - boilerplate_generation
  - formatting
  - simple_automation
  - preliminary_analysis

# Optimize for quick tasks
parameters:
  temperature: 0.2
  max_tokens: 2048
  fast_mode: true
  local_processing: true

# Headless mode settings
headless:
  auto_approve: true
  use_defaults: true
  no_interactive: true
  batch_processing: true

# Resource allocation
resources:
  memory: "1GB"
  timeout: 120
  retry_attempts: 1
EOF
}

# Create agent-specific task runners
create_agent_runners() {
    echo "Creating agent-specific task runners..."
    
    # Gemini runner
    cat << 'EOF' > "$PROJECT_DIR/.forge/agents/run_gemini.sh"
#!/bin/bash
# Gemini task runner for headless mode

TASK="$1"
OUTPUT_DIR="$2"
LOG_FILE="$3"

echo "Gemini processing: $TASK" >> "$LOG_FILE"
gemini -p headless --model gemini-2.0-flash-exp --temperature 0.7 "$TASK" >> "$LOG_FILE" 2>&1
EOF

    # Codex runner
    cat << 'EOF' > "$PROJECT_DIR/.forge/agents/run_codex.sh"
#!/bin/bash
# Codex task runner for headless mode

TASK="$1"
OUTPUT_DIR="$2"
LOG_FILE="$3"

echo "Codex processing: $TASK" >> "$LOG_FILE"
codex -p headless -m o3-mini --temperature 0.1 "$TASK" >> "$LOG_FILE" 2>&1
EOF

    # OpenCode runner
    cat << 'EOF' > "$PROJECT_DIR/.forge/agents/run_opencode.sh"
#!/bin/bash
# OpenCode task runner for headless mode

TASK="$1"
OUTPUT_DIR="$2"
LOG_FILE="$3"

echo "OpenCode processing: $TASK" >> "$LOG_FILE"
opencode --headless --background --task "$TASK" >> "$LOG_FILE" 2>&1
EOF

    # Fabric runner
    cat << 'EOF' > "$PROJECT_DIR/.forge/agents/run_fabric.sh"
#!/bin/bash
# Fabric task runner for headless mode

PATTERN="$1"
INPUT="$2"
OUTPUT="$3"
LOG_FILE="$4"

echo "Fabric processing pattern: $PATTERN" >> "$LOG_FILE"
fabric-lite --pattern "$PATTERN" --input "$INPUT" --output "$OUTPUT" --headless >> "$LOG_FILE" 2>&1
EOF

    # Ollama runner
    cat << 'EOF' > "$PROJECT_DIR/.forge/agents/run_ollama.sh"
#!/bin/bash
# Ollama task runner for headless mode

TASK="$1"
OUTPUT_DIR="$2"
LOG_FILE="$3"

echo "Ollama processing: $TASK" >> "$LOG_FILE"
ollama run llama3.1:8b "$TASK" >> "$LOG_FILE" 2>&1
EOF

    # Make all runners executable
    chmod +x "$PROJECT_DIR/.forge/agents/run_*.sh"
}

# Create agent orchestration matrix
create_orchestration_matrix() {
    echo "Creating agent orchestration matrix..."
    
    cat << 'EOF' > "$PROJECT_DIR/.forge/agents/orchestration_matrix.yaml"
# Agent orchestration matrix for optimal task assignment
# Maps phases and tasks to specialized agents

phases:
  discovery:
    primary: gemini
    secondary: ollama
    tasks:
      research:
        agent: gemini
        reason: "Large context window + search integration"
      quick_analysis:
        agent: ollama
        reason: "Fast local processing for preliminary work"
      market_research:
        agent: gemini
        reason: "Access to current web information"
        
  planning:
    primary: opencode
    secondary: gemini
    tasks:
      architecture_design:
        agent: opencode
        reason: "Strong architectural reasoning"
      task_breakdown:
        agent: opencode
        reason: "Excellent planning capabilities"
      validation:
        agent: gemini
        reason: "Best practices review"
        
  implementation:
    primary: codex
    secondary: ollama
    tasks:
      code_generation:
        agent: codex
        reason: "Advanced reasoning for complex code"
      boilerplate:
        agent: ollama
        reason: "Quick template generation"
      refactoring:
        agent: codex
        reason: "Deep code understanding"
        
  testing:
    primary: gemini
    secondary: codex
    tasks:
      test_strategy:
        agent: gemini
        reason: "Comprehensive coverage analysis"
      test_implementation:
        agent: codex
        reason: "Precise test code generation"
        
  deployment:
    primary: fabric
    secondary: opencode
    tasks:
      documentation:
        agent: fabric
        reason: "Pattern-based documentation"
      deployment_pipeline:
        agent: opencode
        reason: "Infrastructure as code"
      release_notes:
        agent: fabric
        reason: "Standardized output format"

# Agent collaboration patterns
collaboration:
  - phase: discovery_to_planning
    from: gemini
    to: opencode
    handoff: "discovery_report.md"
    
  - phase: planning_to_implementation
    from: opencode
    to: codex
    handoff: "implementation_spec.md"
    
  - phase: implementation_to_testing
    from: codex
    to: gemini
    handoff: "code_summary.md"
    
  - phase: testing_to_deployment
    from: gemini
    to: fabric
    handoff: "test_results.md"
EOF
}

# Main setup function
main() {
    echo "=========================================="
    echo "  Setting up Agent Configurations"
    echo "=========================================="
    echo ""
    
    # Create directories
    mkdir -p "$PROJECT_DIR/.forge/agents"
    mkdir -p "$PROJECT_DIR/.forge/logs"
    mkdir -p "$PROJECT_DIR/.forge/artifacts"
    mkdir -p "$PROJECT_DIR/.forge/processes"
    
    # Setup each agent
    setup_gemini
    setup_codex
    setup_opencode
    setup_fabric
    setup_ollama
    
    # Create runners and orchestration
    create_agent_runners
    create_orchestration_matrix
    
    echo ""
    echo "✅ All agent configurations created successfully!"
    echo ""
    echo "Agent configurations located at: $PROJECT_DIR/.forge/agents/"
    echo "Orchestration matrix: $PROJECT_DIR/.forge/agents/orchestration_matrix.yaml"
    echo ""
    echo "To run the workflow:"
    echo "  $PROJECT_DIR/headless_workflow.sh all"
    echo ""
    echo "To monitor background processes:"
    echo "  $PROJECT_DIR/headless_workflow.sh monitor"
}

main "$@"