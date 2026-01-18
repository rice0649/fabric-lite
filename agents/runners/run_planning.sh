#!/bin/bash
# Planning Pipeline Runner
# Run this in a separate terminal
# Executes: Planning → Review → Code Analysis → Final Review

set -e

PROJECT_DIR="/home/oak38/projects/fabric-lite"
FABRIC_DIR="/home/oak38/projects/fabric"
OUTPUT_DIR="$PROJECT_DIR/agents/outputs"

echo "============================================"
echo "  Planning Pipeline"
echo "============================================"
echo ""
echo "This runs the full planning pipeline:"
echo "  1. Planning Agent → 01_planning.md"
echo "  2. Review Agent → 02_review.md"
echo "  3. Code Analysis → 03_code_analysis.md"
echo "  4. Final Review → 04_final_spec.md + FINAL_SUMMARY.md"
echo ""
echo "============================================"
echo ""

# Ensure output directory exists
mkdir -p "$OUTPUT_DIR"

# Change to fabric directory for analysis
cd "$FABRIC_DIR"

echo "Starting Gemini Code session..."
echo ""
echo "Instructions for Gemini:"
echo ""
echo "Execute these agents IN ORDER:"
echo ""
echo "STEP 1: Read agents/planning_agent.md (from fabric-lite)"
echo "        Analyze the Fabric codebase"
echo "        Write output to: $OUTPUT_DIR/01_planning.md"
echo ""
echo "STEP 2: Read agents/review_agent.md"
echo "        Review the planning output"
echo "        Write output to: $OUTPUT_DIR/02_review.md"
echo ""
echo "STEP 3: Read agents/code_analysis_agent.md"
echo "        Analyze key source files"
echo "        Write output to: $OUTPUT_DIR/03_code_analysis.md"
echo ""
echo "STEP 4: Read agents/final_review_agent.md"
echo "        Consolidate all outputs"
echo "        Write to: $OUTPUT_DIR/04_final_spec.md"
echo "        Write summary to: $OUTPUT_DIR/FINAL_SUMMARY.md"
echo ""

# Launch Gemini
gemini "You are running the Planning Pipeline. Start with $PROJECT_DIR/agents/planning_agent.md. Execute each agent in sequence, writing outputs to $OUTPUT_DIR/. Finish by creating FINAL_SUMMARY.md."
