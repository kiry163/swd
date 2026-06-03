package variant

import (
	"strings"
	"unicode"
)

// SymbolProcessor handles symbol interference filtering
type SymbolProcessor struct {
	removeSymbols bool
}

// NewSymbolProcessor creates a new symbol processor
func NewSymbolProcessor() *SymbolProcessor {
	return &SymbolProcessor{
		removeSymbols: true,
	}
}

// Process removes or normalizes symbols from text
func (p *SymbolProcessor) Process(text string) string {
	if !p.removeSymbols {
		return text
	}

	var builder strings.Builder
	builder.Grow(len(text))

	for _, r := range text {
		// Keep alphanumeric and CJK characters
		if unicode.IsLetter(r) || unicode.IsDigit(r) || isCJK(r) {
			builder.WriteRune(r)
		} else if unicode.IsSpace(r) {
			builder.WriteRune(' ')
		}
		// Skip other symbols
	}

	return builder.String()
}

// Name returns the processor name
func (p *SymbolProcessor) Name() string {
	return "symbol"
}

// isCJK checks if a rune is a CJK character
func isCJK(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
		(r >= 0x3400 && r <= 0x4DBF) || // CJK Extension A
		(r >= 0x20000 && r <= 0x2A6DF) || // CJK Extension B
		(r >= 0x2A700 && r <= 0x2B73F) || // CJK Extension C
		(r >= 0x2B740 && r <= 0x2B81F) || // CJK Extension D
		(r >= 0x2B820 && r <= 0x2CEAF) // CJK Extension E
}

// RemoveSymbols is a utility function to remove symbols from text
func RemoveSymbols(text string) string {
	processor := NewSymbolProcessor()
	return processor.Process(text)
}

// NormalizeWhitespace replaces multiple spaces with a single space
func NormalizeWhitespace(text string) string {
	return strings.Join(strings.Fields(text), " ")
}


