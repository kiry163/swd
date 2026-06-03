package gosensitive

import "time"

// Options contains configuration options for the detector
type Options struct {
	// Algorithm specifies which matching algorithm to use
	Algorithm AlgorithmType

	// CaseSensitive determines if matching should be case-sensitive
	CaseSensitive bool

	// EnablePinyin enables pinyin variant detection
	EnablePinyin bool

	// EnableTraditional enables traditional Chinese variant detection
	EnableTraditional bool

	// EnableSymbolFilter enables filtering of symbol interference
	EnableSymbolFilter bool

	// EnableSimilarChar enables similar character detection
	EnableSimilarChar bool

	// ReplaceChar is the default character used for replacement
	ReplaceChar rune

	// Categories filters detection to only these categories (nil means all)
	Categories []Category

	// MinLevel is the minimum severity level to detect
	MinLevel Level

	// MaxMatchCount limits the maximum number of matches to return (0 means no limit)
	MaxMatchCount int

	// WatchFile enables automatic reloading when word files change
	WatchFile bool

	// WatchInterval is the interval for checking file changes
	WatchInterval time.Duration
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

// Category represents the category of sensitive words using bit flags
type Category int

const (
	// CategoryPolitical represents political sensitive words
	CategoryPolitical Category = 1 << iota
	// CategoryPornographic represents pornographic content
	CategoryPornographic
	// CategoryViolence represents violence and gore
	CategoryViolence
	// CategoryAbuse represents abusive and insulting words
	CategoryAbuse
	// CategoryAd represents advertisement and spam
	CategoryAd
	// CategoryIllegal represents illegal activities
	CategoryIllegal
	// CategoryOther represents other categories
	CategoryOther
)

// Level represents the severity level of a sensitive word
type Level int

const (
	// LevelLow indicates a low severity word
	LevelLow Level = iota
	// LevelMedium indicates a medium severity word
	LevelMedium
	// LevelHigh indicates a high severity word
	LevelHigh
	// LevelCritical indicates a critical severity word
	LevelCritical
)

// DefaultOptions returns the default options
func DefaultOptions() *Options {
	return &Options{
		Algorithm:          AlgorithmAuto,
		CaseSensitive:      false,
		EnablePinyin:       false,
		EnableTraditional:  false,
		EnableSymbolFilter: false,
		EnableSimilarChar:  false,
		ReplaceChar:        '*',
		Categories:         nil,
		MinLevel:           LevelLow,
		MaxMatchCount:      0,
		WatchFile:          false,
		WatchInterval:      time.Second * 30,
	}
}


