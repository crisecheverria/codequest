package challenge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateWorkspace(ch Challenge) (string, error) {
	// Create workspace directory
	workDir := fmt.Sprintf("challenge-%s", ch.Slug)
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// Create solution file
	solutionFile := getSolutionFileName(ch)
	solutionPath := filepath.Join(workDir, solutionFile)
	if err := os.WriteFile(solutionPath, []byte(ch.Template), 0644); err != nil {
		return "", fmt.Errorf("failed to create solution file: %w", err)
	}

	// Create README with challenge description
	readmePath := filepath.Join(workDir, "README.md")
	readme := generateReadme(ch)
	if err := os.WriteFile(readmePath, []byte(readme), 0644); err != nil {
		return "", fmt.Errorf("failed to create README: %w", err)
	}

	// Create challenge metadata file
	metaPath := filepath.Join(workDir, ".challenge.json")
	metaContent := fmt.Sprintf(`{
  "slug": "%s",
  "language": "%s",
  "functionName": "%s",
  "solutionFile": "%s"
}`, ch.Slug, ch.Language, ch.FunctionName, solutionFile)
	if err := os.WriteFile(metaPath, []byte(metaContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create metadata file: %w", err)
	}

	return workDir, nil
}

func getSolutionFileName(ch Challenge) string {
	switch ch.Language {
	case "typescript":
		return "solution.ts"
	case "javascript":
		return "solution.js"
	case "go":
		return "solution.go"
	case "python":
		return "solution.py"
	default:
		return "solution.txt"
	}
}

func generateReadme(ch Challenge) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("# %s\n\n", ch.Title))
	builder.WriteString(fmt.Sprintf("**Language:** %s  \n", ch.Language))
	builder.WriteString(fmt.Sprintf("**Difficulty:** %s  \n", ch.Difficulty))
	builder.WriteString(fmt.Sprintf("**Concepts:** %s\n\n", strings.Join(ch.ConceptTags, ", ")))

	if ch.Description != "" {
		builder.WriteString(fmt.Sprintf("## Description\n\n%s\n\n", ch.Description))
	}

	builder.WriteString("## Function Signature\n\n")
	builder.WriteString(fmt.Sprintf("- **Function:** `%s`\n", ch.FunctionName))
	builder.WriteString(fmt.Sprintf("- **Parameters:** `%s`\n", strings.Join(ch.ParameterTypes, ", ")))
	builder.WriteString(fmt.Sprintf("- **Return Type:** `%s`\n\n", ch.ReturnType))

	builder.WriteString("## Test Cases\n\n")
	for i, testCase := range ch.TestCases {
		builder.WriteString(fmt.Sprintf("**Test %d:** %s\n", i+1, testCase.Description))
		builder.WriteString(fmt.Sprintf("- Input: `%v`\n", testCase.Input))
		builder.WriteString(fmt.Sprintf("- Expected: `%v`\n\n", testCase.Expected))
	}

	builder.WriteString("## Commands\n\n")
	builder.WriteString("```bash\n")
	builder.WriteString("# Test your solution\n")
	builder.WriteString("codequest test\n\n")
	builder.WriteString("# Submit your solution\n")
	builder.WriteString("codequest submit\n")
	builder.WriteString("```\n")

	return builder.String()
}
