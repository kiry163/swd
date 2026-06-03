package builtin

import (
	_ "embed"
	"strings"

	"github.com/Karrecy/sensitive-go/dict"
)

//go:embed data/default.txt
var defaultWords string

// GetDefaultWords returns the built-in default word dictionary
func GetDefaultWords() []dict.Word {
	return parseWords(defaultWords)
}

// parseWords parses text content into Word slice
func parseWords(content string) []dict.Word {
	lines := strings.Split(content, "\n")
	words := make([]dict.Word, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		words = append(words, dict.Word{
			Text:     line,
			Category: dict.CategoryOther,
			Level:    dict.LevelMedium,
		})
	}

	return words
}
