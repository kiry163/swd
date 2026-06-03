package filter

import "testing"

func TestWhitelist_ShouldFilter(t *testing.T) {
	whitelist := NewWhitelist([]string{"测试", "示例"})

	tests := []struct {
		name     string
		word     string
		expected bool
	}{
		{"Word in whitelist", "测试", true},
		{"Word not in whitelist", "敏感", false},
		{"Case insensitive", "测试", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := whitelist.ShouldFilter(tt.word)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestWhitelist_Add(t *testing.T) {
	whitelist := NewWhitelist([]string{})

	whitelist.Add("新词")
	if !whitelist.Contains("新词") {
		t.Error("Failed to add word to whitelist")
	}
}

func TestWhitelist_Remove(t *testing.T) {
	whitelist := NewWhitelist([]string{"测试"})

	whitelist.Remove("测试")
	if whitelist.Contains("测试") {
		t.Error("Failed to remove word from whitelist")
	}
}

func TestWhitelist_Clear(t *testing.T) {
	whitelist := NewWhitelist([]string{"测试", "示例"})

	whitelist.Clear()
	if whitelist.Contains("测试") || whitelist.Contains("示例") {
		t.Error("Failed to clear whitelist")
	}
}


