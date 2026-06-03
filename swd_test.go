package swd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEngine_FindAll_ReturnsRunePositions(t *testing.T) {
	eng := New(
		WithIgnoreSymbol(true),
		WithIgnoreWidth(true),
		WithIgnoreCase(true),
		WithTraditional(true),
		WithSimilarChar(true),
		WithPinyin(true),
	)

	if err := eng.AddWord(Word{Text: "测试", Type: "custom"}); err != nil {
		t.Fatalf("AddWord() error = %v", err)
	}
	if err := eng.Build(); err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	matches := eng.FindAll("这是一段 测@試 的文本")
	if len(matches) != 1 {
		t.Fatalf("FindAll() len = %d, want 1", len(matches))
	}

	got := matches[0]
	if got.Word != "测试" {
		t.Fatalf("Word = %q, want %q", got.Word, "测试")
	}
	if got.Type != "custom" {
		t.Fatalf("Type = %q, want %q", got.Type, "custom")
	}
	if got.StartPos != 5 || got.EndPos != 8 {
		t.Fatalf("StartPos/EndPos = %d/%d, want 5/8", got.StartPos, got.EndPos)
	}
}

func TestEngine_AddFile_LoadsWords(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "words.txt")
	if err := os.WriteFile(file, []byte("诈骗,crime\n赌博,crime\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	eng := New()
	if err := eng.AddFile(file); err != nil {
		t.Fatalf("AddFile() error = %v", err)
	}
	if err := eng.Build(); err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if !eng.Contains("这里有诈骗信息") {
		t.Fatal("Contains() = false, want true")
	}
}

func TestEngine_Replace(t *testing.T) {
	eng := New()
	if err := eng.AddWords([]Word{{Text: "坏词", Type: "custom"}}); err != nil {
		t.Fatalf("AddWords() error = %v", err)
	}
	if err := eng.Build(); err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	got := eng.Replace("这是坏词", "*")
	if got != "这是**" {
		t.Fatalf("Replace() = %q, want %q", got, "这是**")
	}
}

func TestEngine_FindAll_MultipleMatches(t *testing.T) {
	eng := New(WithIgnoreSymbol(true))
	if err := eng.AddWords([]Word{
		{Text: "测试", Type: "a"},
		{Text: "文本", Type: "b"},
	}); err != nil {
		t.Fatalf("AddWords() error = %v", err)
	}
	if err := eng.Build(); err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	matches := eng.FindAll("测试 文本")
	if len(matches) != 2 {
		t.Fatalf("FindAll() len = %d, want 2", len(matches))
	}
}

func TestEngine_FindAll_TraditionalVariant(t *testing.T) {
	eng := New(WithTraditional(true), WithIgnoreSymbol(true))
	if err := eng.AddWord(Word{Text: "测试", Type: "custom"}); err != nil {
		t.Fatalf("AddWord() error = %v", err)
	}
	if err := eng.Build(); err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if !eng.Contains("測試") {
		t.Fatal("Contains() = false, want true")
	}
}
