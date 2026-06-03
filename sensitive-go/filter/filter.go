package filter

// Filter is the interface for filtering matched words
type Filter interface {
	// ShouldFilter returns true if the word should be filtered out (not reported)
	ShouldFilter(word string) bool

	// Name returns the filter name
	Name() string
}


