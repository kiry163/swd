package dict

import "testing"

func TestCategory_String(t *testing.T) {
	tests := []struct {
		category Category
		expected string
	}{
		{CategoryPolitical, "political"},
		{CategoryPornographic, "pornographic"},
		{CategoryViolence, "violence"},
	}

	for _, tt := range tests {
		if tt.category.String() != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, tt.category.String())
		}
	}
}

func TestCategory_Has(t *testing.T) {
	cat := CategoryPolitical | CategoryViolence

	if !cat.Has(CategoryPolitical) {
		t.Error("Should have Political category")
	}

	if cat.Has(CategoryPornographic) {
		t.Error("Should not have Pornographic category")
	}
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{LevelLow, "low"},
		{LevelMedium, "medium"},
		{LevelHigh, "high"},
		{LevelCritical, "critical"},
	}

	for _, tt := range tests {
		if tt.level.String() != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, tt.level.String())
		}
	}
}


