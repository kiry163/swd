package loader

import "github.com/Karrecy/sensitive-go/dict"

// MemoryLoader loads sensitive words from an in-memory slice
type MemoryLoader struct {
	words []string
}

// NewMemoryLoader creates a new memory loader
func NewMemoryLoader(words []string) *MemoryLoader {
	return &MemoryLoader{words: words}
}

// Load converts string slice to Word slice
func (l *MemoryLoader) Load() ([]dict.Word, error) {
	words := make([]dict.Word, 0, len(l.words))

	for _, w := range l.words {
		if w == "" {
			continue
		}

		words = append(words, dict.Word{
			Text:     w,
			Category: dict.CategoryOther,
			Level:    dict.LevelMedium,
		})
	}

	return words, nil
}


