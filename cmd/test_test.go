package cmd

import (
	"strings"
	"testing"

	"github.com/crisecheverria/codequest/internal/challenge"
)

func TestGenerateGoTestCode(t *testing.T) {
	ch := challenge.Challenge{
		Language:     "go",
		FunctionName: "add",
	}

	userCode := `package main

import "fmt"

func add(a, b int) int {
    return a + b
}

func main() {
    fmt.Println("Hello")
}`

	testCase := challenge.TestCase{
		Input:       []interface{}{float64(2), float64(3)},
		Expected:    "5",
		Description: "should add two numbers",
	}

	result := generateGoTestCode(ch, userCode, testCase)

	// Check that the generated code has the correct structure
	if !strings.Contains(result, "package main") {
		t.Error("Generated code should contain 'package main'")
	}

	if !strings.Contains(result, "func add(a, b int) int") {
		t.Error("Generated code should contain the user's function")
	}

	if !strings.Contains(result, "result := add(2, 3)") {
		t.Error("Generated code should call the function with test inputs")
	}

	if !strings.Contains(result, "expected := \"5\"") {
		t.Error("Generated code should set the expected value")
	}

	if strings.Contains(result, "func main() {") && strings.Count(result, "func main()") > 1 {
		t.Error("Generated code should not contain multiple main functions")
	}

	// Check that user's imports are removed
	if strings.Contains(result, "import \"fmt\"") {
		t.Error("Generated code should not contain user's imports")
	}

	// Check that user's original main function is removed
	if strings.Contains(result, "fmt.Println(\"Hello\")") {
		t.Error("Generated code should not contain user's main function content")
	}
}

func TestGenerateTypeScriptTestCode(t *testing.T) {
	ch := challenge.Challenge{
		Language:     "typescript",
		FunctionName: "multiply",
	}

	userCode := `function multiply(a: number, b: number): number {
    return a * b;
}`

	testCase := challenge.TestCase{
		Input:       []interface{}{float64(4), float64(5)},
		Expected:    float64(20),
		Description: "should multiply two numbers",
	}

	result := generateTypeScriptTestCode(ch, userCode, testCase)

	// Check that type annotations are removed
	if strings.Contains(result, ": number") {
		t.Error("Generated code should not contain TypeScript type annotations")
	}

	if !strings.Contains(result, "function multiply(a, b)") {
		t.Error("Generated code should contain the function without type annotations")
	}

	if !strings.Contains(result, "const result = multiply(4, 5)") {
		t.Error("Generated code should call the function with test inputs")
	}

	if !strings.Contains(result, "const expected = 20") {
		t.Error("Generated code should set the expected value")
	}
}

func TestGeneratePHPTestCode(t *testing.T) {
	ch := challenge.Challenge{
		Language:     "php",
		FunctionName: "divide",
	}

	userCode := `<?php
function divide($a, $b) {
    return $a / $b;
}
?>`

	testCase := challenge.TestCase{
		Input:       []interface{}{float64(10), float64(2)},
		Expected:    float64(5),
		Description: "should divide two numbers",
	}

	result := generatePHPTestCode(ch, userCode, testCase)

	if !strings.Contains(result, "<?php") {
		t.Error("Generated code should start with PHP opening tag")
	}

	if !strings.Contains(result, "$result = divide(10, 2)") {
		t.Error("Generated code should call the function with test inputs")
	}

	if !strings.Contains(result, "$expected = 5") {
		t.Error("Generated code should set the expected value")
	}
}

func TestCleanGoUserCode(t *testing.T) {
	userCode := `package main

import (
    "fmt"
    "strconv"
)

import "os"

func add(a, b int) int {
    return a + b
}

func helper() string {
    return "helper"
}

func main() {
    fmt.Println("This should be removed")
    if true {
        fmt.Println("Nested content")
    }
}

func anotherFunc() int {
    return 42
}`

	cleaned := cleanGoUserCode(userCode)

	// Should remove package declaration
	if strings.Contains(cleaned, "package main") {
		t.Error("Cleaned code should not contain package declaration")
	}

	// Should remove import statements
	if strings.Contains(cleaned, "import") {
		t.Error("Cleaned code should not contain import statements")
	}

	// Should remove main function
	if strings.Contains(cleaned, "func main()") {
		t.Error("Cleaned code should not contain main function")
	}
	if strings.Contains(cleaned, "This should be removed") {
		t.Error("Cleaned code should not contain main function content")
	}

	// Should keep other functions
	if !strings.Contains(cleaned, "func add(a, b int) int") {
		t.Error("Cleaned code should contain the add function")
	}

	if !strings.Contains(cleaned, "func helper() string") {
		t.Error("Cleaned code should contain the helper function")
	}

	if !strings.Contains(cleaned, "func anotherFunc() int") {
		t.Error("Cleaned code should contain the anotherFunc function")
	}
}

func TestCleanGoUserCodeWithMainFunctionVariations(t *testing.T) {
	// Test main function with opening brace on same line
	userCode1 := `package main

func add(a, b int) int {
    return a + b
}

func main() {
    fmt.Println("Hello")
}`

	cleaned1 := cleanGoUserCode(userCode1)
	if strings.Contains(cleaned1, "func main()") || strings.Contains(cleaned1, "Hello") {
		t.Error("Should remove main function with brace on same line")
	}

	// Test main function with opening brace on next line
	userCode2 := `package main

func add(a, b int) int {
    return a + b
}

func main()
{
    fmt.Println("Hello")
}`

	cleaned2 := cleanGoUserCode(userCode2)
	if strings.Contains(cleaned2, "func main()") || strings.Contains(cleaned2, "Hello") {
		t.Error("Should remove main function with brace on next line")
	}
}
