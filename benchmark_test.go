package swd

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

var benchmarkSink any

func BenchmarkLoadMemory(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		words := benchmarkWords(size)
		b.Run(fmt.Sprintf("words_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				eng := New()
				if err := eng.Load(context.Background(), NewMemoryLoader(words)); err != nil {
					b.Fatal(err)
				}
				benchmarkSink = eng
			}
		})
	}
}

func BenchmarkContains(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		eng := benchmarkEngine(b, size)
		text := benchmarkText(size, false)
		b.Run(fmt.Sprintf("words_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				benchmarkSink = eng.Contains(text)
			}
		})
	}
}

func BenchmarkFindAll(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		eng := benchmarkEngine(b, size)
		text := benchmarkText(size, true)
		b.Run(fmt.Sprintf("words_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				benchmarkSink = eng.FindAll(text)
			}
		})
	}
}

func BenchmarkReplace(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		eng := benchmarkEngine(b, size)
		text := benchmarkText(size, true)
		b.Run(fmt.Sprintf("words_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				benchmarkSink = eng.Replace(text, "*")
			}
		})
	}
}

func BenchmarkAddWordRuntimeUpdate(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		words := benchmarkWords(size)
		b.Run(fmt.Sprintf("base_words_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				eng := New()
				if err := eng.Load(context.Background(), NewMemoryLoader(words)); err != nil {
					b.Fatal(err)
				}
				if err := eng.AddWord(Word{Text: fmt.Sprintf("runtime-%d", i), Type: "runtime"}); err != nil {
					b.Fatal(err)
				}
				benchmarkSink = eng
			}
		})
	}
}

func benchmarkEngine(b *testing.B, size int) *Engine {
	b.Helper()
	eng := New()
	if err := eng.Load(context.Background(), NewMemoryLoader(benchmarkWords(size))); err != nil {
		b.Fatal(err)
	}
	return eng
}

func benchmarkWords(size int) []Word {
	words := make([]Word, size)
	for i := range words {
		words[i] = Word{
			Text: fmt.Sprintf("敏感词%06d", i),
			Type: fmt.Sprintf("type-%02d", i%10),
		}
	}
	return words
}

func benchmarkText(size int, manyMatches bool) string {
	var b strings.Builder
	b.WriteString("这是一段用于性能测试的普通文本。")
	if manyMatches {
		for i := 0; i < 20; i++ {
			if i > 0 {
				b.WriteString("，")
			}
			b.WriteString(fmt.Sprintf("敏感词%06d", (i*997)%size))
		}
	} else {
		b.WriteString(fmt.Sprintf("这里包含敏感词%06d。", size/2))
	}
	b.WriteString("结尾还有一些普通内容。")
	return b.String()
}
