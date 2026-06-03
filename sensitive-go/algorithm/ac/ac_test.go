package ac

import (
	"testing"

	"github.com/Karrecy/sensitive-go/dict"
)

func TestACMatcher_Build(t *testing.T) {
	matcher := NewACMatcher(true) // case-sensitive
	words := []dict.Word{
		{Text: "敏感词", Category: dict.CategoryOther, Level: dict.LevelMedium},
		{Text: "测试", Category: dict.CategoryOther, Level: dict.LevelLow},
	}

	err := matcher.Build(words)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if matcher.root == nil {
		t.Fatal("Root node is nil")
	}
}

func TestACMatcher_Match(t *testing.T) {
	matcher := NewACMatcher(true) // case-sensitive
	words := []dict.Word{
		{Text: "敏感", Category: dict.CategoryAbuse, Level: dict.LevelHigh},
		{Text: "测试", Category: dict.CategoryOther, Level: dict.LevelLow},
		{Text: "敏感词", Category: dict.CategoryAbuse, Level: dict.LevelHigh},
	}

	err := matcher.Build(words)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	tests := []struct {
		name          string
		text          string
		expectedCount int
	}{
		{"Single match", "这是一个敏感词测试", 3}, // "敏感", "敏感词", "测试"
		{"Multiple matches", "敏感词和测试都是敏感内容", 4}, // "敏感", "敏感词", "测试", "敏感"
		{"No match", "正常文本", 0},
		{"Empty text", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := matcher.Match(tt.text)
			if len(matches) != tt.expectedCount {
				t.Errorf("Expected %d matches, got %d", tt.expectedCount, len(matches))
			}
		})
	}
}

func TestACMatcher_Replace(t *testing.T) {
	matcher := NewACMatcher(true) // case-sensitive
	words := []dict.Word{
		{Text: "敏感", Category: dict.CategoryAbuse, Level: dict.LevelHigh},
		{Text: "测试", Category: dict.CategoryOther, Level: dict.LevelLow},
	}

	err := matcher.Build(words)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	tests := []struct {
		name     string
		text     string
		repl     rune
		expected string
	}{
		{"Replace with asterisk", "这是敏感测试", '*', "这是****"},
		{"Replace with X", "敏感内容测试", 'X', "XX内容XX"},
		{"No sensitive words", "正常文本", '*', "正常文本"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.Replace(tt.text, tt.repl)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestACMatcher_Validate(t *testing.T) {
	matcher := NewACMatcher(true) // case-sensitive
	words := []dict.Word{
		{Text: "敏感", Category: dict.CategoryAbuse, Level: dict.LevelHigh},
		{Text: "测试", Category: dict.CategoryOther, Level: dict.LevelLow},
	}

	err := matcher.Build(words)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"Contains sensitive word", "这是敏感内容", false},
		{"Clean text", "正常文本", true},
		{"Empty text", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.Validate(tt.text)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func BenchmarkACMatcher_Match(b *testing.B) {
	matcher := NewACMatcher(true) // case-sensitive
	words := make([]dict.Word, 1000)
	for i := 0; i < 1000; i++ {
		words[i] = dict.Word{
			Text:     "测试词" + string(rune(i)),
			Category: dict.CategoryOther,
			Level:    dict.LevelLow,
		}
	}

	matcher.Build(words)
	text := "这是一段包含测试词100的文本内容，用于性能基准测试。"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matcher.Match(text)
	}
}

