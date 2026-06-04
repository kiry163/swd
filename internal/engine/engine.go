package engine

import (
	"context"
	"fmt"
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

func New(opts ...Option) *Engine {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
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
