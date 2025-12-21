package parser

import (
	"strings"
	"unicode"
)

// Command represents a parsed download command
type Command struct {
	RemotePath string
	LocalTarget string
}

var prefixes = []string{":d", ":download", "~d", "~download"}

// Parse checks if the input line matches a download command and extracts the arguments
func Parse(line string) (*Command, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return nil, false
	}

	lower := strings.ToLower(trimmed)
	
	// Check if line starts with any of the prefixes (case-insensitive)
	var matchedPrefix string
	var prefixLen int
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, strings.ToLower(prefix)) {
			matchedPrefix = prefix
			prefixLen = len(prefix)
			break
		}
	}

	if matchedPrefix == "" {
		return nil, false
	}

	// Extract the remainder after the prefix
	remainder := trimmed[prefixLen:]
	remainder = strings.TrimSpace(remainder)

	if remainder == "" {
		return nil, false
	}

	// Tokenize the arguments (split on whitespace, but respect quoted strings)
	tokens := tokenize(remainder)
	
	if len(tokens) == 0 {
		return nil, false
	}

	cmd := &Command{
		RemotePath: tokens[0],
	}

	if len(tokens) > 1 {
		cmd.LocalTarget = tokens[1]
	}

	return cmd, true
}

// tokenize splits a string into tokens, handling quoted strings
func tokenize(s string) []string {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(s); i++ {
		char := s[i]
		
		if (char == '"' || char == '\'') && (i == 0 || s[i-1] != '\\') {
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else {
				current.WriteByte(char)
			}
		} else if unicode.IsSpace(rune(char)) && !inQuotes {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		} else {
			current.WriteByte(char)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}


