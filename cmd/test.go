package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/crisecheverria/codequest/internal/challenge"
	"github.com/crisecheverria/codequest/internal/native"
	"github.com/spf13/cobra"
)

type ChallengeMetadata struct {
	Slug         string `json:"slug"`
	Language     string `json:"language"`
	FunctionName string `json:"functionName"`
	SolutionFile string `json:"solutionFile"`
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test your solution against the challenge test cases",
	Long: `Execute your solution locally using native language runtimes and validate 
against the challenge test cases. Requires Go and/or Node.js to be installed 
depending on the challenge language.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Look for challenge metadata
		metadataPath := ".challenge.json"
		if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
			return fmt.Errorf("not in a challenge directory. Run 'codequest fetch <challenge>' first")
		}

		// Load metadata
		metadata, err := loadChallengeMetadata(metadataPath)
		if err != nil {
			return fmt.Errorf("failed to load challenge metadata: %w", err)
		}

		// Load the challenge from data
		challenges, err := challenge.LoadChallenges()
		if err != nil {
			return fmt.Errorf("failed to load challenges: %w", err)
		}

		ch, found := challenge.FindBySlug(challenges, metadata.Slug)
		if !found {
			return fmt.Errorf("challenge '%s' not found", metadata.Slug)
		}

		// Read solution file
		solutionCode, err := os.ReadFile(metadata.SolutionFile)
		if err != nil {
			return fmt.Errorf("failed to read solution file '%s': %w", metadata.SolutionFile, err)
		}

		// Create native executor
		executor, err := native.NewExecutor()
		if err != nil {
			return fmt.Errorf("failed to create native executor: %w", err)
		}
		defer executor.Close()

		// Test the solution
		fmt.Printf("Testing solution for '%s'...\n\n", ch.Title)

		success := true
		for i, testCase := range ch.TestCases {
			fmt.Printf("Test %d: %s\n", i+1, testCase.Description)

			// Create test code that calls the function with test inputs
			testCode := generateTestCode(ch, string(solutionCode), testCase)

			// Use longer timeout for Go due to compilation overhead
			timeout := ch.TimeLimit
			if ch.Language == "go" {
				timeout = 15000 // 15 seconds for Go compilation + execution
			}
			result, err := executor.ExecuteCode(ch.Language, testCode, timeout)
			if err != nil {
				fmt.Printf("  ❌ Execution error: %v\n\n", err)
				success = false
				continue
			}

			if result.Success {
				fmt.Printf("  ✅ Passed (%.2fms)\n", float64(result.Duration.Nanoseconds())/1e6)
			} else {
				fmt.Printf("  ❌ Failed\n")
				if result.Error != "" {
					fmt.Printf("     Error: %s\n", result.Error)
				}
				if result.Output != "" {
					fmt.Printf("     Output: %s\n", result.Output)
				}
				success = false
			}
			fmt.Println()
		}

		if success {
			fmt.Println("🎉 All tests passed! Run 'codequest submit' to submit your solution.")
		} else {
			fmt.Println("❌ Some tests failed. Please fix your solution and try again.")
		}

		return nil
	},
}

func loadChallengeMetadata(path string) (*ChallengeMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata ChallengeMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func generateTestCode(ch challenge.Challenge, solutionCode string, testCase challenge.TestCase) string {
	switch ch.Language {
	case "typescript":
		return generateTypeScriptTestCode(ch, solutionCode, testCase)
	case "javascript":
		return generateJavaScriptTestCode(ch, solutionCode, testCase)
	case "php":
		return generatePHPTestCode(ch, solutionCode, testCase)
	case "go":
		return generateGoTestCode(ch, solutionCode, testCase)
	default:
		return solutionCode
	}
}

func generateTypeScriptTestCode(ch challenge.Challenge, solutionCode string, testCase challenge.TestCase) string {
	// Convert TypeScript function to JavaScript by removing type annotations
	jsCode := strings.ReplaceAll(solutionCode, ": number", "")
	jsCode = strings.ReplaceAll(jsCode, ": string", "")
	jsCode = strings.ReplaceAll(jsCode, ": boolean", "")

	// Convert input array to function arguments
	args := make([]string, len(testCase.Input))
	for i, input := range testCase.Input {
		switch v := input.(type) {
		case string:
			args[i] = fmt.Sprintf(`"%s"`, v)
		case float64:
			if v == float64(int(v)) {
				args[i] = fmt.Sprintf("%d", int(v))
			} else {
				args[i] = fmt.Sprintf("%f", v)
			}
		default:
			args[i] = fmt.Sprintf("%v", v)
		}
	}

	expectedStr := fmt.Sprintf("%v", testCase.Expected)
	if _, ok := testCase.Expected.(string); ok {
		expectedStr = fmt.Sprintf(`"%v"`, testCase.Expected)
	}

	argsStr := strings.Join(args, ", ")

	return fmt.Sprintf(`%s

const result = %s(%s);
const expected = %s;

if (JSON.stringify(result) === JSON.stringify(expected)) {
  console.log("Test passed");
  process.exit(0);
} else {
  console.log("Expected:", expected, "Got:", result);
  process.exit(1);
}`, jsCode, ch.FunctionName, argsStr, expectedStr)
}

func generateJavaScriptTestCode(ch challenge.Challenge, solutionCode string, testCase challenge.TestCase) string {
	return generateTypeScriptTestCode(ch, solutionCode, testCase) // Same logic for JS
}

func generatePHPTestCode(ch challenge.Challenge, solutionCode string, testCase challenge.TestCase) string {
	args := make([]string, len(testCase.Input))
	for i, input := range testCase.Input {
		switch v := input.(type) {
		case string:
			args[i] = fmt.Sprintf(`"%s"`, v)
		case float64:
			args[i] = fmt.Sprintf("%v", v)
		default:
			args[i] = fmt.Sprintf("%v", v)
		}
	}

	argsStr := strings.Join(args, ", ")

	return fmt.Sprintf(`<?php
%s

$result = %s(%s);
$expected = %v;

if ($result === $expected) {
    echo "Test passed\n";
    exit(0);
} else {
    echo "Expected: " . print_r($expected, true) . " Got: " . print_r($result, true) . "\n";
    exit(1);
}
?>`, solutionCode, ch.FunctionName, argsStr, testCase.Expected)
}

func generateGoTestCode(ch challenge.Challenge, solutionCode string, testCase challenge.TestCase) string {
	args := make([]string, len(testCase.Input))
	for i, input := range testCase.Input {
		switch v := input.(type) {
		case string:
			args[i] = fmt.Sprintf(`"%s"`, v)
		case float64:
			if v == float64(int(v)) {
				args[i] = fmt.Sprintf("%d", int(v))
			} else {
				args[i] = fmt.Sprintf("%f", v)
			}
		default:
			args[i] = fmt.Sprintf("%v", v)
		}
	}

	// Join args properly for function call
	argsStr := ""
	if len(args) > 0 {
		argsStr = strings.Join(args, ", ")
	}

	// Clean user code by removing package declaration, imports, and main function
	cleanedCode := cleanGoUserCode(solutionCode)

	// Handle multiple return values by checking if expected is an array
	var testLogic string
	if expectedSlice, ok := testCase.Expected.([]interface{}); ok {
		// Multiple return values - create variables to capture each return value
		if len(expectedSlice) == 2 {
			testLogic = fmt.Sprintf(`	result1, result2 := %s(%s)
	expected1, expected2 := %v, %v
	
	if result1 == expected1 && result2 == expected2 {
		fmt.Println("Test passed")
		os.Exit(0)
	} else {
		fmt.Printf("Expected: [%%v, %%v] Got: [%%v, %%v]\n", expected1, expected2, result1, result2)
		os.Exit(1)
	}`, ch.FunctionName, argsStr, expectedSlice[0], expectedSlice[1])
		} else {
			// Fallback for other multiple return value counts
			var expectedStr string
			if _, ok := testCase.Expected.(string); ok {
				expectedStr = fmt.Sprintf(`"%s"`, testCase.Expected)
			} else {
				expectedStr = fmt.Sprintf("%v", testCase.Expected)
			}
			
			testLogic = fmt.Sprintf(`	result := %s(%s)
	expected := %s
	
	if fmt.Sprintf("%%v", result) == fmt.Sprintf("%%v", expected) {
		fmt.Println("Test passed")
		os.Exit(0)
	} else {
		fmt.Printf("Expected: %%v Got: %%v\n", expected, result)
		os.Exit(1)
	}`, ch.FunctionName, argsStr, expectedStr)
		}
	} else {
		// Single return value
		var expectedStr string
		if _, ok := testCase.Expected.(string); ok {
			expectedStr = fmt.Sprintf(`"%s"`, testCase.Expected)
		} else {
			expectedStr = fmt.Sprintf("%v", testCase.Expected)
		}
		
		testLogic = fmt.Sprintf(`	result := %s(%s)
	expected := %s
	
	if result == expected {
		fmt.Println("Test passed")
		os.Exit(0)
	} else {
		fmt.Printf("Expected: %%v Got: %%v\n", expected, result)
		os.Exit(1)
	}`, ch.FunctionName, argsStr, expectedStr)
	}

	return fmt.Sprintf(`package main

import (
	"fmt"
	"os"
)

%s

func main() {
%s
}`, cleanedCode, testLogic)
}

func cleanGoUserCode(code string) string {
	lines := strings.Split(code, "\n")
	var cleanedLines []string

	inImportBlock := false
	inMainFunction := false
	braceCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip package declarations
		if strings.HasPrefix(trimmed, "package ") {
			continue
		}

		// Skip imports (we provide our own imports)
		if strings.HasPrefix(trimmed, "import ") && !strings.Contains(trimmed, "(") {
			continue
		}

		if strings.HasPrefix(trimmed, "import (") {
			inImportBlock = true
			continue
		}

		if inImportBlock {
			if trimmed == ")" {
				inImportBlock = false
			}
			continue
		}

		// Skip user's main function
		if strings.HasPrefix(trimmed, "func main(") {
			inMainFunction = true
			braceCount = 0
			// Count opening brace on the same line
			for _, r := range line {
				if r == '{' {
					braceCount++
				}
			}
			continue
		}

		if inMainFunction {
			// Count braces to know when main function ends
			for _, r := range line {
				if r == '{' {
					braceCount++
				} else if r == '}' {
					braceCount--
				}
			}

			if braceCount <= 0 {
				inMainFunction = false
			}
			continue
		}

		cleanedLines = append(cleanedLines, line)
	}

	// Remove empty lines at the beginning and end
	for len(cleanedLines) > 0 && strings.TrimSpace(cleanedLines[0]) == "" {
		cleanedLines = cleanedLines[1:]
	}
	for len(cleanedLines) > 0 && strings.TrimSpace(cleanedLines[len(cleanedLines)-1]) == "" {
		cleanedLines = cleanedLines[:len(cleanedLines)-1]
	}

	return strings.Join(cleanedLines, "\n")
}

func init() {
	rootCmd.AddCommand(testCmd)
}
