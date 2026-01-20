package providers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ExecutableProvider implements Provider for custom executable scripts
type ExecutableProvider struct {
	name       string
	executable string
	args       []string
	env        map[string]string
	workDir    string
	timeout    time.Duration
}

// NewExecutableProvider creates a new executable-based provider
func NewExecutableProvider(name string, config map[string]any) (*ExecutableProvider, error) {
	executable := getConfigString(config, "executable", "")
	if executable == "" {
		return nil, fmt.Errorf("executable path is required for executable provider")
	}

	// Expand ~ to home directory
	if strings.HasPrefix(executable, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			executable = filepath.Join(home, executable[2:])
		}
	}

	// Parse args
	var args []string
	if argsVal, ok := config["args"].([]any); ok {
		for _, arg := range argsVal {
			if s, ok := arg.(string); ok {
				args = append(args, s)
			}
		}
	}

	// Parse environment variables
	env := make(map[string]string)
	if envVal, ok := config["env"].(map[string]any); ok {
		for k, v := range envVal {
			if s, ok := v.(string); ok {
				env[k] = s
			}
		}
	}

	workDir := getConfigString(config, "work_dir", "")
	timeout := time.Duration(getConfigInt(config, "timeout_seconds", 300)) * time.Second

	return &ExecutableProvider{
		name:       name,
		executable: executable,
		args:       args,
		env:        env,
		workDir:    workDir,
		timeout:    timeout,
	}, nil
}

func (p *ExecutableProvider) Name() string {
	return p.name
}

func (p *ExecutableProvider) IsAvailable() bool {
	// Check if executable exists and is executable
	info, err := os.Stat(p.executable)
	if err != nil {
		// Try to find in PATH
		_, err = exec.LookPath(p.executable)
		return err == nil
	}
	// Check if it's executable
	return info.Mode()&0111 != 0
}

func (p *ExecutableProvider) GetModels() []string {
	// Executable providers typically don't have models
	return []string{"default"}
}

func (p *ExecutableProvider) Execute(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	start := time.Now()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	// Build command
	args := append(p.args, request.Prompt)
	cmd := exec.CommandContext(ctx, p.executable, args...)

	// Set environment
	cmd.Env = os.Environ()
	for k, v := range p.env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Add system prompt as environment variable if present
	if request.System != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("FABRIC_SYSTEM=%s", request.System))
	}

	// Set working directory
	if p.workDir != "" {
		cmd.Dir = p.workDir
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Provide input via stdin
	cmd.Stdin = strings.NewReader(request.Prompt)

	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("execution timed out after %v", p.timeout)
		}
		return nil, fmt.Errorf("execution failed: %w\nstderr: %s", err, stderr.String())
	}

	return &CompletionResponse{
		Content:  stdout.String(),
		Model:    "executable:" + filepath.Base(p.executable),
		Duration: time.Since(start),
	}, nil
}

func (p *ExecutableProvider) ExecuteStream(ctx context.Context, request CompletionRequest) (<-chan StreamChunk, error) {
	// Executable providers don't support streaming natively
	// Return the full result as a single chunk
	chunks := make(chan StreamChunk, 1)

	go func() {
		defer close(chunks)
		resp, err := p.Execute(ctx, request)
		if err != nil {
			chunks <- StreamChunk{Error: err, Done: true}
			return
		}
		chunks <- StreamChunk{Content: resp.Content, Done: true}
	}()

	return chunks, nil
}
