package core

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

const maxLoaderLineBytes = 1024 * 1024

type Loader interface {
	Load(context.Context) ([]Word, error)
}

type memoryLoader struct {
	words []Word
}

func NewMemoryLoader(words []Word) Loader {
	cp := append([]Word(nil), words...)
	return memoryLoader{words: cp}
}

func (l memoryLoader) Load(context.Context) ([]Word, error) {
	return append([]Word(nil), l.words...), nil
}

type fileLoader struct {
	path string
}

func NewFileLoader(path string) Loader {
	return fileLoader{path: path}
}

func (l fileLoader) Load(ctx context.Context) ([]Word, error) {
	f, err := os.Open(l.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return loadWordsFromReader(ctx, f)
}

type readerLoader struct {
	r io.Reader
}

func NewReaderLoader(r io.Reader) Loader {
	return readerLoader{r: r}
}

func (l readerLoader) Load(ctx context.Context) ([]Word, error) {
	return loadWordsFromReader(ctx, l.r)
}

func loadWordsFromReader(ctx context.Context, r io.Reader) ([]Word, error) {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0, 64*1024), maxLoaderLineBytes)
	words := make([]Word, 0)
	lineNo := 0
	for s.Scan() {
		lineNo++
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ",", 2)
		w := Word{Text: strings.TrimSpace(parts[0])}
		if w.Text == "" {
			return nil, fmt.Errorf("line %d: word text is empty", lineNo)
		}
		if len(parts) > 1 {
			w.Type = strings.TrimSpace(parts[1])
		}
		words = append(words, w)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return words, nil
}
