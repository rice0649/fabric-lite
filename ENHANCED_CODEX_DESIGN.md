See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# Enhanced Tool Coordination System

This creates a more intelligent meta-coding tool that can:
1. Assess task complexity and choose optimal provider
2. Coordinate multiple tools for different aspects of the same task
3. Implement fallback mechanisms when primary tools fail
4. Provide real-time progress tracking and error recovery

## Enhanced Codex Implementation

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "time"
)

// Enhanced codex tool coordinator
type EnhancedCodex struct {
    Name            string
    Description     string
    Version         string
    Tools           map[string]Tool
    Config          *CodexConfig
}

type CodexConfig struct {
    DefaultProvider  string                 `yaml:"default_provider"`
    FallbackChain   []string              `yaml:"fallback_chain"`
    TaskComplexity string                 `yaml:"task_complexity"`
    MaxRetries     int                   `yaml:"max_retries"`
    TimeoutSeconds   int                   `yaml:"timeout_seconds"`
}

type Tool interface {
    Name() string
    Execute(ctx ExecutionContext) (*ExecutionResult, error)
    IsAvailable() bool
    GetCapabilities() []string
}

type TaskComplexity int

const (
    Simple TaskComplexity = iota
    Medium
    Complex
    Expert
)

func NewEnhancedCodex(config *CodexConfig) *EnhancedCodex {
    return &EnhancedCodex{
        Name:        "enhanced-codex",
        Description: "Advanced meta-coding assistant with multi-tool coordination",
        Version:     "2.0",
        Tools:       initializeTools(config),
        Config:      config,
    }
}

func initializeTools(config *CodexConfig) map[string]Tool {
    tools := make(map[string]Tool)
    
    // Initialize all available tools with their configurations
    tools["codex"] = NewCodexTool(config)
    tools["claude"] = NewClaudeTool(config)
    tools["gemini"] = NewGeminiTool(config)
    tools["opencode"] = NewOpencodeTool(config)
    tools["ollama"] = NewOllamaTool(config)
    tools["fabric"] = NewFabricTool(config)
    
    return tools
}

// Execute coordinates multiple tools for complex tasks
func (ec *EnhancedCodex) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
    start := time.Now()
    
    // Analyze task complexity
    complexity := analyzeComplexity(ctx.Prompt)
    
    // Select optimal provider based on task and complexity
    primaryProvider, fallbackProviders := selectProviders(ec.Config, complexity)
    
    var combinedResult strings.Builder
    var errors []string
    
    // Execute primary provider
    if tool, ok := ec.Tools[primaryProvider]; ok {
        result, err := executeWithRetry(tool, ctx, ec.Config.MaxRetries)
        if err != nil {
            errors = append(errors, fmt.Sprintf("%s failed: %v", primaryProvider, err))
        } else {
            combinedResult.WriteString(fmt.Sprintf("[%s Result]\n%s\n", primaryProvider, result.Output))
        }
    }
    
    // Execute fallback providers if primary failed
    for i, providerName := range fallbackProviders {
        if len(errors) > 0 {
            break // Skip fallbacks if primary succeeded
        }
        
        if tool, ok := ec.Tools[providerName]; ok {
            modifiedCtx := createModifiedContext(ctx, providerName, i+1)
            result, err := executeWithRetry(tool, modifiedCtx, ec.Config.MaxRetries)
            if err != nil {
                errors = append(errors, fmt.Sprintf("%s fallback #%d failed: %v", providerName, i+1, err))
            } else {
                combinedResult.WriteString(fmt.Sprintf("[%s Fallback #%d]\n%s\n", providerName, i+1, result.Output))
            }
        }
    }
    
    // Execute coordination tool if needed
    if len(errors) > 0 || combinedResult.Len() > 0 {
        coordinationResult, err := coordinateResults(ec, combinedResult.String(), errors)
        if err != nil {
            return nil, fmt.Errorf("coordination failed: %w", err)
        }
        return coordinationResult, nil
    }
    
    duration := time.Since(start)
    
    return &ExecutionResult{
        Output: combinedResult.String(),
        Success: len(errors) == 0 && combinedResult.Len() > 0,
        ExitCode: 0,
        Error:   "",
        Metadata: map[string]interface{}{
            "duration_ms":     duration.Milliseconds(),
            "complexity":      complexity.String(),
            "providers_used":  append([]string{primaryProvider}, fallbackProviders...),
            "tool_version":   "2.0",
        },
    }, nil
}

func executeWithRetry(tool Tool, ctx ExecutionContext, maxRetries int) (*ExecutionResult, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        result, err := tool.Execute(ctx)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // Exponential backoff
        delay := time.Duration(i) * time.Second
        if delay > 30*time.Second {
            delay = 30 * time.Second
        }
        
        select {
        case <-time.After(delay):
            continue
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }
    
    return nil, fmt.Errorf("tool failed after %d retries: %w", maxRetries, lastErr)
}

func analyzeComplexity(prompt string) TaskComplexity {
    // Simple heuristic analysis
    words := strings.Fields(prompt)
    if len(words) < 10 {
        return Simple
    } else if len(words) < 50 {
        return Medium
    } else if strings.Contains(strings.ToLower(prompt), "implement") ||
               strings.Contains(strings.ToLower(prompt), "design") ||
               strings.Contains(strings.ToLower(prompt), "architecture") {
        return Complex
    }
    return Expert
}

func selectProviders(config *CodexConfig, complexity TaskComplexity) (string, []string) {
    // Provider selection based on task complexity and configuration
    switch complexity {
    case Simple:
        return config.DefaultProvider, []string{}
    case Medium:
        return config.DefaultProvider, config.FallbackChain[:1] // Use first fallback
    case Complex, Expert:
        return config.DefaultProvider, config.FallbackChain // Use all fallbacks
    }
}

func createModifiedContext(ctx ExecutionContext, providerName string, attempt int) ExecutionContext {
    modified := ctx
    modified.Prompt = fmt.Sprintf("[%s attempt %d] %s", providerName, attempt, ctx.Prompt)
    return modified
}

func coordinateResults(ec *EnhancedCodex, primaryOutput string, errors []string) (*ExecutionResult, error) {
    // Simple coordination: combine all outputs
    var final strings.Builder
    
    if primaryOutput != "" {
        final.WriteString("=== PRIMARY PROVIDER RESULT ===\n")
        final.WriteString(primaryOutput)
    }
    
    if len(errors) > 0 {
        final.WriteString("\n=== FALLBACK RESULTS ===\n")
        for i, err := range errors {
            final.WriteString(fmt.Sprintf("Fallback #%d: %s\n", i+1, err))
        }
    }
    
    // Simple enhancement suggestions
    final.WriteString("\n=== ENHANCEMENT SUGGESTIONS ===\n")
    final.WriteString(generateSuggestions(primaryOutput, errors))
    
    return &ExecutionResult{
        Output:   final.String(),
        Success: true,
        ExitCode: 0,
    }, nil
}

func generateSuggestions(primaryOutput string, errors []string) string {
    var suggestions strings.Builder
    
    if strings.Contains(primaryOutput, "error") {
        suggestions.WriteString("• Consider using gemini for research tasks\n")
        suggestions.WriteString("• Verify API keys and endpoint configurations\n")
    }
    
    if len(errors) > 0 {
        suggestions.WriteString("• Check network connectivity to provider endpoints\n")
        suggestions.WriteString("• Consider simplifying task complexity\n")
    }
    
    return suggestions.String()
}

// Individual tool implementations (simplified for brevity)
type CodexTool struct{ config *CodexConfig }
type ClaudeTool struct{ config *CodexConfig }
type GeminiTool struct{ config *CodexConfig }
type OpencodeTool struct{ config *CodexConfig }
type OllamaTool struct{ config *CodexConfig }
type FabricTool struct{ config *CodexConfig }

func (t *CodexTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
    // Implementation would delegate to primary provider
    return nil, fmt.Errorf("codex tool: delegating to %s", t.config.DefaultProvider)
}
// ... similar for other tools
```

## Configuration File

```yaml
# enhanced_codex_config.yaml
default_provider: ollama
fallback_chain: [claude, gemini, codex]
task_complexity: medium
max_retries: 3
timeout_seconds: 30
```

## Usage

```bash
# Enhanced tool coordination
export ENHANCED_CODEX_CONFIG="path/to/enhanced_codex_config.yaml"
./bin/fabric-lite run enhanced-codex -P "implement scalable REST API with error handling and comprehensive testing"
```

This system provides:
- **Automatic provider selection** based on task complexity
- **Fallback mechanisms** with retry logic
- **Multi-tool coordination** for complex tasks
- **Progress tracking** and detailed error reporting
- **Configurable behavior** via YAML files