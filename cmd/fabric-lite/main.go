package main

import (
	"fmt"
	"os"
)

// Version is set at build time via ldflags
var Version = "0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// TODO: Implement CLI using cobra
	// For now, just print version and usage
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("fabric-lite version %s\n", Version)
			return nil
		case "--help", "-h":
			printUsage()
			return nil
		case "--list", "-l":
			fmt.Println("Available patterns:")
			fmt.Println("  - summarize")
			fmt.Println("\nMore patterns coming soon!")
			return nil
		}
	}

	printUsage()
	return nil
}

func printUsage() {
	fmt.Printf(`fabric-lite v%s - A lightweight AI augmentation framework

Usage:
  fabric-lite [flags]
  fabric-lite --pattern <name> < input.txt
  echo "text" | fabric-lite --pattern <name>

Flags:
  -p, --pattern <name>   Pattern to use (e.g., summarize)
  -m, --model <model>    Model to use (default: gpt-4o-mini)
  -l, --list             List available patterns
  -v, --version          Show version
  -h, --help             Show this help

Examples:
  # Summarize a file
  fabric-lite --pattern summarize < article.txt

  # Summarize piped input
  cat README.md | fabric-lite --pattern summarize

  # Use a specific model
  echo "Hello world" | fabric-lite -p summarize -m gpt-4o

For more information, visit: https://github.com/rice0649/fabric-lite
`, Version)
}
