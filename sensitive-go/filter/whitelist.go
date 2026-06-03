package filter

import "strings"

// Whitelist filters out words that are in the whitelist
type Whitelist struct {
	words map[string]bool
}

// NewWhitelist creates a new whitelist filter
func NewWhitelist(words []string) *Whitelist {
	wordMap := make(map[string]bool)
	for _, word := range words {
		wordMap[strings.ToLower(word)] = true
	}

	return &Whitelist{
		words: wordMap,
	}
}

// ShouldFilter returns true if the word is in the whitelist
func (w *Whitelist) ShouldFilter(word string) bool {
	return w.words[strings.ToLower(word)]
}

// Name returns the filter name
func (w *Whitelist) Name() string {
	return "whitelist"
}

// Add adds a word to the whitelist
func (w *Whitelist) Add(word string) {
	w.words[strings.ToLower(word)] = true
}

// Remove removes a word from the whitelist
func (w *Whitelist) Remove(word string) {
	delete(w.words, strings.ToLower(word))
}

// Contains checks if a word is in the whitelist
func (w *Whitelist) Contains(word string) bool {
	return w.words[strings.ToLower(word)]
}

// Clear removes all words from the whitelist
func (w *Whitelist) Clear() {
	w.words = make(map[string]bool)
}


