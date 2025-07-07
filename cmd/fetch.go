package cmd

import (
	"fmt"

	"github.com/crisecheverria/codequest/internal/challenge"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [challenge-slug]",
	Short: "Fetch a specific challenge and create local files",
	Long: `Fetch a coding challenge by its slug and create the necessary template files 
for local development. This creates a working directory with the challenge template,
test cases, and instructions.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		slug := args[0]

		challenges, err := challenge.LoadChallenges()
		if err != nil {
			return fmt.Errorf("failed to load challenges: %w", err)
		}

		ch, found := challenge.FindBySlug(challenges, slug)
		if !found {
			return fmt.Errorf("challenge '%s' not found", slug)
		}

		workDir, err := challenge.CreateWorkspace(ch)
		if err != nil {
			return fmt.Errorf("failed to create workspace: %w", err)
		}

		fmt.Printf("Challenge '%s' fetched successfully!\n", ch.Title)
		fmt.Printf("Working directory: %s\n", workDir)
		fmt.Printf("Edit the solution file and run 'codequest test' to validate.\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
