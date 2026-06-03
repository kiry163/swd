package ac

import (
	"strings"
	"unicode"

	"github.com/Karrecy/sensitive-go/algorithm"
	"github.com/Karrecy/sensitive-go/dict"
)

// ACMatcher implements the Aho-Corasick algorithm for multi-pattern matching
type ACMatcher struct {
	root          *Node // Root of the trie
	caseSensitive bool  // Whether matching is case-sensitive
}

// NewACMatcher creates a new AC matcher instance
func NewACMatcher(caseSensitive bool) *ACMatcher {
	return &ACMatcher{
		root:          newNode(),
		caseSensitive: caseSensitive,
	}
}

// Build constructs the AC automaton from the given words
func (m *ACMatcher) Build(words []dict.Word) error {
	// Build the trie
	for i := range words {
		m.insert(&words[i])
	}

	// Build failure pointers using BFS
	m.buildFailurePointers()

	return nil
}

// insert adds a word to the trie
func (m *ACMatcher) insert(word *dict.Word) {
	node := m.root
	text := word.Text

	// Convert to lowercase if case-insensitive
	if !m.caseSensitive {
		text = strings.ToLower(text)
	}

	runes := []rune(text)

	for _, r := range runes {
		node = node.addChild(r)
	}

	node.setWord(word)
}

// buildFailurePointers constructs failure pointers for the AC automaton
func (m *ACMatcher) buildFailurePointers() {
	queue := make([]*Node, 0)

	// Initialize: all children of root have failure pointer to root
	for _, child := range m.root.children {
		child.fail = m.root
		queue = append(queue, child)
	}

	// BFS to build failure pointers
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for r, child := range current.children {
			queue = append(queue, child)

			// Find failure pointer
			failNode := current.fail
			for failNode != nil {
				if next, exists := failNode.getChild(r); exists {
					child.fail = next
					break
				}
				if failNode == m.root {
					child.fail = m.root
					break
				}
				failNode = failNode.fail
			}
		}
	}
}

// Match finds all sensitive words in the text
func (m *ACMatcher) Match(text string) []algorithm.MatchResult {
	results := make([]algorithm.MatchResult, 0)

	// Convert to lowercase if case-insensitive
	if !m.caseSensitive {
		text = strings.Map(unicode.ToLower, text)
	}

	runes := []rune(text)
	node := m.root

	for i, r := range runes {
		// Follow failure pointers until we find a match or reach root
		for node != m.root && !node.hasChild(r) {
			node = node.fail
		}

		// Try to match the character
		if child, exists := node.getChild(r); exists {
			node = child
		} else {
			node = m.root
			continue
		}

		// Check for matches at current node and all its failure nodes
		tempNode := node
		for tempNode != m.root {
			if tempNode.isEnd && tempNode.word != nil {
				wordLen := len([]rune(tempNode.word.Text))
				results = append(results, algorithm.MatchResult{
					Word:     tempNode.word.Text,
					Start:    i - wordLen + 1,
					End:      i + 1,
					Category: tempNode.word.Category,
					Level:    tempNode.word.Level,
				})
			}
			tempNode = tempNode.fail
		}
	}

	return results
}

// Replace replaces all sensitive words with the given replacement rune
func (m *ACMatcher) Replace(text string, repl rune) string {
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
func (m *ACMatcher) Validate(text string) bool {
	// Convert to lowercase if case-insensitive
	if !m.caseSensitive {
		text = strings.Map(unicode.ToLower, text)
	}

	runes := []rune(text)
	node := m.root

	for _, r := range runes {
		// Follow failure pointers
		for node != m.root && !node.hasChild(r) {
			node = node.fail
		}

		// Try to match
		if child, exists := node.getChild(r); exists {
			node = child
		} else {
			node = m.root
			continue
		}

		// Check for match
		tempNode := node
		for tempNode != m.root {
			if tempNode.isEnd {
				return false // Found a sensitive word
			}
			tempNode = tempNode.fail
		}
	}

	return true // No sensitive words found
}
