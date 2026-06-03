package loader

import (
	"os"
	"testing"

	"github.com/Karrecy/sensitive-go/dict"
)

func TestMemoryLoader_Load(t *testing.T) {
	words := []string{"敏感词", "测试", ""}
	loader := NewMemoryLoader(words)

	result, err := loader.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Empty strings should be filtered out
	if len(result) != 2 {
		t.Errorf("Expected 2 words, got %d", len(result))
	}

	if result[0].Text != "敏感词" {
		t.Errorf("Expected '敏感词', got '%s'", result[0].Text)
	}
}

func TestFileLoader_LoadTXT(t *testing.T) {
	// Create temporary test file
	tmpFile, err := os.CreateTemp("", "test_words_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "敏感词\n测试\n# 这是注释\n\n另一个词"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	loader := NewFileLoader(tmpFile.Name())
	result, err := loader.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should load 3 words (comment and empty line skipped)
	if len(result) != 3 {
		t.Errorf("Expected 3 words, got %d", len(result))
	}
}

func TestFileLoader_LoadJSON(t *testing.T) {
	// Create temporary test file
	tmpFile, err := os.CreateTemp("", "test_words_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `[
		{"Text": "敏感词", "Category": 1, "Level": 2},
		{"Text": "测试", "Category": 2, "Level": 1}
	]`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	loader := NewFileLoader(tmpFile.Name())
	result, err := loader.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 words, got %d", len(result))
	}

	if result[0].Level != dict.LevelHigh {
		t.Errorf("Expected level %v, got %v", dict.LevelHigh, result[0].Level)
	}
}

func TestFileLoader_NonExistentFile(t *testing.T) {
	loader := NewFileLoader("/nonexistent/file.txt")
	_, err := loader.Load()
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}


