package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Engine struct {
	writeMu sync.Mutex
	mu      sync.RWMutex
	cfg     Config
	words   []Word
	m       *matcher
}

func New(options Options) *Engine {
	cfg := defaultConfig()
	cfg.IgnoreSymbol = options.IgnoreSymbol
	cfg.IgnoreWidth = options.IgnoreWidth
	cfg.IgnoreCase = options.IgnoreCase
	cfg.Traditional = options.Traditional
	cfg.SimilarChar = options.SimilarChar
	return &Engine{cfg: cfg}
}

func (e *Engine) AddWord(w Word) error {
	return e.AddWords([]Word{w})
}

func (e *Engine) AddWords(words []Word) error {
	if err := validateWords(words); err != nil {
		return err
	}
	e.writeMu.Lock()
	defer e.writeMu.Unlock()

	e.mu.RLock()
	next := appendOrOverrideWords(e.words, words)
	e.mu.RUnlock()
	return e.replaceWords(next)
}

func (e *Engine) Load(ctx context.Context, loaders ...Loader) error {
	loaded := make([]Word, 0)
	for _, loader := range loaders {
		if loader == nil {
			return fmt.Errorf("loader is nil")
		}
		words, err := loader.Load(ctx)
		if err != nil {
			return err
		}
		loaded = append(loaded, words...)
	}
	return e.AddWords(loaded)
}

func (e *Engine) ReplaceWords(words []Word) error {
	e.writeMu.Lock()
	defer e.writeMu.Unlock()
	return e.replaceWords(appendOrOverrideWords(nil, words))
}

func (e *Engine) Words() []Word {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return append([]Word(nil), e.words...)
}

func (e *Engine) Export(ctx context.Context, w io.Writer) error {
	if w == nil {
		return fmt.Errorf("writer is nil")
	}
	words := e.Words()
	if err := validateExportWords(words); err != nil {
		return err
	}
	return exportWords(ctx, words, w)
}

func (e *Engine) ExportFile(ctx context.Context, path string) error {
	words := e.Words()
	if err := validateExportWords(words); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	dir := filepath.Dir(path)
	f, err := os.CreateTemp(dir, ".swd-export-*")
	if err != nil {
		return err
	}
	tmpPath := f.Name()
	removeTmp := true
	defer func() {
		if removeTmp {
			_ = os.Remove(tmpPath)
		}
	}()

	err = exportWords(ctx, words, f)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}
	removeTmp = false
	return nil
}

func (e *Engine) RemoveWord(text string) error {
	return e.RemoveWords([]string{text})
}

func (e *Engine) RemoveWords(texts []string) error {
	remove := make(map[string]struct{}, len(texts))
	for _, text := range texts {
		text = strings.TrimSpace(text)
		if text == "" {
			return fmt.Errorf("word text is empty")
		}
		remove[text] = struct{}{}
	}

	e.writeMu.Lock()
	defer e.writeMu.Unlock()

	e.mu.RLock()
	next := make([]Word, 0, len(e.words))
	for _, w := range e.words {
		if _, ok := remove[w.Text]; ok {
			continue
		}
		next = append(next, w)
	}
	e.mu.RUnlock()
	return e.replaceWords(next)
}

func (e *Engine) Clear() error {
	e.writeMu.Lock()
	defer e.writeMu.Unlock()
	return e.replaceWords(nil)
}

func (e *Engine) FindAll(text string) []Match {
	e.mu.RLock()
	m := e.m
	e.mu.RUnlock()
	if m == nil {
		return nil
	}
	return m.findAll(text)
}

func (e *Engine) Find(text string) *Match {
	matches := e.FindAll(text)
	if len(matches) == 0 {
		return nil
	}
	return &matches[0]
}

func (e *Engine) Contains(text string) bool {
	e.mu.RLock()
	m := e.m
	e.mu.RUnlock()
	if m == nil {
		return false
	}
	return m.contains(text)
}

func (e *Engine) Replace(text, mask string) string {
	e.mu.RLock()
	m := e.m
	e.mu.RUnlock()
	if m == nil || mask == "" {
		return text
	}
	return m.replace(text, mask)
}

func (e *Engine) replaceWords(words []Word) error {
	if err := validateWords(words); err != nil {
		return err
	}
	next := append([]Word(nil), words...)
	m, err := newMatcher(e.cfg, next)
	if err != nil {
		return err
	}

	e.mu.Lock()
	e.words = next
	e.m = m
	e.mu.Unlock()
	return nil
}

func validateWords(words []Word) error {
	for _, w := range words {
		if strings.TrimSpace(w.Text) == "" {
			return fmt.Errorf("word text is empty")
		}
	}
	return nil
}

func appendOrOverrideWords(base, updates []Word) []Word {
	next := append([]Word(nil), base...)
	index := make(map[string]int, len(next)+len(updates))
	for i, w := range next {
		index[w.Text] = i
	}
	for _, w := range updates {
		if i, ok := index[w.Text]; ok {
			next[i] = w
			continue
		}
		index[w.Text] = len(next)
		next = append(next, w)
	}
	return next
}

func exportWords(ctx context.Context, words []Word, w io.Writer) error {
	for _, word := range words {
		if err := ctx.Err(); err != nil {
			return err
		}
		line := word.Text
		if word.Type != "" {
			line += "," + word.Type
		}
		if _, err := io.WriteString(w, line+"\n"); err != nil {
			return err
		}
	}
	return ctx.Err()
}

func validateExportWords(words []Word) error {
	for _, w := range words {
		if strings.ContainsAny(w.Text, ",\r\n") {
			return fmt.Errorf("word %q cannot be exported in simple format", w.Text)
		}
		if strings.ContainsAny(w.Type, ",\r\n") {
			return fmt.Errorf("word type %q cannot be exported in simple format", w.Type)
		}
	}
	return nil
}
