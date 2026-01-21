package cli

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewRootCmdValidation(t *testing.T) {
	version := "1.0.0-test"
	cmd := NewRootCmd(version)

	// Test that it's a proper cobra command
	if cmd == nil {
		t.Fatal("Expected command to be non-nil")
	}

	// Test command properties
	if cmd.Use == "" {
		t.Error("Expected command use to be set")
	}
	if cmd.Short == "" {
		t.Error("Expected command short description to be set")
	}
	if cmd.Long == "" {
		t.Error("Expected command long description to be set")
	}
	if cmd.Version != version {
		t.Errorf("Expected version '%s', got '%s'", version, cmd.Version)
	}

	// Test that it has subcommands
	if len(cmd.Commands()) == 0 {
		t.Error("Expected root command to have subcommands")
	}
}

func TestRootCmdFlags(t *testing.T) {
	cmd := NewRootCmd("test")

	// Test persistent flags exist
	flags := cmd.PersistentFlags()

	requiredFlags := []string{"pattern", "model", "provider", "stream", "verbose"}
	for _, flagName := range requiredFlags {
		flag := flags.Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to be present", flagName)
		}
	}
}

func TestRootCmdFlagDefaults(t *testing.T) {
	cmd := NewRootCmd("test")
	flags := cmd.PersistentFlags()

	// Test default model value
	modelFlag := flags.Lookup("model")
	if modelFlag != nil && modelFlag.DefValue != "gpt-4o-mini" {
		t.Errorf("Expected default model 'gpt-4o-mini', got '%s'", modelFlag.DefValue)
	}

	// Test default stream value
	streamFlag := flags.Lookup("stream")
	if streamFlag != nil && streamFlag.DefValue != "false" {
		t.Errorf("Expected default stream 'false', got '%s'", streamFlag.DefValue)
	}
}

func TestRootCmdFlagShorthands(t *testing.T) {
	cmd := NewRootCmd("test")
	flags := cmd.PersistentFlags()

	testCases := []struct {
		flagName string
		expected string
	}{
		{"pattern", "p"},
		{"model", "m"},
		{"provider", "v"},
		{"stream", "s"},
		{"verbose", "V"},
	}

	for _, tc := range testCases {
		flag := flags.Lookup(tc.flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to be present", tc.flagName)
			continue
		}
		if flag.Shorthand != tc.expected {
			t.Errorf("Expected flag '%s' shorthand to be '%s', got '%s'", tc.flagName, tc.expected, flag.Shorthand)
		}
	}
}

func TestRootCmdExecuteHelp(t *testing.T) {
	cmd := NewRootCmd("test")

	// Capture output
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"fabric-lite", "--help"}

	err := cmd.Execute()
	// Help command should exit with nil error or help-related error
	if err != nil && !isHelpError(err) {
		t.Errorf("Unexpected error executing help: %v", err)
	}
}

func TestRootCmdExecuteInvalidFlag(t *testing.T) {
	cmd := NewRootCmd("test")

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"fabric-lite", "--invalid-flag"}

	err := cmd.Execute()
	// Should get an error for invalid flag
	if err == nil {
		t.Error("Expected error for invalid flag")
	}
}

func isHelpError(err error) bool {
	if err == nil {
		return false
	}
	// Check if error is related to help (common pattern)
	errStr := err.Error()
	return contains(errStr, "help") || contains(errStr, "usage")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestGetStatusIcon(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		current  bool
		expected string
	}{
		{
			name:     "current phase",
			status:   "in_progress",
			current:  true,
			expected: "▶",
		},
		{
			name:     "completed",
			status:   "completed",
			current:  false,
			expected: "✓",
		},
		{
			name:     "in progress",
			status:   "in_progress",
			current:  false,
			expected: "◐",
		},
		{
			name:     "pending",
			status:   "pending",
			current:  false,
			expected: "○",
		},
		{
			name:     "unknown status",
			status:   "unknown",
			current:  false,
			expected: "○",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStatusIcon(tt.status, tt.current)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "just now",
			duration: 30 * time.Second,
			expected: "just now",
		},
		{
			name:     "5 minutes",
			duration: 5 * time.Minute,
			expected: "5 minutes ago",
		},
		{
			name:     "2 hours",
			duration: 2 * time.Hour,
			expected: "2 hours ago",
		},
		{
			name:     "3 days",
			duration: 72 * time.Hour,
			expected: "3 days ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		n        int
		expected string
	}{
		{
			name:     "pad short string",
			input:    "hello",
			n:        10,
			expected: "hello     ",
		},
		{
			name:     "exact length",
			input:    "hello",
			n:        5,
			expected: "hello",
		},
		{
			name:     "truncate long string",
			input:    "hello world",
			n:        5,
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    "",
			n:        5,
			expected: "     ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padRight(tt.input, tt.n)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
			if len(result) != tt.n {
				t.Errorf("Expected length %d, got %d", tt.n, len(result))
			}
		})
	}
}

func TestNewInitCmd(t *testing.T) {
	cmd := newInitCmd()
	if cmd == nil {
		t.Fatal("Expected init command to be non-nil")
	}
	if cmd.Use != "init" {
		t.Errorf("Expected command use to be 'init', got '%s'", cmd.Use)
	}
}

func TestNewPhaseCmd(t *testing.T) {
	cmd := newPhaseCmd()
	if cmd == nil {
		t.Fatal("Expected phase command to be non-nil")
	}
	if cmd.Use != "phase" {
		t.Errorf("Expected command use to be 'phase', got '%s'", cmd.Use)
	}

	// Test that it has subcommands
	subcommands := cmd.Commands()
	if len(subcommands) == 0 {
		t.Error("Expected phase command to have subcommands")
	}

	// Check for expected subcommands
	expectedSubs := []string{"list", "start", "complete", "info"}
	for _, expected := range expectedSubs {
		found := false
		for _, sub := range subcommands {
			if sub.Use == expected || strings.HasPrefix(sub.Use, expected+" ") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' to be present", expected)
		}
	}
}

func TestNewStatusCmd(t *testing.T) {
	cmd := newStatusCmd()
	if cmd == nil {
		t.Fatal("Expected status command to be non-nil")
	}
	if cmd.Use != "status" {
		t.Errorf("Expected command use to be 'status', got '%s'", cmd.Use)
	}
}

func TestNewSessionCmd(t *testing.T) {
	cmd := newSessionCmd()
	if cmd == nil {
		t.Fatal("Expected session command to be non-nil")
	}
	if cmd.Use != "session" {
		t.Errorf("Expected command use to be 'session', got '%s'", cmd.Use)
	}

	// Test that it has subcommands
	subcommands := cmd.Commands()
	if len(subcommands) == 0 {
		t.Error("Expected session command to have subcommands")
	}

	// Check for expected subcommands
	expectedSubs := []string{"save", "resume", "show"}
	for _, expected := range expectedSubs {
		found := false
		for _, sub := range subcommands {
			if sub.Use == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' to be present", expected)
		}
	}
}

func TestNewAutoCmd(t *testing.T) {
	cmd := newAutoCmd()
	if cmd == nil {
		t.Fatal("Expected auto command to be non-nil")
	}
	if cmd.Use != "auto" {
		t.Errorf("Expected command use to be 'auto', got '%s'", cmd.Use)
	}

	// Test flags
	dryRunFlag := cmd.Flags().Lookup("dry-run")
	if dryRunFlag == nil {
		t.Error("Expected 'dry-run' flag to be present")
	}
}

func TestNewRunCmd(t *testing.T) {
	cmd := newRunCmd()
	if cmd == nil {
		t.Fatal("Expected run command to be non-nil")
	}
	// Check that command use starts with "run"
	if !strings.HasPrefix(cmd.Use, "run") {
		t.Errorf("Expected command use to start with 'run', got '%s'", cmd.Use)
	}
}

func TestNewListCmd(t *testing.T) {
	cmd := newListCmd()
	if cmd == nil {
		t.Fatal("Expected list command to be non-nil")
	}
	if cmd.Use != "list" {
		t.Errorf("Expected command use to be 'list', got '%s'", cmd.Use)
	}
}

func TestNewConfigCmd(t *testing.T) {
	cmd := newConfigCmd()
	if cmd == nil {
		t.Fatal("Expected config command to be non-nil")
	}
	if cmd.Use != "config" {
		t.Errorf("Expected command use to be 'config', got '%s'", cmd.Use)
	}
}

func TestNewVersionCmd(t *testing.T) {
	version := "test-version-1.0.0"
	cmd := newVersionCmd(version)
	if cmd == nil {
		t.Fatal("Expected version command to be non-nil")
	}
	if cmd.Use != "version" {
		t.Errorf("Expected command use to be 'version', got '%s'", cmd.Use)
	}
}
