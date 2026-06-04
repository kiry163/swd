package engine

import (
	"context"
	"strings"
	"testing"
)

func TestReaderLoader_LoadsWordsAndSkipsComments(t *testing.T) {
	loader := NewReaderLoader(strings.NewReader(`
# comment
诈骗,risk
赌博

  违规  ,  policy
`))

	words, err := loader.Load(context.Background())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(words) != 3 {
		t.Fatalf("len(words) = %d, want 3: %#v", len(words), words)
	}
	if words[0].Text != "诈骗" || words[0].Type != "risk" {
		t.Fatalf("words[0] = %#v, want 诈骗/risk", words[0])
	}
	if words[1].Text != "赌博" || words[1].Type != "" {
		t.Fatalf("words[1] = %#v, want 赌博", words[1])
	}
	if words[2].Text != "违规" || words[2].Type != "policy" {
		t.Fatalf("words[2] = %#v, want 违规/policy", words[2])
	}
}

func TestReaderLoader_ReturnsLineNumberForEmptyWord(t *testing.T) {
	loader := NewReaderLoader(strings.NewReader("正常,ok\n ,bad\n"))

	_, err := loader.Load(context.Background())
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Fatalf("Load() error = %q, want line number", err.Error())
	}
	if !strings.Contains(err.Error(), "word text is empty") {
		t.Fatalf("Load() error = %q, want empty word message", err.Error())
	}
}

func TestReaderLoader_AllowsLongLines(t *testing.T) {
	longWord := strings.Repeat("长", 70*1024)
	loader := NewReaderLoader(strings.NewReader(longWord + ",long\n"))

	words, err := loader.Load(context.Background())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(words) != 1 {
		t.Fatalf("len(words) = %d, want 1", len(words))
	}
	if words[0].Text != longWord {
		t.Fatalf("len(words[0].Text) = %d, want %d", len(words[0].Text), len(longWord))
	}
	if words[0].Type != "long" {
		t.Fatalf("Type = %q, want long", words[0].Type)
	}
}
