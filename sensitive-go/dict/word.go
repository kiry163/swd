package dict

// Word represents a sensitive word with metadata
type Word struct {
	Text     string   // The sensitive word text
	Category Category // Category of the word
	Level    Level    // Severity level
	Tags     []string // Custom tags for the word
}

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

// String returns the string representation of the level
func (l Level) String() string {
	switch l {
	case LevelLow:
		return "low"
	case LevelMedium:
		return "medium"
	case LevelHigh:
		return "high"
	case LevelCritical:
		return "critical"
	default:
		return "unknown"
	}
}


