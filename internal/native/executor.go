package native

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ExecutionResult struct {
	Success  bool
	Output   string
	Error    string
	Duration time.Duration
	ExitCode int
}

type Executor struct {
	workDir string
}

func NewExecutor() (*Executor, error) {
	// Create a temporary working directory for code execution
	workDir, err := os.MkdirTemp("", "codequest-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	return &Executor{workDir: workDir}, nil
}

func (e *Executor) Close() error {
	return os.RemoveAll(e.workDir)
}

func (e *Executor) ExecuteCode(language, code string, timeLimit int) (*ExecutionResult, error) {
	start := time.Now()

	// Create language-specific executor
	var executor LanguageExecutor
	switch language {
	case "go":
		executor = &GoExecutor{workDir: e.workDir}
	case "javascript", "typescript":
		executor = &NodeExecutor{workDir: e.workDir}
	case "python":
		executor = &PythonExecutor{workDir: e.workDir}
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	// Check if the language runtime is available
	if err := executor.CheckAvailability(); err != nil {
		return nil, fmt.Errorf("language runtime not available: %w", err)
	}

	// Execute the code
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLimit)*time.Millisecond)
	defer cancel()

	result, err := executor.Execute(ctx, code)
	if err != nil {
		return nil, err
	}

	result.Duration = time.Since(start)
	return result, nil
}

// LanguageExecutor interface for different language executors
type LanguageExecutor interface {
	CheckAvailability() error
	Execute(ctx context.Context, code string) (*ExecutionResult, error)
}

// GoExecutor implements Go code execution
type GoExecutor struct {
	workDir string
}

func (g *GoExecutor) CheckAvailability() error {
	_, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("Go runtime not found. Please install Go from https://golang.org/dl/")
	}
	return nil
}

func (g *GoExecutor) Execute(ctx context.Context, code string) (*ExecutionResult, error) {
	// Create a unique subdirectory for this execution
	execDir := filepath.Join(g.workDir, fmt.Sprintf("go-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(execDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create execution directory: %w", err)
	}
	defer os.RemoveAll(execDir)

	// Write code to main.go
	mainFile := filepath.Join(execDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write Go code: %w", err)
	}

	// Initialize go module
	modCmd := exec.CommandContext(ctx, "go", "mod", "init", "solution")
	modCmd.Dir = execDir
	if err := modCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to initialize Go module: %w", err)
	}

	// Execute go run
	cmd := exec.CommandContext(ctx, "go", "run", "main.go")
	cmd.Dir = execDir

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	result := &ExecutionResult{
		Success:  err == nil,
		Output:   outputStr,
		ExitCode: 0,
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1
		}
		result.Error = err.Error()
	}

	return result, nil
}

// NodeExecutor implements Node.js code execution
type NodeExecutor struct {
	workDir string
}

func (n *NodeExecutor) CheckAvailability() error {
	_, err := exec.LookPath("node")
	if err != nil {
		return fmt.Errorf("Node.js runtime not found. Please install Node.js from https://nodejs.org/")
	}
	return nil
}

func (n *NodeExecutor) Execute(ctx context.Context, code string) (*ExecutionResult, error) {
	// Create a unique subdirectory for this execution
	execDir := filepath.Join(n.workDir, fmt.Sprintf("node-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(execDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create execution directory: %w", err)
	}
	defer os.RemoveAll(execDir)

	// For TypeScript, we'll transpile to JavaScript (basic approach)
	jsCode := code
	if strings.Contains(code, ": number") || strings.Contains(code, ": string") || strings.Contains(code, ": boolean") {
		jsCode = transpileTypeScript(code)
	}

	// Write code to solution.js
	jsFile := filepath.Join(execDir, "solution.js")
	if err := os.WriteFile(jsFile, []byte(jsCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write JavaScript code: %w", err)
	}

	// Execute node
	cmd := exec.CommandContext(ctx, "node", "solution.js")
	cmd.Dir = execDir

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	result := &ExecutionResult{
		Success:  err == nil,
		Output:   outputStr,
		ExitCode: 0,
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1
		}
		result.Error = err.Error()
	}

	return result, nil
}

// Basic TypeScript to JavaScript transpilation
func transpileTypeScript(code string) string {
	// Remove common TypeScript type annotations
	jsCode := strings.ReplaceAll(code, ": number", "")
	jsCode = strings.ReplaceAll(jsCode, ": string", "")
	jsCode = strings.ReplaceAll(jsCode, ": boolean", "")
	jsCode = strings.ReplaceAll(jsCode, ": void", "")
	jsCode = strings.ReplaceAll(jsCode, ": any", "")
	
	// Remove array type annotations like number[], string[]
	jsCode = strings.ReplaceAll(jsCode, ": number[]", "")
	jsCode = strings.ReplaceAll(jsCode, ": string[]", "")
	jsCode = strings.ReplaceAll(jsCode, ": boolean[]", "")
	
	return jsCode
}

// PythonExecutor implements Python code execution
type PythonExecutor struct {
	workDir string
}

func (p *PythonExecutor) CheckAvailability() error {
	_, err := exec.LookPath("python3")
	if err != nil {
		// Try fallback to python
		_, err = exec.LookPath("python")
		if err != nil {
			return fmt.Errorf("Python runtime not found. Please install Python from https://python.org/downloads/")
		}
	}
	return nil
}

func (p *PythonExecutor) Execute(ctx context.Context, code string) (*ExecutionResult, error) {
	// Create a unique subdirectory for this execution
	execDir := filepath.Join(p.workDir, fmt.Sprintf("python-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(execDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create execution directory: %w", err)
	}
	defer os.RemoveAll(execDir)

	// Write code to solution.py
	pyFile := filepath.Join(execDir, "solution.py")
	if err := os.WriteFile(pyFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write Python code: %w", err)
	}

	// Try python3 first, then python
	var cmd *exec.Cmd
	if _, err := exec.LookPath("python3"); err == nil {
		cmd = exec.CommandContext(ctx, "python3", "solution.py")
	} else {
		cmd = exec.CommandContext(ctx, "python", "solution.py")
	}
	cmd.Dir = execDir

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	result := &ExecutionResult{
		Success:  err == nil,
		Output:   outputStr,
		ExitCode: 0,
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1
		}
		result.Error = err.Error()
	}

	return result, nil
}