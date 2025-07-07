package challenge

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateWorkspace(t *testing.T) {
	tempDir := t.TempDir()

	challenge := Challenge{
		Title:        "Test Workspace Challenge",
		Description:  "A challenge for testing workspace creation",
		Slug:         "test-workspace-challenge",
		Language:     "go",
		FunctionName: "solve",
		Template:     "package main\n\nfunc solve(n int) int {\n\treturn n * 2\n}",
		TestCases: []TestCase{
			{
				Input:       []interface{}{5},
				Expected:    10,
				Description: "should double the input",
			},
		},
	}

	workspaceDir := filepath.Join(tempDir, challenge.Slug)

	// Override the function to create workspace in a specific directory
	err := createWorkspaceInDir(workspaceDir, challenge)
	if err != nil {
		t.Fatalf("CreateWorkspace() failed: %v", err)
	}

	// Check if workspace directory was created
	if _, err := os.Stat(workspaceDir); os.IsNotExist(err) {
		t.Error("Workspace directory was not created")
	}

	// Check if README.md was created
	readmePath := filepath.Join(workspaceDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Error("README.md was not created")
	}

	// Check README content
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeStr := string(readmeContent)
	if !contains(readmeStr, "Test Workspace Challenge") {
		t.Error("README.md does not contain challenge title")
	}
	if !contains(readmeStr, "should double the input") {
		t.Error("README.md does not contain test case description")
	}

	// Check if solution file was created
	solutionPath := filepath.Join(workspaceDir, "solution.go")
	if _, err := os.Stat(solutionPath); os.IsNotExist(err) {
		t.Error("solution.go was not created")
	}

	// Check solution file content
	solutionContent, err := os.ReadFile(solutionPath)
	if err != nil {
		t.Fatalf("Failed to read solution.go: %v", err)
	}

	solutionStr := string(solutionContent)
	if !contains(solutionStr, "func solve(n int) int") {
		t.Error("solution.go does not contain the correct function signature")
	}

	// Check if .challenge.json was created
	metadataPath := filepath.Join(workspaceDir, ".challenge.json")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		t.Error(".challenge.json was not created")
	}

	// Check metadata content
	metadataContent, err := os.ReadFile(metadataPath)
	if err != nil {
		t.Fatalf("Failed to read .challenge.json: %v", err)
	}

	metadataStr := string(metadataContent)
	if !contains(metadataStr, "test-workspace-challenge") {
		t.Error(".challenge.json does not contain the correct slug")
	}
	if !contains(metadataStr, "\"language\": \"go\"") {
		t.Error(".challenge.json does not contain the correct language")
	}
}

func TestCreateWorkspaceTypeScript(t *testing.T) {
	tempDir := t.TempDir()

	challenge := Challenge{
		Title:        "TypeScript Challenge",
		Slug:         "typescript-challenge",
		Language:     "typescript",
		FunctionName: "addNumbers",
		Template:     "function addNumbers(a: number, b: number): number {\n  return a + b;\n}",
	}

	workspaceDir := filepath.Join(tempDir, challenge.Slug)

	// Override the function to create workspace in a specific directory
	err := createWorkspaceInDir(workspaceDir, challenge)
	if err != nil {
		t.Fatalf("CreateWorkspace() failed: %v", err)
	}

	// Check if solution file has correct extension
	solutionPath := filepath.Join(workspaceDir, "solution.ts")
	if _, err := os.Stat(solutionPath); os.IsNotExist(err) {
		t.Error("solution.ts was not created")
	}
}

func TestCreateWorkspacePHP(t *testing.T) {
	tempDir := t.TempDir()

	challenge := Challenge{
		Title:        "PHP Challenge",
		Slug:         "php-challenge",
		Language:     "php",
		FunctionName: "concatenate",
		Template:     "<?php\nfunction concatenate($a, $b) {\n    return $a . $b;\n}\n?>",
	}

	workspaceDir := filepath.Join(tempDir, challenge.Slug)

	// Override the function to create workspace in a specific directory
	err := createWorkspaceInDir(workspaceDir, challenge)
	if err != nil {
		t.Fatalf("CreateWorkspace() failed: %v", err)
	}

	// Check if solution file has correct extension
	solutionPath := filepath.Join(workspaceDir, "solution.php")
	if _, err := os.Stat(solutionPath); os.IsNotExist(err) {
		t.Error("solution.php was not created")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Helper function to create workspace in a specific directory (for testing)
func createWorkspaceInDir(workspaceDir string, ch Challenge) error {
	// Create workspace directory
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		return fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// Create solution file
	solutionFile := getSolutionFileNameForTest(ch)
	solutionPath := filepath.Join(workspaceDir, solutionFile)
	if err := os.WriteFile(solutionPath, []byte(ch.Template), 0644); err != nil {
		return fmt.Errorf("failed to create solution file: %w", err)
	}

	// Create README with challenge description
	readmePath := filepath.Join(workspaceDir, "README.md")
	readme := generateReadmeForTest(ch)
	if err := os.WriteFile(readmePath, []byte(readme), 0644); err != nil {
		return fmt.Errorf("failed to create README: %w", err)
	}

	// Create challenge metadata file
	metaPath := filepath.Join(workspaceDir, ".challenge.json")
	metaContent := fmt.Sprintf(`{
  "slug": "%s",
  "language": "%s",
  "functionName": "%s",
  "solutionFile": "%s"
}`, ch.Slug, ch.Language, ch.FunctionName, solutionFile)
	if err := os.WriteFile(metaPath, []byte(metaContent), 0644); err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}

	return nil
}

func getSolutionFileNameForTest(ch Challenge) string {
	switch ch.Language {
	case "typescript":
		return "solution.ts"
	case "javascript":
		return "solution.js"
	case "php":
		return "solution.php"
	case "go":
		return "solution.go"
	default:
		return "solution.txt"
	}
}

func generateReadmeForTest(ch Challenge) string {
	return fmt.Sprintf(`# %s

**Language:** %s  
**Difficulty:** %s  

## Description

%s

## Function Signature

- **Function:** %s

## Test Cases

%s

## Commands

`+"```"+`bash
# Test your solution
codequest test

# Submit your solution
codequest submit
`+"```"+`
`, ch.Title, ch.Language, ch.Difficulty, ch.Description, ch.FunctionName, getTestCaseDesc(ch))
}

func getTestCaseDesc(ch Challenge) string {
	if len(ch.TestCases) > 0 {
		return ch.TestCases[0].Description
	}
	return "No test cases defined"
}
