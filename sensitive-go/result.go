package gosensitive

import "github.com/Karrecy/sensitive-go/dict"

// Result represents the detection result
type Result struct {
	Found        bool    // Whether any sensitive words were found
	Matches      []Match // List of all matches
	FilteredText string  // Text with sensitive words filtered/replaced
}

// Match represents a single sensitive word match
type Match struct {
	Word     string        // The matched sensitive word
	Start    int           // Start position in runes (not bytes)
	End      int           // End position in runes (not bytes)
	Category dict.Category // Category of the matched word
	Level    dict.Level    // Severity level of the matched word
}

// HasCategory checks if the result contains matches of the specified category
func (r *Result) HasCategory(category dict.Category) bool {
	for _, match := range r.Matches {
		if match.Category.Has(category) {
			return true
		}
	}
	return false
}

// HasLevel checks if the result contains matches of the specified level
func (r *Result) HasLevel(level dict.Level) bool {
	for _, match := range r.Matches {
		if match.Level == level {
			return true
		}
	}
	return false
}

// FilterByCategory returns matches that belong to the specified category
func (r *Result) FilterByCategory(category dict.Category) []Match {
	var filtered []Match
	for _, match := range r.Matches {
		if match.Category.Has(category) {
			filtered = append(filtered, match)
		}
	}
	return filtered
}

// FilterByLevel returns matches that have the specified severity level
func (r *Result) FilterByLevel(level dict.Level) []Match {
	var filtered []Match
	for _, match := range r.Matches {
		if match.Level == level {
			filtered = append(filtered, match)
		}
	}
	return filtered
}


