package algorithm

import "github.com/Karrecy/sensitive-go/dict"

// Matcher is the interface that all matching algorithms must implement
type Matcher interface {
	// Build constructs the detection structure from the given words
	Build(words []dict.Word) error

	// Match finds all sensitive words in the text and returns their positions
	Match(text string) []MatchResult

	// Replace replaces all sensitive words in the text with the given replacement rune
	Replace(text string, repl rune) string

	// Validate checks if the text contains any sensitive words
	Validate(text string) bool
}

// MatchResult represents a single match result from the algorithm
type MatchResult struct {
	Word     string        // The matched word
	Start    int           // Start position (rune index)
	End      int           // End position (rune index)
	Category dict.Category // Category of the word
	Level    dict.Level    // Severity level
}

// AlgorithmType represents the type of matching algorithm
type AlgorithmType int

const (
	// AlgorithmAuto automatically selects the best algorithm based on word count
	AlgorithmAuto AlgorithmType = iota
	// AlgorithmDFA uses Deterministic Finite Automaton
	AlgorithmDFA
	// AlgorithmAC uses Aho-Corasick algorithm
	AlgorithmAC
)

// String returns the string representation of the algorithm type
func (a AlgorithmType) String() string {
	switch a {
	case AlgorithmAuto:
		return "auto"
	case AlgorithmDFA:
		return "dfa"
	case AlgorithmAC:
		return "ac"
	default:
		return "unknown"
	}
}


