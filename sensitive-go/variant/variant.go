package variant

// Processor is the interface for text variant processing
type Processor interface {
	// Process transforms the text to handle variants
	Process(text string) string

	// Name returns the name of the processor
	Name() string
}


