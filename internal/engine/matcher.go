package engine

import (
	ahocorasick "github.com/pgavlin/aho-corasick"
)

type pattern struct {
	word Word
	key  string
}

type matcher struct {
	cfg      Config
	patterns []pattern
	ac       ahocorasick.AhoCorasick
}

func newMatcher(cfg Config, words []Word) (*matcher, error) {
	patterns := make([]pattern, 0, len(words))
	keys := make([]string, 0, len(words))
	for _, w := range words {
		k := normalizeText(w.Text, cfg)
		if k == "" {
			continue
		}
		patterns = append(patterns, pattern{word: w, key: k})
		keys = append(keys, k)
	}
	builder := ahocorasick.NewAhoCorasickBuilder(ahocorasick.Opts{
		AsciiCaseInsensitive: false,
		MatchKind:            ahocorasick.StandardMatch,
		DFA:                  true,
	})
	ac := builder.Build(keys)
	return &matcher{cfg: cfg, patterns: patterns, ac: ac}, nil
}

func (m *matcher) findAll(text string) []Match {
	norm, mapping := normalizeWithMap(text, m.cfg)
	if norm == "" {
		return nil
	}
	matches := m.ac.FindAll(norm)
	if len(matches) == 0 {
		return nil
	}
	out := make([]Match, 0, len(matches))
	for _, mt := range matches {
		if mt.Pattern() < 0 || mt.Pattern() >= len(m.patterns) {
			continue
		}
		p := m.patterns[mt.Pattern()]
		startB := mt.Start()
		endB := mt.End()
		if startB < 0 || endB <= 0 || startB >= len(mapping) {
			continue
		}
		if endB-1 >= len(mapping) {
			endB = len(mapping)
		}
		startPos := mapping[startB]
		endPos := mapping[endB-1] + 1
		if startPos < 0 || endPos <= startPos {
			continue
		}
		out = append(out, Match{Word: p.word.Text, Type: p.word.Type, Text: runeSlice(text, startPos, endPos), StartPos: startPos, EndPos: endPos})
	}
	return out
}

func (m *matcher) replace(text, mask string) string {
	matches := m.findAll(text)
	if len(matches) == 0 {
		return text
	}
	runes := []rune(text)
	maskRune := []rune(mask)[0]
	for _, mt := range matches {
		for i := mt.StartPos; i < mt.EndPos && i < len(runes); i++ {
			runes[i] = maskRune
		}
	}
	return string(runes)
}

func runeSlice(text string, start, end int) string {
	r := []rune(text)
	if start < 0 {
		start = 0
	}
	if end > len(r) {
		end = len(r)
	}
	if start >= end {
		return ""
	}
	return string(r[start:end])
}
