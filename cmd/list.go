package cmd

import (
	"fmt"

	"github.com/crisecheverria/codequest/internal/challenge"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available coding challenges",
	Long:  `Display all available coding challenges with their difficulty and language.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		language, _ := cmd.Flags().GetString("language")
		difficulty, _ := cmd.Flags().GetString("difficulty")

		challenges, err := challenge.LoadChallenges()
		if err != nil {
			return fmt.Errorf("failed to load challenges: %w", err)
		}

		filtered := challenge.FilterChallenges(challenges, language, difficulty)
		challenge.DisplayChallengeList(filtered)

		return nil
	},
}

func init() {
	listCmd.Flags().StringP("language", "l", "", "Filter by language (typescript, php, go)")
	listCmd.Flags().StringP("difficulty", "d", "", "Filter by difficulty (easy, medium, hard)")
	rootCmd.AddCommand(listCmd)
}
