package loader

import "github.com/Karrecy/sensitive-go/dict"

// Loader is the interface for loading sensitive words from various sources
type Loader interface {
	// Load loads sensitive words and returns them as a slice
	Load() ([]dict.Word, error)
}


