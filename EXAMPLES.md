# Fabric-Lite Examples & Workflows

## Quick Start Examples

### 1. Basic Tool Usage
```bash
# Test a simple math problem
./bin/fabric-lite run ollama -P "What is 15+27?" --provider ollama

# Get coding help
./bin/fabric-lite run codex -P "Create a Python function to reverse a string" --provider openai

# Research assistance
./bin/fabric-lite run gemini -P "Explain quantum computing in simple terms" --provider gemini
```

### 2. Pattern-Based Workflows
```bash
# Summarize text using built-in pattern
echo "Long text content..." | ./bin/fabric-lite --pattern summarize --provider openai

# Extract ideas from content
echo "Meeting notes..." | ./bin/fabric-lite --pattern extract_ideas --provider anthropic

# Explain code
cat main.go | ./bin/fabric-lite --pattern explain_code --provider openai
```

### 3. YouTube Content Analysis
```bash
# Analyze YouTube transcript for ADHD-friendly summary
python3 tools/youtube_analyzer.py transcripts/video.srt

# Save analysis to file
python3 tools/youtube_analyzer.py transcripts/video.srt --output analysis.md

# Get JSON output for automation
python3 tools/youtube_analyzer.py transcripts/video.srt --json > analysis.json
```

### 4. ADHD-Optimized Workflows
```bash
# Start ADHD daemon for background processing
./scripts/adhd_daemon.sh start

# Check daemon status
./scripts/adhd_daemon.sh --status

# Stop daemon when done
./scripts/adhd_daemon.sh stop
```

## Advanced Workflows

### 5. Multi-Tool Coordination
```bash
# Research → Design → Implementation workflow
./bin/fabric-lite run gemini -P "Research REST API best practices" --provider gemini
./bin/fabric-lite run codex -P "Design API structure based on research" --provider openai  
./bin/fabric-lite run opencode -P "Implement the API endpoints" --provider openai
```

### 6. Error Recovery & Fallbacks
```bash
# This will automatically try alternative providers if one fails
./bin/fabric-lite run codex -P "Complex task requiring reliable execution" --provider ollama

# Provider-specific execution
./bin/fabric-lite run codex -P "Task requiring advanced reasoning" --provider anthropic
```

### 7. Session Management
```bash
# Save session context
./bin/fabric-lite run codex -P "Project planning" --save-session --provider openai

# Use saved context for follow-up tasks
./bin/fabric-lite run codex -P "Continue from previous session" --context ./user_context/session_* --provider openai
```

## Real-World Use Cases

### Development Workflow
```bash
# 1. Plan project
./bin/fabric-lite run gemini -P "Plan a microservices architecture for e-commerce" --provider gemini

# 2. Generate code structure  
./bin/fabric-lite run codex -P "Create Go service structure for user management" --provider openai

# 3. Test implementation
./bin/fabric-lite run codex -P "Write unit tests for the user service" --provider openai

# 4. Document API
./bin/fabric-lite run opencode -P "Generate OpenAPI documentation" --provider openai
```

### Content Analysis Pipeline
```bash
# 1. Process YouTube content
python3 tools/youtube_analyzer.py transcripts/lecture.srt --output lecture_analysis.md

# 2. Extract key concepts
cat lecture_analysis.md | ./bin/fabric-lite --pattern extract_ideas --provider openai

# 3. Create study guide
./bin/fabric-lite run gemini -P "Create ADHD-friendly study guide from these concepts" --provider gemini --context lecture_analysis.md
```

### Automated Workflow with Daemon
```bash
# Start continuous monitoring
./scripts/adhd_daemon.sh start

# Daemon will automatically:
# - Monitor for new transcripts
# - Analyze content with appropriate tools  
# - Save results to ./user_context/
# - Maintain session state across restarts
```

## Troubleshooting Examples

### Provider Failures
```bash
# If OpenAI fails, try local alternative
./bin/fabric-lite run codex -P "Help debug this issue" --provider ollama

# Check which providers are configured
./bin/fabric-lite config
```

### Tool Errors
```bash
# Verbose output for debugging
./bin/fabric-lite run codex -P "Test task" --provider openai --verbose

# Check pattern availability
./bin/fabric-lite list
```

These examples demonstrate the key capabilities and workflows available in fabric-lite v2.0!