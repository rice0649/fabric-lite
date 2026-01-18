#!/bin/bash
# Implementation Agent Runner
# Run this AFTER the planning pipeline completes

set -e

PROJECT_DIR="/home/oak38/projects/fabric-lite"
OUTPUT_DIR="$PROJECT_DIR/agents/outputs"
SPEC_FILE="$OUTPUT_DIR/04_final_spec.md"
PATTERNS_QUEUE="$PROJECT_DIR/agents/queue/patterns_to_implement.md"

echo "============================================"
echo "  Implementation Agent"
echo "============================================"
echo ""
echo "Project: $PROJECT_DIR"
echo "Spec:    $SPEC_FILE"
echo "Patterns: $PATTERNS_QUEUE"
echo ""
echo "============================================"
echo ""

# Check if planning is complete
if [ ! -f "$SPEC_FILE" ]; then
    echo "ERROR: Final spec not found at $SPEC_FILE"
    echo "Run the planning pipeline first: ./run_planning.sh"
    exit 1
fi

if [ ! -f "$OUTPUT_DIR/FINAL_SUMMARY.md" ]; then
    echo "WARNING: FINAL_SUMMARY.md not found"
    echo "Planning may not be complete"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Change to project directory
cd "$PROJECT_DIR"

echo "Starting Gemini Code session..."
echo ""
echo "Instructions for Gemini:"
echo "1. Read your agent instructions: agents/implementation_agent.md"
echo "2. Read the final spec: $SPEC_FILE"
echo "3. Check patterns queue: $PATTERNS_QUEUE"
echo "4. Implement Phase 1 → 2 → 3 → 4"
echo "5. Commit after each phase"
echo "6. Update agents/outputs/implementation_progress.md"
echo ""

# Launch Gemini
gemini "You are the Implementation Agent. Read agents/implementation_agent.md for instructions. Read the spec from $SPEC_FILE. Begin implementing fabric-lite core functionality."
