package dict

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

// String returns the string representation of the category
func (c Category) String() string {
	switch c {
	case CategoryPolitical:
		return "political"
	case CategoryPornographic:
		return "pornographic"
	case CategoryViolence:
		return "violence"
	case CategoryAbuse:
		return "abuse"
	case CategoryAd:
		return "ad"
	case CategoryIllegal:
		return "illegal"
	case CategoryOther:
		return "other"
	default:
		return "unknown"
	}
}

// Has checks if the category contains the specified flag
func (c Category) Has(flag Category) bool {
	return c&flag != 0
}

// Add adds a category flag
func (c Category) Add(flag Category) Category {
	return c | flag
}

// Remove removes a category flag
func (c Category) Remove(flag Category) Category {
	return c &^ flag
}


