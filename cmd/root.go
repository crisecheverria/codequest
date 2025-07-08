package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version will be set during build time via ldflags
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "codequest",
	Version: Version,
	Short:   "CodeQuest CLI - Interactive coding challenges in your terminal",
	Long: `CodeQuest CLI allows you to fetch, solve, and submit coding challenges 
directly from your terminal with native language execution.

Examples:
  codequest list                    # List available challenges
  codequest fetch two-sum           # Fetch a specific challenge
  codequest test                    # Test your solution locally`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to CodeQuest CLI!")
		fmt.Println("Use 'codequest --help' to see available commands.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.codequest.yaml)")
}
