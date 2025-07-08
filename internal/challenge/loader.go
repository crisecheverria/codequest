package challenge

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func LoadChallenges() ([]Challenge, error) {
	// First try to use embedded data
	if len(embeddedChallenges) > 0 {
		var challenges []Challenge
		if err := json.Unmarshal(embeddedChallenges, &challenges); err != nil {
			return nil, fmt.Errorf("failed to parse embedded challenges JSON: %w", err)
		}
		return challenges, nil
	}

	// Fallback to external file for development
	possiblePaths := []string{
		"data/challenges.json",                  // From binary location (standalone)
		"../data/challenges.json",              // From challenge directory
		"../../data/challenges.json",           // From nested challenge directory
	}

	var dataPath string
	var err error

	for _, path := range possiblePaths {
		if _, err = os.Stat(path); err == nil {
			dataPath = path
			break
		}
	}

	if dataPath == "" {
		return nil, fmt.Errorf("could not find challenges.json file. Tried paths: %v", possiblePaths)
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read challenges file: %w", err)
	}

	var challenges []Challenge
	if err := json.Unmarshal(data, &challenges); err != nil {
		return nil, fmt.Errorf("failed to parse challenges JSON: %w", err)
	}

	return challenges, nil
}

func FindBySlug(challenges []Challenge, slug string) (Challenge, bool) {
	for _, ch := range challenges {
		if ch.Slug == slug {
			return ch, true
		}
	}
	return Challenge{}, false
}

func FilterChallenges(challenges []Challenge, language, difficulty string) []Challenge {
	var filtered []Challenge

	for _, ch := range challenges {
		if language != "" && !strings.EqualFold(ch.Language, language) {
			continue
		}
		if difficulty != "" && !strings.EqualFold(ch.Difficulty, difficulty) {
			continue
		}
		filtered = append(filtered, ch)
	}

	return filtered
}

func FilterByLanguage(challenges []Challenge, language string) []Challenge {
	var filtered []Challenge

	for _, ch := range challenges {
		if strings.EqualFold(ch.Language, language) {
			filtered = append(filtered, ch)
		}
	}

	return filtered
}

func DisplayChallengeList(challenges []Challenge) {
	if len(challenges) == 0 {
		fmt.Println("No challenges found matching the criteria.")
		return
	}

	fmt.Printf("%-40s %-12s %-10s %s\n", "TITLE", "LANGUAGE", "DIFFICULTY", "SLUG")
	fmt.Println(strings.Repeat("-", 80))

	for _, ch := range challenges {
		fmt.Printf("%-40s %-12s %-10s %s\n",
			truncateString(ch.Title, 40),
			ch.Language,
			ch.Difficulty,
			ch.Slug,
		)
	}

	fmt.Printf("\nTotal: %d challenges\n", len(challenges))
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
