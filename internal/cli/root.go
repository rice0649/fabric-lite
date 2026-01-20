package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/rice0649/fabric-lite/internal/core"
	"github.com/rice0649/fabric-lite/internal/executor"
	"github.com/rice0649/fabric-lite/internal/tools"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "fabric-lite",
		Short:   "Lightweight AI augmentation framework",
		Long:    `Fabric-Lite is a personal CLI tool that runs AI prompts (patterns) against text input.`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig()
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringP("pattern", "p", "", "Pattern to use (e.g., summarize)")
	rootCmd.PersistentFlags().StringP("model", "m", "gpt-4o-mini", "Model to use")
	rootCmd.PersistentFlags().StringP("provider", "v", "", "AI provider to use (openai, anthropic, ollama)")
	rootCmd.PersistentFlags().BoolP("stream", "s", false, "Enable streaming responses")
	rootCmd.PersistentFlags().BoolP("verbose", "V", false, "Verbose output")

	viper.BindPFlag("pattern", rootCmd.PersistentFlags().Lookup("pattern"))
	viper.BindPFlag("model", rootCmd.PersistentFlags().Lookup("model"))
	viper.BindPFlag("provider", rootCmd.PersistentFlags().Lookup("provider"))
	viper.BindPFlag("stream", rootCmd.PersistentFlags().Lookup("stream"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Add subcommands
	rootCmd.AddCommand(newRunCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newVersionCmd(version))

	return rootCmd
}

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/fabric-lite")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("FABRIC_LITE")

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		// Config file not found is OK, we'll use defaults
	}

	return nil
}

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [input_file|tool_name]",
		Short: "Execute a pattern against input or run a tool",
		Long:  `Execute a fabric-lite pattern against input text or run a specific tool.`,
		Args:  cobra.MaximumNArgs(2),
		RunE:  runCommand,
	}

	cmd.Flags().String("pattern", "", "Pattern name (required)")
	cmd.Flags().StringP("prompt", "P", "", "Prompt for tool execution")
	cmd.Flags().String("model", "", "Model to use")
	cmd.Flags().String("provider", "", "Provider to use")
	cmd.Flags().Bool("stream", false, "Stream response")
	cmd.Flags().Bool("save-session", false, "Save conversation to session")

	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	// Check if first arg is a tool name (for tool execution mode)
	if len(args) > 0 {
		if tool, err := tools.GetTool(args[0]); err == nil {
			return executeTool(cmd, args, tool)
		}
	}

	// Otherwise, treat as pattern execution
	return executePattern(cmd, args)
}

func executeTool(cmd *cobra.Command, args []string, tool tools.Tool) error {
	// Parse tool-specific flags
	prompt := cmd.Flag("prompt").Value.String()
	if prompt == "" {
		// For tools, prompt might be passed directly without --prompt flag
		// Look for -P flag in remaining args
		for i := 1; i < len(args); i++ {
			if args[i] == "-P" && i+1 < len(args) {
				prompt = args[i+1]
				break
			}
		}
	}

	if prompt == "" {
		return fmt.Errorf("prompt is required for tool execution (use -p \"prompt\")")
	}

	// Create execution context
	executionContext := tools.ExecutionContext{
		Prompt:  prompt,
		Args:    args[1:], // Skip tool name
		Env:     make(map[string]string),
		WorkDir: ".",
	}

	// Execute the tool
	result, err := tool.Execute(executionContext)
	if err != nil {
		return fmt.Errorf("tool execution failed: %w", err)
	}

	// Display the result
	fmt.Print(result.Output)
	return nil
}

func executePattern(cmd *cobra.Command, args []string) error {
	patternName := cmd.Flag("pattern").Value.String()
	if patternName == "" {
		patternName = viper.GetString("pattern")
	}

	if patternName == "" {
		return fmt.Errorf("pattern is required (use --pattern or set in config)")
	}

	// Get input
	var input string
	if len(args) > 0 {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
		input = string(data)
	} else {
		// Read from stdin
		stat, _ := os.Stdin.Stat()
		if stat != nil && (stat.Mode()&os.ModeCharDevice) != 0 {
			return fmt.Errorf("no input file provided and stdin is not available")
		}

		var buf strings.Builder
		for {
			var chunk [1024]byte
			n, err := os.Stdin.Read(chunk[:])
			if err != nil && err.Error() != "EOF" {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			if n == 0 {
				break
			}
			buf.Write(chunk[:n])
		}
		input = buf.String()
	}

	// Initialize config manager
	configManager := core.NewConfigManager("")
	config, err := configManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize provider manager
	providerManager := core.NewProviderManager(config)
	if err := providerManager.InitializeAll(); err != nil {
		return fmt.Errorf("failed to initialize providers: %w", err)
	}

	// Determine provider name
	providerName := cmd.Flag("provider").Value.String()
	if providerName == "" {
		providerName = viper.GetString("provider")
	}
	if providerName == "" {
		providerName, _ = configManager.GetDefaultProvider()
	}

	// Initialize executor with provider manager
	executor := executor.NewPatternExecutor()

	// Load specific provider into executor
	provider, err := providerManager.Get(providerName)
	if err != nil {
		return fmt.Errorf("failed to get provider %s: %w", providerName, err)
	}
	executor.LoadProviderDirect(providerName, provider)

	response, err := executor.Execute(cmd.Context(), patternName, input, providerName)
	if err != nil {
		return fmt.Errorf("failed to execute pattern: %w", err)
	}

	fmt.Print(response.Content)
	return nil
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available patterns",
		Long:  `List all available patterns in the patterns directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			patterns, err := executor.NewPatternExecutor().ListPatterns()
			if err != nil {
				return fmt.Errorf("failed to list patterns: %w", err)
			}

			fmt.Println("Available patterns:")
			for _, pattern := range patterns {
				fmt.Printf("  - %s (%s)\n", pattern.Name, pattern.Description)
			}

			return nil
		},
	}
}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show configuration",
		Long:  `Display current fabric-lite configuration.`,
		RunE:  showConfig,
	}
	return cmd
}

func showConfig(cmd *cobra.Command, args []string) error {
	fmt.Println("Fabric-Lite Configuration:")
	fmt.Printf("  Config File: %s\n", viper.ConfigFileUsed())
	fmt.Printf("  Provider: %s\n", viper.GetString("provider"))
	fmt.Printf("  Model: %s\n", viper.GetString("model"))
	fmt.Printf("  Patterns Dir: %s\n", executor.NewPatternExecutor().GetPatternsDir())

	// Show loaded providers
	configManager := core.NewConfigManager("")
	config, err := configManager.Load()
	if err == nil {
		fmt.Println("  Available Providers:")
		for _, provider := range config.Providers {
			fmt.Printf("    - %s (%s)\n", provider.Name, provider.Type)
		}
	}

	return nil
}

func newVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Display version information for fabric-lite.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("fabric-lite version %s\n", version)
		},
	}
}
