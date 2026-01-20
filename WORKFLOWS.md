# Multi-Tool Coordination Workflows

This guide demonstrates how to use multiple AI tools in coordination for complex tasks, leveraging the combined capabilities of fabric-lite's tool ecosystem.

## Workflow Principles

1. **Tool Specialization** - Use each tool for its strengths
2. **Sequential Coordination** - Pass outputs between tools as inputs
3. **Background Processing** - Execute tools in headless mode for batch operations
4. **State Persistence** - Save intermediate results for resumption
5. **Error Recovery** - Implement fallback mechanisms when tools fail

## Workflow Examples

### Workflow 1: Research → Analysis → Implementation

**Use Case**: Analyze competitor's API and implement improved version

```bash
# Step 1: Research with gemini
fabric-lite run gemini -P "analyze competitor's API features, pricing, and technical architecture" > research/gemini_analysis.md

# Step 2: Architecture design with claude  
fabric-lite run claude -P "design improved API architecture based on research findings" --input research/gemini_analysis.md > architecture/claude_design.md

# Step 3: Implementation with codex
fabric-lite run codex -P "implement REST API with following architecture: $(cat architecture/claude_design.md)" > implementation/codex_implementation.go

# Step 4: Local testing with ollama
fabric-lite run ollama -P "validate and test the implemented API" --input implementation/codex_implementation.go
```

### Workflow 2: Documentation Generation Pipeline

**Use Case**: Generate comprehensive documentation from source code

```bash
# Extract features with gemini
fabric-lite run gemini -P "extract key features from source code" --input src/main.go > features/gemini_features.md

# Structure documentation with claude
fabric-lite run claude -P "organize features into structured documentation" --input features/gemini_features.md > docs/api_reference.md

# Generate examples with codex
fabric-lite run codex -P "create usage examples for the API" --input docs/api_reference.md > examples/codex_examples.md

# Review and refine with opencode (interactive)
fabric-lite run opencode -P "review and improve documentation" --input docs/api_reference.md examples/codex_examples.md
```

### Workflow 3: Multi-Tool Code Review

**Use Case**: Comprehensive code review using multiple perspectives

```bash
# Security analysis with gemini
fabric-lite run gemini -P "analyze security vulnerabilities and best practices" --input src/ > security/gemini_analysis.md

# Architecture review with claude
fabric-lite run claude -P "review code architecture and design patterns" --input src/ > architecture/claude_review.md

# Code quality with codex
fabric-lite run codex -P "analyze code quality, test coverage, and performance" --input src/ > quality/codex_analysis.md

# Local testing with ollama
fabric-lite run ollama -P "run comprehensive tests and validate functionality" --input src/ > testing/ollama_results.md

# Compile results
cat security/gemini_analysis.md architecture/claude_review.md quality/codex_analysis.md testing/ollama_results.md > review/comprehensive_report.md
```

### Workflow 4: Background Batch Processing

**Use Case**: Process multiple files simultaneously using tool coordination

```bash
#!/bin/bash
# Background processing script
input_dir="./src"
output_dir="./processed"

# Process each file
for file in "$input_dir"/*.go; do
    echo "Processing $file with gemini..."
    fabric-lite run gemini -P "extract key information and summarize" --input "$file" > "$output_dir/$(basename "$file" .go)_summary.md"
    
    echo "Analyzing $file with codex..."
    fabric-lite run codex -P "suggest improvements and optimizations" --input "$file" >> "$output_dir/$(basename "$file" .go)_improvements.md"
done

echo "Batch processing completed. Results in $output_dir"
```

### Workflow 5: Headless Execution Protocols

**Use Case**: Execute long-running tasks in background

```bash
# Start background research task
nohup fabric-lite run gemini -P "comprehensive market research" > research/market_analysis.md 2>&1 &

# Monitor progress
tail -f research/market_analysis.md &

# Check completion when process finishes
if pgrep -f "fabric-lite run gemini" > /dev/null; then
    echo "Research task completed"
    cat research/market_analysis.md
fi
```

## Tool Coordination Patterns

### Sequential Processing
- Output of Tool A becomes input for Tool B
- Use intermediate files to pass state between tools
- Validate each step before proceeding

### Parallel Processing  
- Run multiple tools simultaneously on different aspects
- Combine results using a coordination tool (codex/claude)
- Use fabric-lite patterns to standardize output formats

### Feedback Loops
- Use tool outputs to refine subsequent tool inputs
- Iterate between tools until quality threshold is met
- Use ollama for quick validation of intermediate results

## Error Handling Strategies

### Tool Unavailability
```bash
# Check tool availability before execution
if ! fabric-lite run gemini --check-available 2>/dev/null; then
    echo "Gemini unavailable, using claude fallback..."
    fabric-lite run claude -P "$PROMPT"
else
    fabric-lite run gemini -P "$PROMPT"
fi
```

### Partial Failures
```bash
# Save intermediate state
fabric-lite run codex -P "implement feature A" > state/feature_a_progress.md

# If fails, resume from checkpoint
if [ $? -ne 0 ]; then
    fabric-lite run codex -P "continue from checkpoint" --input state/feature_a_progress.md --resume-checkpoint state/checkpoint_a.yaml
fi
```

## Performance Optimization

### Batch Operations
- Process multiple items in single tool call when possible
- Use local tools (ollama) for rapid iteration
- Cache results to avoid redundant API calls

### Resource Management
- Use different providers for parallel operations
- Implement rate limiting for cloud providers
- Prioritize local tools for sensitive data

## Integration with Forge

### Phase Mapping
- Discovery: gemini (research) + fabric-lite (documentation)
- Planning: claude (architecture) + opencode (interactive design)
- Design: claude (detailed design) + codex (implementation planning)
- Implementation: codex (coding) + ollama (local testing)
- Testing: gemini (test planning) + ollama (quick validation)
- Deployment: fabric (documentation) + all tools (final validation)

### Checkpoint Integration
```bash
# Forge automatically tracks tool usage in .forge/state.yaml
# Each phase completion creates checkpoints
# Resume from any checkpoint with: forge phase resume --from checkpoint_name
```

These workflows demonstrate how to leverage the full fabric-lite tool ecosystem for complex, multi-stage tasks while maintaining resumable progress and error recovery.