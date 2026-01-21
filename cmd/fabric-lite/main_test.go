package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rice0649/fabric-lite/internal/cli"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close stdout: %v", err)
	}
	os.Stdout = origStdout

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	return string(out)
}

func TestRunVersionFlag(t *testing.T) {
	origArgs := os.Args
	origVersion := Version
	t.Cleanup(func() {
		os.Args = origArgs
		Version = origVersion
	})

	Version = "test-version"
	os.Args = []string{"fabric-lite", "--version"}

	output := captureStdout(t, func() {
		rootCmd := cli.NewRootCmd(Version)
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("run: %v", err)
		}
	})

	if !strings.Contains(output, "fabric-lite version test-version") {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestRunDefaultUsage(t *testing.T) {
	origArgs := os.Args
	origVersion := Version
	t.Cleanup(func() {
		os.Args = origArgs
		Version = origVersion
	})

	Version = "test-version"
	os.Args = []string{"fabric-lite"}

	output := captureStdout(t, func() {
		rootCmd := cli.NewRootCmd(Version)
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("run: %v", err)
		}
	})

	if !strings.Contains(output, "Usage:") {
		t.Fatalf("expected usage output, got: %q", output)
	}
	if !strings.Contains(output, "fabric-lite [command]") {
		t.Fatalf("expected command usage in output, got: %q", output)
	}
}

func TestRunListCommand(t *testing.T) {
	origArgs := os.Args
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() {
		os.Args = origArgs
		if chdirErr := os.Chdir(origDir); chdirErr != nil {
			t.Fatalf("restore cwd: %v", chdirErr)
		}
	})

	repoRoot := filepath.Dir(filepath.Dir(origDir))
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("chdir to repo root: %v", err)
	}

	// Use 'list' subcommand instead of '--list' flag
	os.Args = []string{"fabric-lite", "list"}

	output := captureStdout(t, func() {
		rootCmd := cli.NewRootCmd(Version)
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("run: %v", err)
		}
	})

	if !strings.Contains(output, "Available patterns:") {
		t.Fatalf("expected patterns header, got: %q", output)
	}
}
