package swd

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestEngine_ConcurrentReadAndWrite(t *testing.T) {
	eng := New()
	if err := eng.AddWords([]Word{
		{Text: "基础词", Type: "base"},
		{Text: "保留词", Type: "base"},
	}); err != nil {
		t.Fatalf("AddWords() error = %v", err)
	}

	const readers = 8
	const iterations = 100

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-start
			for j := 0; j < iterations; j++ {
				text := fmt.Sprintf("reader-%d-%d 基础词 保留词 动态词-%d", id, j, j%10)
				_ = eng.Contains(text)
				_ = eng.FindAll(text)
				_ = eng.Replace(text, "*")
			}
		}(i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < iterations; i++ {
			word := fmt.Sprintf("动态词-%d", i%10)
			if err := eng.AddWord(Word{Text: word, Type: "dynamic"}); err != nil {
				t.Errorf("AddWord(%q) error = %v", word, err)
				return
			}
			if i%3 == 0 {
				if err := eng.RemoveWord(word); err != nil {
					t.Errorf("RemoveWord(%q) error = %v", word, err)
					return
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < iterations/10; i++ {
			if err := eng.Clear(); err != nil {
				t.Errorf("Clear() error = %v", err)
				return
			}
			if err := eng.Load(context.Background(), NewMemoryLoader([]Word{
				{Text: "基础词", Type: "base"},
				{Text: "保留词", Type: "base"},
			})); err != nil {
				t.Errorf("Load() error = %v", err)
				return
			}
		}
	}()

	close(start)
	wg.Wait()

	if err := eng.AddWord(Word{Text: "最终词", Type: "final"}); err != nil {
		t.Fatalf("AddWord(final) error = %v", err)
	}
	if !eng.Contains("这里有最终词") {
		t.Fatal("Contains(final word) = false, want true")
	}
}
