package builtin

import "testing"

func TestGetDefaultWords(t *testing.T) {
	words := GetDefaultWords()

	if len(words) == 0 {
		t.Error("Default words should not be empty")
	}

	// Check that comments are filtered out
	for _, word := range words {
		if word.Text == "" {
			t.Error("Word text should not be empty")
		}
		if word.Text[0] == '#' {
			t.Error("Comments should be filtered out")
		}
	}
}

func TestParseWords(t *testing.T) {
	content := `# Comment
word1
word2

# Another comment
word3`

	words := parseWords(content)

	if len(words) != 3 {
		t.Errorf("Expected 3 words, got %d", len(words))
	}

	expected := []string{"word1", "word2", "word3"}
	for i, word := range words {
		if word.Text != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], word.Text)
		}
	}
}
