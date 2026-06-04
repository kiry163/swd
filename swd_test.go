package swd

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEngine_FindAll_ReturnsRunePositions(t *testing.T) {
	eng := New(Options{
		IgnoreSymbol: true,
		IgnoreWidth:  true,
		IgnoreCase:   true,
		Traditional:  true,
		SimilarChar:  true,
	})

	if err := eng.AddWord(Word{Text: "测试", Type: "custom"}); err != nil {
		t.Fatalf("AddWord() error = %v", err)
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

func TestEngine_Load_LoadsWordsFromMultipleSources(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "words.txt")
	if err := os.WriteFile(file, []byte("诈骗,crime\n赌博,crime\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	eng := New(Options{})
	err := eng.Load(
		context.Background(),
		NewFileLoader(file),
		NewMemoryLoader([]Word{{Text: "敏感词", Type: "custom"}}),
		NewReaderLoader(strings.NewReader("违规,policy\n")),
	)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !eng.Contains("这里有诈骗信息") {
		t.Fatal("Contains(file word) = false, want true")
	}
	if !eng.Contains("这里有敏感词") {
		t.Fatal("Contains(memory word) = false, want true")
	}
	if !eng.Contains("这里有违规内容") {
		t.Fatal("Contains(reader word) = false, want true")
	}
}

func TestEngine_Replace(t *testing.T) {
	eng := New(Options{})
	if err := eng.AddWords([]Word{{Text: "坏词", Type: "custom"}}); err != nil {
		t.Fatalf("AddWords() error = %v", err)
	}

	got := eng.Replace("这是坏词", "*")
	if got != "这是**" {
		t.Fatalf("Replace() = %q, want %q", got, "这是**")
	}
}

func TestEngine_Replace_UsesFullMaskForEachMatchedRune(t *testing.T) {
	eng := New(Options{})
	if err := eng.AddWords([]Word{{Text: "坏词", Type: "custom"}}); err != nil {
		t.Fatalf("AddWords() error = %v", err)
	}

	got := eng.Replace("这是坏词", "[x]")
	if got != "这是[x][x]" {
		t.Fatalf("Replace() = %q, want %q", got, "这是[x][x]")
	}
}

func TestEngine_FindAll_MultipleMatches(t *testing.T) {
	eng := New(Options{IgnoreSymbol: true})
	if err := eng.AddWords([]Word{
		{Text: "测试", Type: "a"},
		{Text: "文本", Type: "b"},
	}); err != nil {
		t.Fatalf("AddWords() error = %v", err)
	}

	matches := eng.FindAll("测试 文本")
	if len(matches) != 2 {
		t.Fatalf("FindAll() len = %d, want 2", len(matches))
	}
}

func TestEngine_FindAll_PrefersLongestMatchAtSameStart(t *testing.T) {
	eng := New(Options{})
	if err := eng.AddWords([]Word{
		{Text: "赌博", Type: "short"},
		{Text: "网络赌博", Type: "long"},
	}); err != nil {
		t.Fatalf("AddWords() error = %v", err)
	}

	matches := eng.FindAll("涉及网络赌博内容")
	if len(matches) != 1 {
		t.Fatalf("FindAll() len = %d, want 1: %#v", len(matches), matches)
	}
	if matches[0].Word != "网络赌博" || matches[0].Type != "long" {
		t.Fatalf("Match = %#v, want 网络赌博/long", matches[0])
	}
}

func TestEngine_FindAll_FiltersOverlappingShortMatches(t *testing.T) {
	eng := New(Options{})
	if err := eng.AddWords([]Word{
		{Text: "abcd", Type: "long"},
		{Text: "bc", Type: "inside"},
		{Text: "ef", Type: "next"},
	}); err != nil {
		t.Fatalf("AddWords() error = %v", err)
	}

	matches := eng.FindAll("xxabcdef")
	if len(matches) != 2 {
		t.Fatalf("FindAll() len = %d, want 2: %#v", len(matches), matches)
	}
	if matches[0].Word != "abcd" || matches[1].Word != "ef" {
		t.Fatalf("Words = %q/%q, want abcd/ef", matches[0].Word, matches[1].Word)
	}
	if matches[0].StartPos > matches[1].StartPos {
		t.Fatalf("matches are not sorted by position: %#v", matches)
	}
}

func TestEngine_AddWords_OverridesExistingWordByText(t *testing.T) {
	eng := New(Options{})
	if err := eng.AddWord(Word{Text: "重复", Type: "old"}); err != nil {
		t.Fatalf("AddWord(old) error = %v", err)
	}
	if err := eng.AddWord(Word{Text: "重复", Type: "new"}); err != nil {
		t.Fatalf("AddWord(new) error = %v", err)
	}

	matches := eng.FindAll("这里有重复词")
	if len(matches) != 1 {
		t.Fatalf("FindAll() len = %d, want 1", len(matches))
	}
	if matches[0].Type != "new" {
		t.Fatalf("Type = %q, want %q", matches[0].Type, "new")
	}
}

func TestEngine_FindAll_TraditionalVariant(t *testing.T) {
	eng := New(Options{Traditional: true, IgnoreSymbol: true})
	if err := eng.AddWord(Word{Text: "测试", Type: "custom"}); err != nil {
		t.Fatalf("AddWord() error = %v", err)
	}

	if !eng.Contains("測試") {
		t.Fatal("Contains() = false, want true")
	}
}

func TestEngine_FindAll_TraditionalPhraseUsesOpenCC(t *testing.T) {
	eng := New(Options{Traditional: true})
	if err := eng.AddWord(Word{Text: "软体", Type: "custom"}); err != nil {
		t.Fatalf("AddWord() error = %v", err)
	}

	matches := eng.FindAll("這個軟體很好")
	if len(matches) != 1 {
		t.Fatalf("FindAll() len = %d, want 1: %#v", len(matches), matches)
	}
	if matches[0].Text != "軟體" || matches[0].StartPos != 2 || matches[0].EndPos != 4 {
		t.Fatalf("Match = %#v, want text 軟體 at 2..4", matches[0])
	}
}

func TestEngine_FindAll_IgnoreWidthMatchesHalfwidthKatakanaMarks(t *testing.T) {
	eng := New(Options{IgnoreWidth: true})
	if err := eng.AddWord(Word{Text: "ガス", Type: "custom"}); err != nil {
		t.Fatalf("AddWord() error = %v", err)
	}

	matches := eng.FindAll("これはｶﾞｽです")
	if len(matches) != 1 {
		t.Fatalf("FindAll() len = %d, want 1: %#v", len(matches), matches)
	}
	if matches[0].Text != "ｶﾞｽ" || matches[0].StartPos != 3 || matches[0].EndPos != 6 {
		t.Fatalf("Match = %#v, want text ｶﾞｽ at 3..6", matches[0])
	}
}

func TestEngine_Contains_IgnoreSymbolSkipsZeroWidthCharacters(t *testing.T) {
	eng := New(Options{IgnoreSymbol: true})
	if err := eng.AddWord(Word{Text: "敏感词", Type: "custom"}); err != nil {
		t.Fatalf("AddWord() error = %v", err)
	}

	if !eng.Contains("敏\u200b感\u200d词") {
		t.Fatal("Contains() = false, want true")
	}
}

func TestEngine_Contains_SimilarCharMatchesCommonLeetspeak(t *testing.T) {
	eng := New(Options{SimilarChar: true, IgnoreCase: true})
	if err := eng.AddWords([]Word{
		{Text: "bad", Type: "custom"},
		{Text: "test", Type: "custom"},
	}); err != nil {
		t.Fatalf("AddWords() error = %v", err)
	}

	if !eng.Contains("b@d") {
		t.Fatal("Contains(b@d) = false, want true")
	}
	if !eng.Contains("t3st") {
		t.Fatal("Contains(t3st) = false, want true")
	}
}

func TestEngine_Contains_UsesSameNormalizationAsFind(t *testing.T) {
	eng := New(Options{
		IgnoreSymbol: true,
		IgnoreWidth:  true,
		IgnoreCase:   true,
		Traditional:  true,
		SimilarChar:  true,
	})
	if err := eng.AddWord(Word{Text: "测试", Type: "custom"}); err != nil {
		t.Fatalf("AddWord() error = %v", err)
	}

	text := "这是一段 測@試 的文本"
	if !eng.Contains(text) {
		t.Fatal("Contains() = false, want true")
	}
	if eng.Find(text) == nil {
		t.Fatal("Find() = nil, want match")
	}
}

func TestEngine_RemoveWord_RebuildsMatcher(t *testing.T) {
	eng := New(Options{})
	if err := eng.AddWords([]Word{
		{Text: "保留", Type: "a"},
		{Text: "移除", Type: "b"},
		{Text: "也移除", Type: "b"},
	}); err != nil {
		t.Fatalf("AddWords() error = %v", err)
	}

	if err := eng.RemoveWords([]string{"移除", "也移除"}); err != nil {
		t.Fatalf("RemoveWords() error = %v", err)
	}

	if eng.Contains("这里有移除") {
		t.Fatal("Contains(removed word) = true, want false")
	}
	if eng.Contains("这里也移除") {
		t.Fatal("Contains(second removed word) = true, want false")
	}
	if !eng.Contains("这里有保留") {
		t.Fatal("Contains(kept word) = false, want true")
	}
}

func TestEngine_Clear_RebuildsEmptyMatcher(t *testing.T) {
	eng := New(Options{})
	if err := eng.AddWord(Word{Text: "敏感", Type: "custom"}); err != nil {
		t.Fatalf("AddWord() error = %v", err)
	}

	if err := eng.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	if eng.Contains("这里有敏感内容") {
		t.Fatal("Contains() = true after Clear, want false")
	}
}

func TestEngine_LoadError_KeepsExistingWords(t *testing.T) {
	eng := New(Options{})
	if err := eng.AddWord(Word{Text: "旧词", Type: "old"}); err != nil {
		t.Fatalf("AddWord() error = %v", err)
	}

	err := eng.Load(context.Background(), NewMemoryLoader([]Word{
		{Text: "新词", Type: "new"},
		{Text: " ", Type: "bad"},
	}))
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}

	if !eng.Contains("这里有旧词") {
		t.Fatal("Contains(old word) = false after failed Load, want true")
	}
	if eng.Contains("这里有新词") {
		t.Fatal("Contains(new word) = true after failed Load, want false")
	}
}
