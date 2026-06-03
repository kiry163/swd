package dfa

import (
	"strings"
	"unicode"

	"github.com/Karrecy/sensitive-go/algorithm"
	"github.com/Karrecy/sensitive-go/dict"
)

// DFAMatcher implements a Deterministic Finite Automaton for pattern matching
type DFAMatcher struct {
	root          *State // Root state of the DFA
	caseSensitive bool   // Whether matching is case-sensitive
}

// NewDFAMatcher creates a new DFA matcher instance
func NewDFAMatcher(caseSensitive bool) *DFAMatcher {
	return &DFAMatcher{
		root:          newState(),
		caseSensitive: caseSensitive,
	}
}

// Build constructs the DFA from the given words
func (m *DFAMatcher) Build(words []dict.Word) error {
	for i := range words {
		m.insert(&words[i])
	}
	return nil
}

// insert adds a word to the DFA
func (m *DFAMatcher) insert(word *dict.Word) {
	state := m.root
	text := word.Text

	// Convert to lowercase if case-insensitive
	if !m.caseSensitive {
		text = strings.ToLower(text)
	}

	runes := []rune(text)

	for _, r := range runes {
		if next, exists := state.transition(r); exists {
			state = next
		} else {
			nextState := newState()
			state.addTransition(r, nextState)
			state = nextState
		}
	}

	state.setWord(word)
}

// Match finds all sensitive words in the text
func (m *DFAMatcher) Match(text string) []algorithm.MatchResult {
	results := make([]algorithm.MatchResult, 0)

	// Convert to lowercase if case-insensitive
	if !m.caseSensitive {
		text = strings.Map(unicode.ToLower, text)
	}

	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		state := m.root
		j := i

		// Try to match from position i
		for j < len(runes) {
			next, exists := state.transition(runes[j])
			if !exists {
				break
			}

			state = next
			j++

			// Check if we've reached an end state
			if state.isEndState() {
				word := state.getWord()
				if word != nil {
					results = append(results, algorithm.MatchResult{
						Word:     word.Text,
						Start:    i,
						End:      j,
						Category: word.Category,
						Level:    word.Level,
					})
				}
			}
		}
	}

	return results
}

// Replace replaces all sensitive words with the given replacement rune
func (m *DFAMatcher) Replace(text string, repl rune) string {
	// Convert to lowercase if case-insensitive for consistent matching
	processedText := text
	if !m.caseSensitive {
		processedText = strings.Map(unicode.ToLower, text)
	}

	runes := []rune(processedText)
	matches := m.Match(text)

	// Mark positions to replace
	toReplace := make([]bool, len(runes))
	for _, match := range matches {
		for i := match.Start; i < match.End; i++ {
			toReplace[i] = true
		}
	}

	// Replace marked positions
	for i := range runes {
		if toReplace[i] {
			runes[i] = repl
		}
	}

	return string(runes)
}

// Validate checks if the text contains any sensitive words
func (m *DFAMatcher) Validate(text string) bool {
	// Convert to lowercase if case-insensitive
	if !m.caseSensitive {
		text = strings.Map(unicode.ToLower, text)
	}

	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		state := m.root
		j := i

		// Try to match from position i
		for j < len(runes) {
			next, exists := state.transition(runes[j])
			if !exists {
				break
			}

			state = next
			j++

			// If we find a match, text is invalid
			if state.isEndState() {
				return false
			}
		}
	}

	return true // No sensitive words found
}
