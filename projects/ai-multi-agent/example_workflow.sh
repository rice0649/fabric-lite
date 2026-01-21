#!/bin/bash
# Example: Building a Web API with Multi-Agent Headless Workflow
# Demonstrates the full power of forge + fabric-lite with background agents

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKFLOW_SCRIPT="$PROJECT_DIR/headless_workflow.sh"
SETUP_SCRIPT="$PROJECT_DIR/setup_agents.sh"

echo "================================================"
echo "  Multi-Agent Headless Workflow Example"
echo "================================================"
echo "Project: RESTful API for Task Management"
echo "Architecture: Node.js + Express + PostgreSQL"
echo "Mode: All agents running in headless background"
echo ""

# Ensure setup is complete
if [[ ! -f "$PROJECT_DIR/.forge/agents/orchestration_matrix.yaml" ]]; then
    echo "🔧 Setting up agent configurations..."
    "$SETUP_SCRIPT"
fi

# Create project requirements
create_project_requirements() {
    echo "📝 Creating project requirements..."
    
    cat << 'EOF' > "$PROJECT_DIR/project_requirements.md"
# Task Management API Project Requirements

## Overview
Build a RESTful API for task management with the following features:
- User authentication and authorization
- CRUD operations for tasks
- Task categorization and filtering
- Due date management
- Task status tracking
- API documentation

## Technical Requirements
- **Backend**: Node.js with Express.js
- **Database**: PostgreSQL with Prisma ORM
- **Authentication**: JWT tokens
- **Validation**: Joi or Zod
- **Documentation**: OpenAPI/Swagger
- **Testing**: Jest with Supertest
- **Deployment**: Docker with docker-compose

## API Endpoints
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/tasks` - List tasks (with filters)
- `POST /api/tasks` - Create task
- `GET /api/tasks/:id` - Get task by ID
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task
- `GET /api/docs` - API documentation

## Data Models
### User
- id (UUID, primary)
- email (unique)
- password (hashed)
- name
- created_at
- updated_at

### Task
- id (UUID, primary)
- title
- description
- status (pending, in_progress, completed)
- priority (low, medium, high)
- due_date
- category (string)
- user_id (foreign key)
- created_at
- updated_at

## Non-Functional Requirements
- Response time < 200ms for simple queries
- Support for at least 1000 concurrent users
- 99.9% uptime
- Input validation and sanitization
- Rate limiting
- Comprehensive logging
- Error handling

## Security Requirements
- Password hashing with bcrypt
- JWT token expiration
- Input validation to prevent injection
- CORS configuration
- Rate limiting per user
- HTTPS enforcement in production
EOF

    echo "✅ Project requirements created"
}

# Run the complete workflow
run_complete_workflow() {
    echo "🚀 Starting complete multi-agent workflow..."
    echo ""
    
    # Phase 1: Discovery (Gemini + Ollama)
    echo "🔍 PHASE 1: DISCOVERY"
    echo "Agents: Gemini (research), Ollama (quick analysis)"
    echo ""
    
    "$WORKFLOW_SCRIPT" discovery &
    DISCOVERY_PID=$!
    
    # Wait for discovery to complete
    wait $DISCOVERY_PID
    echo "✅ Discovery phase completed"
    sleep 2
    
    # Phase 2: Planning (OpenCode + Gemini)
    echo ""
    echo "📋 PHASE 2: PLANNING"
    echo "Agents: OpenCode (architecture), Gemini (validation)"
    echo ""
    
    "$WORKFLOW_SCRIPT" planning &
    PLANNING_PID=$!
    
    # Wait for planning to complete
    wait $PLANNING_PID
    echo "✅ Planning phase completed"
    sleep 2
    
    # Phase 3: Implementation (Codex + Ollama)
    echo ""
    echo "💻 PHASE 3: IMPLEMENTATION"
    echo "Agents: Codex (code generation), Ollama (boilerplate)"
    echo ""
    
    "$WORKFLOW_SCRIPT" implementation &
    IMPLEMENTATION_PID=$!
    
    # Wait for implementation to complete
    wait $IMPLEMENTATION_PID
    echo "✅ Implementation phase completed"
    sleep 2
    
    # Phase 4: Testing (Gemini + Codex)
    echo ""
    echo "🧪 PHASE 4: TESTING"
    echo "Agents: Gemini (test strategy), Codex (test implementation)"
    echo ""
    
    "$WORKFLOW_SCRIPT" testing &
    TESTING_PID=$!
    
    # Wait for testing to complete
    wait $TESTING_PID
    echo "✅ Testing phase completed"
    sleep 2
    
    # Phase 5: Deployment (Fabric + OpenCode)
    echo ""
    echo "🚢 PHASE 5: DEPLOYMENT"
    echo "Agents: Fabric (documentation), OpenCode (pipeline)"
    echo ""
    
    "$WORKFLOW_SCRIPT" deployment &
    DEPLOYMENT_PID=$!
    
    # Wait for deployment to complete
    wait $DEPLOYMENT_PID
    echo "✅ Deployment phase completed"
    
    echo ""
    echo "🎉 Complete workflow finished!"
}

# Monitor and show results
show_results() {
    echo ""
    echo "================================================"
    echo "  WORKFLOW RESULTS"
    echo "================================================"
    
    local artifacts_dir="$PROJECT_DIR/.forge/artifacts"
    local logs_dir="$PROJECT_DIR/.forge/logs"
    
    echo ""
    echo "📁 Generated Artifacts:"
    ls -la "$artifacts_dir" | grep -E "\.(md|js|json|yaml|sql)$" | tail -10
    
    echo ""
    echo "📊 Agent Performance:"
    for agent in gemini codex opencode fabric ollama; do
        local log_count=$(ls "$logs_dir"/${agent}_*.log 2>/dev/null | wc -l)
        echo "  $agent: $log_count completed tasks"
    done
    
    echo ""
    echo "📋 Key Deliverables:"
    
    # Show discovery results
    local discovery_file=$(ls "$artifacts_dir"/discovery_*.md 2>/dev/null | tail -1)
    if [[ -n "$discovery_file" ]]; then
        echo "  🔍 Discovery Report: $discovery_file"
        echo "     Key findings: $(grep -c "##" "$discovery_file" 2>/dev/null || echo "0") sections"
    fi
    
    # Show planning results
    local planning_file=$(ls "$artifacts_dir"/planning_*.md 2>/dev/null | tail -1)
    if [[ -n "$planning_file" ]]; then
        echo "  📋 Project Plan: $planning_file"
        echo "     Tasks identified: $(grep -c "###" "$planning_file" 2>/dev/null || echo "0")"
    fi
    
    # Show implementation results
    local impl_file=$(ls "$artifacts_dir"/implementation_*.md 2>/dev/null | tail -1)
    if [[ -n "$impl_file" ]]; then
        echo "  💻 Implementation: $impl_file"
        echo "     Components built: $(grep -c "component\|module\|service" "$impl_file" 2>/dev/null || echo "0")"
    fi
    
    # Show testing results
    local test_file=$(ls "$artifacts_dir"/testing_*.md 2>/dev/null | tail -1)
    if [[ -n "$test_file" ]]; then
        echo "  🧪 Test Suite: $test_file"
        echo "     Test cases: $(grep -c "test\|spec" "$test_file" 2>/dev/null || echo "0")"
    fi
    
    # Show deployment results
    local deploy_file=$(ls "$artifacts_dir"/deployment_*.md 2>/dev/null | tail -1)
    if [[ -n "$deploy_file" ]]; then
        echo "  🚢 Deployment: $deploy_file"
        echo "     Deployment steps: $(grep -c "step\|action\|command" "$deploy_file" 2>/dev/null || echo "0")"
    fi
    
    echo ""
    echo "📄 View detailed logs: $logs_dir/"
    echo "📦 All artifacts: $artifacts_dir/"
}

# Main execution
main() {
    echo "Starting Multi-Agent Headless Workflow Example..."
    echo ""
    
    # Create requirements
    create_project_requirements
    
    # Run complete workflow
    run_complete_workflow
    
    # Show results
    show_results
    
    echo ""
    echo "================================================"
    echo "  EXAMPLE COMPLETED SUCCESSFULLY!"
    echo "================================================"
    echo ""
    echo "What happened:"
    echo "  • 5 AI agents worked in parallel background processes"
    echo "  • Each agent used their specialized capabilities"
    echo "  • No interactive prompts - fully automated"
    echo "  • All collaboration happened via shared files"
    echo ""
    echo "Agent Specialization Used:"
    echo "  • Gemini: Research, discovery, validation, testing strategy"
    echo "  • Codex: Code generation, implementation, test code"
    echo "  • OpenCode: Architecture, planning, deployment pipeline"
    echo "  • Fabric: Documentation, patterns, deployment configs"
    echo "  • Ollama: Quick tasks, boilerplate, preliminary analysis"
    echo ""
    echo "To create your own project:"
    echo "  1. Modify config.yaml for your project needs"
    echo "  2. Update project_requirements.md with your specs"
    echo "  3. Run: ./headless_workflow.sh all"
    echo ""
}

# Run if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi