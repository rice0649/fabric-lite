package main

import (
	"fmt"
	"github.com/rice0649/fabric-lite/internal/cli"
	"os"
)

// Version is set at build time via ldflags
var Version = "0.1.0"

func main() {
	rootCmd := cli.NewRootCmd(Version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
