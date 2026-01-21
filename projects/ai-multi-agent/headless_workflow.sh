#!/bin/bash
# Headless Multi-Agent Launcher for Forge + Fabric-Lite
# Runs all agents in background with no interactive prompts

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONFIG_FILE="$PROJECT_DIR/config.yaml"
LOG_DIR="$PROJECT_DIR/.forge/logs"
ARTIFACTS_DIR="$PROJECT_DIR/.forge/artifacts"
PROCESSES_DIR="$PROJECT_DIR/.forge/processes"

# Create directories
mkdir -p "$LOG_DIR" "$ARTIFACTS_DIR" "$PROCESSES_DIR"

echo "=========================================="
echo "  Headless Multi-Agent Workflow Launcher"
echo "=========================================="
echo "Project: $(basename "$PROJECT_DIR")"
echo "Mode: Headless background processing"
echo ""

# Parse configuration for tools and phases
parse_config() {
    if command -v yq >/dev/null 2>&1; then
        AGENTS=$(yq '.tools | keys | .[]' "$CONFIG_FILE" 2>/dev/null || echo "")
        PHASES=$(yq '.phases | keys | .[]' "$CONFIG_FILE" 2>/dev/null || echo "")
    else
        # Fallback to common agents if yq not available
        AGENTS="gemini codex opencode fabric ollama"
        PHASES="discovery planning implementation testing deployment"
    fi
}

# Agent execution functions
run_gemini() {
    local task="$1"
    local output_file="$2"
    local log_file="$LOG_DIR/gemini_$(date +%s).log"
    
    echo "Starting Gemini (headless) - $task"
    nohup gemini -p headless "$task" > "$log_file" 2>&1 &
    echo $! > "$PROCESSES_DIR/gemini.pid"
    echo "Gemini PID: $(cat "$PROCESSES_DIR/gemini.pid")"
    echo "Log: $log_file"
}

run_codex() {
    local task="$1"
    local output_file="$2"
    local log_file="$LOG_DIR/codex_$(date +%s).log"
    
    echo "Starting Codex (headless) - $task"
    nohup codex -p headless -m o3-mini "$task" > "$log_file" 2>&1 &
    echo $! > "$PROCESSES_DIR/codex.pid"
    echo "Codex PID: $(cat "$PROCESSES_DIR/codex.pid")"
    echo "Log: $log_file"
}

run_opencode() {
    local task="$1"
    local output_file="$2"
    local log_file="$LOG_DIR/opencode_$(date +%s).log"
    
    echo "Starting OpenCode (headless) - $task"
    nohup opencode --headless --background --task "$task" > "$log_file" 2>&1 &
    echo $! > "$PROCESSES_DIR/opencode.pid"
    echo "OpenCode PID: $(cat "$PROCESSES_DIR/opencode.pid")"
    echo "Log: $log_file"
}

run_fabric() {
    local pattern="$1"
    local input="$2"
    local output_file="$3"
    local log_file="$LOG_DIR/fabric_$(date +%s).log"
    
    echo "Starting Fabric (headless) - $pattern"
    nohup fabric-lite --pattern "$pattern" --input "$input" --output "$output_file" --headless > "$log_file" 2>&1 &
    echo $! > "$PROCESSES_DIR/fabric.pid"
    echo "Fabric PID: $(cat "$PROCESSES_DIR/fabric.pid")"
    echo "Log: $log_file"
}

run_ollama() {
    local task="$1"
    local output_file="$2"
    local log_file="$LOG_DIR/ollama_$(date +%s).log"
    
    echo "Starting Ollama (headless) - $task"
    nohup ollama run llama3.1:8b "$task" > "$log_file" 2>&1 &
    echo $! > "$PROCESSES_DIR/ollama.pid"
    echo "Ollama PID: $(cat "$PROCESSES_DIR/ollama.pid")"
    echo "Log: $log_file"
}

# Phase execution workflows
execute_discovery_phase() {
    echo "=========================================="
    echo "  DISCOVERY PHASE (Headless)"
    echo "=========================================="
    
    local timestamp=$(date +%s)
    local discovery_output="$ARTIFACTS_DIR/discovery_$timestamp.md"
    
    # Start Gemini for research and discovery
    run_gemini "Analyze the project requirements and perform comprehensive research. Focus on: 1) Technical requirements, 2) Best practices, 3) Potential challenges, 4) Similar solutions in the wild. Output a detailed discovery report." "$discovery_output"
    
    # Start Ollama for quick preliminary analysis
    run_ollama "Generate a quick project overview based on available documentation. Identify key components and dependencies." "$ARTIFACTS_DIR/quick_overview_$timestamp.md"
    
    echo "Discovery agents started in background"
    echo "Monitoring progress..."
}

execute_planning_phase() {
    echo "=========================================="
    echo "  PLANNING PHASE (Headless)"
    echo "=========================================="
    
    local timestamp=$(date +%s)
    local planning_output="$ARTIFACTS_DIR/planning_$timestamp.md"
    
    # Start OpenCode for master planning
    run_opencode "Create a comprehensive project plan based on discovery results. Include: 1) Architecture overview, 2) Development phases, 3) Task breakdown, 4) Resource requirements, 5) Timeline estimates." "$planning_output"
    
    # Start Gemini for planning validation
    run_gemini "Review and validate the project plan. Check for: 1) Feasibility, 2) Best practices alignment, 3) Risk assessment, 4) Optimization opportunities." "$ARTIFACTS_DIR/planning_review_$timestamp.md"
    
    echo "Planning agents started in background"
}

execute_implementation_phase() {
    echo "=========================================="
    echo "  IMPLEMENTATION PHASE (Headless)"
    echo "=========================================="
    
    local timestamp=$(date +%s)
    
    # Start Codex for code generation
    run_codex "Implement the planned features following the architecture. Generate production-ready code with proper error handling, logging, and documentation." "$ARTIFACTS_DIR/implementation_$timestamp.md"
    
    # Start Ollama for boilerplate and setup
    run_ollama "Generate necessary boilerplate code, configuration files, and setup scripts." "$ARTIFACTS_DIR/boilerplate_$timestamp.md"
    
    echo "Implementation agents started in background"
}

execute_testing_phase() {
    echo "=========================================="
    echo "  TESTING PHASE (Headless)"
    echo "=========================================="
    
    local timestamp=$(date +%s)
    
    # Start Gemini for comprehensive testing strategy
    run_gemini "Generate comprehensive test suite including unit tests, integration tests, and E2E tests. Focus on coverage and edge cases." "$ARTIFACTS_DIR/testing_$timestamp.md"
    
    # Start Codex for test implementation
    run_codex "Implement the test suite with proper fixtures, mocks, and test data. Ensure all tests are runnable and pass." "$ARTIFACTS_DIR/test_implementation_$timestamp.md"
    
    echo "Testing agents started in background"
}

execute_deployment_phase() {
    echo "=========================================="
    echo "  DEPLOYMENT PHASE (Headless)"
    echo "=========================================="
    
    local timestamp=$(date +%s)
    
    # Start Fabric for documentation and deployment configs
    run_fabric "generate_deployment_docs" "Implementation results" "$ARTIFACTS_DIR/deployment_docs_$timestamp.md"
    
    # Start OpenCode for deployment pipeline
    run_opencode "Create deployment pipeline configurations, CI/CD setup, and infrastructure as code." "$ARTIFACTS_DIR/deployment_pipeline_$timestamp.md"
    
    echo "Deployment agents started in background"
}

# Process monitoring
monitor_processes() {
    echo ""
    echo "=========================================="
    echo "  MONITORING BACKGROUND PROCESSES"
    echo "=========================================="
    
    while true; do
        local active=0
        echo "$(date): Checking active processes..."
        
        for pid_file in "$PROCESSES_DIR"/*.pid; do
            if [[ -f "$pid_file" ]]; then
                local pid=$(cat "$pid_file")
                if kill -0 "$pid" 2>/dev/null; then
                    echo "  $(basename "$pid_file" .pid): Running (PID: $pid)"
                    ((active++))
                else
                    echo "  $(basename "$pid_file" .pid): Completed"
                    rm "$pid_file"
                fi
            fi
        done
        
        if [[ $active -eq 0 ]]; then
            echo "All processes completed!"
            break
        else
            echo "$active processes still running..."
            sleep 30
        fi
    done
}

# Main execution
main() {
    parse_config
    
    case "${1:-all}" in
        "discovery")
            execute_discovery_phase
            monitor_processes
            ;;
        "planning")
            execute_planning_phase
            monitor_processes
            ;;
        "implementation")
            execute_implementation_phase
            monitor_processes
            ;;
        "testing")
            execute_testing_phase
            monitor_processes
            ;;
        "deployment")
            execute_deployment_phase
            monitor_processes
            ;;
        "all")
            execute_discovery_phase
            sleep 5
            execute_planning_phase
            sleep 5
            execute_implementation_phase
            sleep 5
            execute_testing_phase
            sleep 5
            execute_deployment_phase
            monitor_processes
            ;;
        "monitor")
            monitor_processes
            ;;
        *)
            echo "Usage: $0 [discovery|planning|implementation|testing|deployment|all|monitor]"
            echo ""
            echo "Examples:"
            echo "  $0 all                    # Run all phases sequentially"
            echo "  $0 discovery              # Run only discovery phase"
            echo "  $0 monitor                # Monitor running processes"
            exit 1
            ;;
    esac
}

main "$@"