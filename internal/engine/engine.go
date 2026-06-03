package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Engine struct {
	cfg   Config
	words []Word
	m     *matcher
}

func New(opts ...Option) *Engine {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Engine{cfg: cfg}
}

func (e *Engine) AddWord(w Word) error {
	if strings.TrimSpace(w.Text) == "" {
		return fmt.Errorf("word text is empty")
	}
	e.words = append(e.words, w)
	return nil
}

func (e *Engine) AddWords(words []Word) error {
	for _, w := range words {
		if err := e.AddWord(w); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) AddFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ",", 2)
		w := Word{Text: strings.TrimSpace(parts[0])}
		if len(parts) > 1 {
			w.Type = strings.TrimSpace(parts[1])
		}
		if err := e.AddWord(w); err != nil {
			return err
		}
	}
	return s.Err()
}

func (e *Engine) Build() error {
	m, err := newMatcher(e.cfg, e.words)
	if err != nil {
		return err
	}
	e.m = m
	return nil
}

func (e *Engine) FindAll(text string) []Match {
	if e.m == nil {
		return nil
	}
	return e.m.findAll(text)
}

func (e *Engine) Find(text string) *Match {
	matches := e.FindAll(text)
	if len(matches) == 0 {
		return nil
	}
	return &matches[0]
}

func (e *Engine) Contains(text string) bool {
	return len(e.FindAll(text)) > 0
}

func (e *Engine) Replace(text, mask string) string {
	if e.m == nil || mask == "" {
		return text
	}
	return e.m.replace(text, mask)
}
