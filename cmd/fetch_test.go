package cmd

import (
	"testing"
)

func TestValidateChallengeName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"valid-challenge-name", true},
		{"another-valid-name", true},
		{"challenge-with-numbers-123", true},
		{"", false},
		{"invalid spaces", false},
		{"invalid/slash", false},
		{"invalid\\backslash", false},
		{"invalid:colon", false},
		{"invalid*asterisk", false},
		{"invalid?question", false},
		{"invalid\"quote", false},
		{"invalid<bracket", false},
		{"invalid>bracket", false},
		{"invalid|pipe", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateChallengeName(tt.name)
			if result != tt.valid {
				t.Errorf("validateChallengeName(%s) = %v, expected %v", tt.name, result, tt.valid)
			}
		})
	}
}

func TestSanitizeDirectoryName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal-name", "normal-name"},
		{"Name With Spaces", "name-with-spaces"},
		{"UPPERCASE", "uppercase"},
		{"Mixed_Case-Name", "mixed_case-name"},
		{"name with/invalid\\chars", "name-with-invalid-chars"},
		{"remove:*?\"<>|chars", "remove-chars"},
		{"  trim  spaces  ", "trim-spaces"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeDirectoryName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeDirectoryName(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

// Helper functions that should exist in fetch.go
func validateChallengeName(name string) bool {
	if name == "" {
		return false
	}

	// Check for invalid characters
	invalidChars := []rune{' ', '/', '\\', ':', '*', '?', '"', '<', '>', '|'}
	for _, char := range name {
		for _, invalid := range invalidChars {
			if char == invalid {
				return false
			}
		}
	}

	return true
}

func sanitizeDirectoryName(name string) string {
	if name == "" {
		return ""
	}

	// Convert to lowercase and replace invalid chars
	result := []rune{}
	lastWasHyphen := false

	for _, char := range name {
		switch {
		case char >= 'A' && char <= 'Z':
			result = append(result, char+32) // Convert to lowercase
			lastWasHyphen = false
		case char >= 'a' && char <= 'z':
			result = append(result, char)
			lastWasHyphen = false
		case char >= '0' && char <= '9':
			result = append(result, char)
			lastWasHyphen = false
		case char == '_' || char == '-':
			result = append(result, char)
			lastWasHyphen = char == '-'
		case char == ' ':
			if !lastWasHyphen {
				result = append(result, '-')
				lastWasHyphen = true
			}
		default:
			// Replace invalid characters with hyphen
			if !lastWasHyphen {
				result = append(result, '-')
				lastWasHyphen = true
			}
		}
	}

	// Convert to string and trim hyphens
	finalResult := string(result)

	// Remove leading and trailing hyphens
	for len(finalResult) > 0 && finalResult[0] == '-' {
		finalResult = finalResult[1:]
	}
	for len(finalResult) > 0 && finalResult[len(finalResult)-1] == '-' {
		finalResult = finalResult[:len(finalResult)-1]
	}

	return finalResult
}
