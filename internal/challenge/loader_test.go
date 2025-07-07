package challenge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadChallengesFromFile(t *testing.T) {
	// Create a temporary challenge file for testing
	tempDir := t.TempDir()
	challengesFile := filepath.Join(tempDir, "challenges.json")

	testChallenges := []Challenge{
		{
			Title:          "Test Challenge",
			Description:    "A test challenge",
			Difficulty:     "easy",
			Language:       "go",
			Slug:           "test-challenge",
			FunctionName:   "testFunc",
			ParameterTypes: []string{"int"},
			ReturnType:     "string",
			Template:       "package main\n\nfunc testFunc(n int) string {\n\treturn \"\"\n}",
			TestCases: []TestCase{
				{
					Input:       []interface{}{42},
					Expected:    "42",
					Description: "should convert int to string",
				},
			},
			ConceptTags: []string{"variables"},
			TimeLimit:   5000,
			MemoryLimit: 128,
		},
	}

	// Write test data to file
	data, err := json.Marshal(testChallenges)
	if err != nil {
		t.Fatalf("Failed to marshal test challenges: %v", err)
	}

	if err := os.WriteFile(challengesFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test challenges file: %v", err)
	}

	// Test loading challenges from file directly
	fileData, err := os.ReadFile(challengesFile)
	if err != nil {
		t.Fatalf("Failed to read test challenges file: %v", err)
	}

	var challenges []Challenge
	if err := json.Unmarshal(fileData, &challenges); err != nil {
		t.Fatalf("Failed to unmarshal challenges: %v", err)
	}

	if len(challenges) != 1 {
		t.Errorf("Expected 1 challenge, got %d", len(challenges))
	}

	challenge := challenges[0]
	if challenge.Title != "Test Challenge" {
		t.Errorf("Expected title 'Test Challenge', got '%s'", challenge.Title)
	}

	if challenge.Slug != "test-challenge" {
		t.Errorf("Expected slug 'test-challenge', got '%s'", challenge.Slug)
	}
}

func TestFindBySlug(t *testing.T) {
	challenges := []Challenge{
		{
			Slug:  "first-challenge",
			Title: "First Challenge",
		},
		{
			Slug:  "second-challenge",
			Title: "Second Challenge",
		},
	}

	// Test finding existing challenge
	challenge, found := FindBySlug(challenges, "first-challenge")
	if !found {
		t.Error("Expected to find challenge with slug 'first-challenge'")
	}
	if challenge.Title != "First Challenge" {
		t.Errorf("Expected title 'First Challenge', got '%s'", challenge.Title)
	}

	// Test finding non-existing challenge
	_, found = FindBySlug(challenges, "non-existent")
	if found {
		t.Error("Expected not to find challenge with slug 'non-existent'")
	}
}

func TestFilterByLanguage(t *testing.T) {
	challenges := []Challenge{
		{
			Slug:     "go-challenge",
			Language: "go",
		},
		{
			Slug:     "ts-challenge",
			Language: "typescript",
		},
		{
			Slug:     "another-go",
			Language: "go",
		},
	}

	goChallenges := FilterByLanguage(challenges, "go")
	if len(goChallenges) != 2 {
		t.Errorf("Expected 2 Go challenges, got %d", len(goChallenges))
	}

	tsChallenges := FilterByLanguage(challenges, "typescript")
	if len(tsChallenges) != 1 {
		t.Errorf("Expected 1 TypeScript challenge, got %d", len(tsChallenges))
	}

	phpChallenges := FilterByLanguage(challenges, "php")
	if len(phpChallenges) != 0 {
		t.Errorf("Expected 0 PHP challenges, got %d", len(phpChallenges))
	}
}
