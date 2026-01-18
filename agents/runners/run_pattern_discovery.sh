#!/bin/bash
# Pattern Discovery Agent Runner
# Run this in a separate terminal

set -e

PROJECT_DIR="/home/oak38/projects/fabric-lite"
AGENT_FILE="$PROJECT_DIR/agents/pattern_discovery_agent.md"
OUTPUT_DIR="$PROJECT_DIR/agents/outputs"
FABRIC_DIR="/home/oak38/projects/fabric"

echo "============================================"
echo "  Pattern Discovery Agent"
echo "============================================"
echo ""
echo "Project: $PROJECT_DIR"
echo "Scanning: $FABRIC_DIR/data/patterns/"
echo "Output:   $OUTPUT_DIR/05_patterns_discovered.md"
echo ""
echo "============================================"
echo ""

# Ensure output directory exists
mkdir -p "$OUTPUT_DIR"
mkdir -p "$PROJECT_DIR/agents/queue"

# Change to fabric directory (where patterns are)
cd "$FABRIC_DIR"

echo "Starting Claude Code session..."
echo ""
echo "Instructions for Claude:"
echo "1. Read the agent file: $AGENT_FILE"
echo "2. Scan all patterns in data/patterns/"
echo "3. Write output to $OUTPUT_DIR/05_patterns_discovered.md"
echo "4. Write queue to $PROJECT_DIR/agents/queue/patterns_to_implement.md"
echo ""

# Launch Claude Code with the agent context
# The user will interact with Claude in this terminal
claude --print "You are the Pattern Discovery Agent. Read your instructions from $AGENT_FILE and begin scanning patterns in data/patterns/. Write your findings to $OUTPUT_DIR/05_patterns_discovered.md"
