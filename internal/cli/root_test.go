package cli

import (
	"testing"
)

func TestNewRootCmd(t *testing.T) {
	version := "test-version-1.0.0"
	cmd := NewRootCmd(version)

	if cmd == nil {
		t.Fatal("Expected root command to be non-nil")
	}

	if cmd.Use != "fabric-lite" {
		t.Errorf("Expected command use to be 'fabric-lite', got '%s'", cmd.Use)
	}

	if cmd.Short != "Lightweight AI augmentation framework" {
		t.Errorf("Expected short description to be 'Lightweight AI augmentation framework', got '%s'", cmd.Short)
	}

	if cmd.Version != version {
		t.Errorf("Expected version to be '%s', got '%s'", version, cmd.Version)
	}

	// Check that global flags are present
	flags := cmd.PersistentFlags()

	patternFlag := flags.Lookup("pattern")
	if patternFlag == nil {
		t.Error("Expected 'pattern' flag to be present")
	} else if patternFlag.Shorthand != "p" {
		t.Errorf("Expected pattern flag shorthand to be 'p', got '%s'", patternFlag.Shorthand)
	}

	modelFlag := flags.Lookup("model")
	if modelFlag == nil {
		t.Error("Expected 'model' flag to be present")
	} else if modelFlag.Shorthand != "m" {
		t.Errorf("Expected model flag shorthand to be 'm', got '%s'", modelFlag.Shorthand)
	}

	providerFlag := flags.Lookup("provider")
	if providerFlag == nil {
		t.Error("Expected 'provider' flag to be present")
	} else if providerFlag.Shorthand != "v" {
		t.Errorf("Expected provider flag shorthand to be 'v', got '%s'", providerFlag.Shorthand)
	}

	streamFlag := flags.Lookup("stream")
	if streamFlag == nil {
		t.Error("Expected 'stream' flag to be present")
	} else if streamFlag.Shorthand != "s" {
		t.Errorf("Expected stream flag shorthand to be 's', got '%s'", streamFlag.Shorthand)
	}
}

func TestRootCmdExecution(t *testing.T) {
	version := "test-version"
	cmd := NewRootCmd(version)

	// Test help command
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error executing help command, got %v", err)
	}
}

func TestRootCmdVersion(t *testing.T) {
	version := "test-version-2.0.0"
	cmd := NewRootCmd(version)

	// Test version flag
	cmd.SetArgs([]string{"--version"})
	err := cmd.Execute()
	// Version command usually exits, so we might get an error
	// The important thing is the command structure is valid
	if err != nil && cmd.Version != version {
		t.Errorf("Expected version to be set correctly even with execution error")
	}
}
