package variant

import "testing"

func TestPinyinProcessor_Process(t *testing.T) {
	processor := NewPinyinProcessor()

	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{"Chinese to pinyin", "测试", "ce"},
		{"Mixed content", "测试ABC", "ce"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Process(tt.input)
			if result == "" {
				t.Errorf("Process returned empty string")
			}
		})
	}
}

func TestTraditionalProcessor_Process(t *testing.T) {
	processor := NewTraditionalProcessor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Traditional to simplified", "測試", "测试"},
		{"Already simplified", "测试", "测试"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ToSimplified(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSymbolProcessor_Process(t *testing.T) {
	processor := NewSymbolProcessor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Remove symbols", "测@试#词", "测试词"},
		{"Keep spaces", "测 试 词", "测 试 词"},
		{"Mixed content", "测试123ABC", "测试123ABC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Process(tt.input)
			// Normalize whitespace for comparison
			result = NormalizeWhitespace(result)
			expected := NormalizeWhitespace(tt.expected)
			if result != expected {
				t.Errorf("Expected %q, got %q", expected, result)
			}
		})
	}
}

func TestSimilarProcessor_Process(t *testing.T) {
	processor := NewSimilarProcessor()

	tests := []struct {
		name  string
		input string
		check func(string) bool
	}{
		{
			"Similar character normalization",
			"傻",
			func(result string) bool {
				return result == "傻"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Process(tt.input)
			if !tt.check(result) {
				t.Errorf("Process result failed check for input %q", tt.input)
			}
		})
	}
}

func TestIsSimilar(t *testing.T) {
	tests := []struct {
		name     string
		r1       rune
		r2       rune
		expected bool
	}{
		{"Same character", '傻', '傻', true},
		{"Different characters", '傻', '测', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSimilar(tt.r1, tt.r2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}


