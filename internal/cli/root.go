package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "forge",
		Short: "AI Project Forge - Orchestrate AI coding assistants",
		Long: `AI Project Forge orchestrates multiple AI coding assistants through
structured development phases to help build any software project.

Supported tools:
  - Gemini CLI: Research and discovery (free tier, 1M context)
  - OpenCode: Planning and design (read-only exploration)
  - Codex CLI: Implementation (advanced reasoning)
  - fabric-lite: Pattern-based document generation`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig()
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default .forge/config.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Add subcommands
	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newRunCmd())
	rootCmd.AddCommand(newPhaseCmd())
	rootCmd.AddCommand(newStatusCmd())
	rootCmd.AddCommand(newHistoryCmd())
	rootCmd.AddCommand(newSessionCmd())

	return rootCmd
}

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".forge")
	viper.AddConfigPath("$HOME/.config/forge")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("FORGE")

	// Read config file if it exists (don't error if not found)
	_ = viper.ReadInConfig()

	return nil
}
