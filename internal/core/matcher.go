package core

import (
	"fmt"
	"sort"
	"strings"

	ahocorasick "github.com/pgavlin/aho-corasick"
)

type pattern struct {
	word Word
	key  string
}

type matcher struct {
	cfg        Config
	normalizer normalizer
	patterns   []pattern
	ac         ahocorasick.AhoCorasick
}

func newMatcher(cfg Config, words []Word) (*matcher, error) {
	normalizer, err := newNormalizer(cfg)
	if err != nil {
		return nil, err
	}
	patterns := make([]pattern, 0, len(words))
	keys := make([]string, 0, len(words))
	for _, w := range words {
		k := normalizer.normalize(w.Text).text
		if k == "" {
			return nil, fmt.Errorf("word %q normalizes to empty", w.Text)
		}
		patterns = append(patterns, pattern{word: w, key: k})
		keys = append(keys, k)
	}
	builder := ahocorasick.NewAhoCorasickBuilder(ahocorasick.Opts{
		AsciiCaseInsensitive: false,
		MatchKind:            ahocorasick.LeftMostLongestMatch,
		DFA:                  true,
	})
	ac := builder.Build(keys)
	return &matcher{cfg: cfg, normalizer: normalizer, patterns: patterns, ac: ac}, nil
}

func (m *matcher) findAll(text string) []Match {
	normalized := m.normalizer.normalize(text)
	if normalized.text == "" {
		return nil
	}
	matches := m.ac.FindAll(normalized.text)
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
		if startB < 0 || endB <= 0 || startB >= len(normalized.mapping) {
			continue
		}
		if endB-1 >= len(normalized.mapping) {
			endB = len(normalized.mapping)
		}
		startPos := normalized.mapping[startB]
		endPos := normalized.mapping[endB-1] + 1
		if startPos < 0 || endPos <= startPos {
			continue
		}
		out = append(out, Match{Word: p.word.Text, Type: p.word.Type, Text: runeSlice(text, startPos, endPos), StartPos: startPos, EndPos: endPos})
	}
	return selectNonOverlappingMatches(out)
}

func (m *matcher) contains(text string) bool {
	normalized := m.normalizer.normalizeTextOnly(text)
	if normalized == "" {
		return false
	}
	return m.ac.Iter(normalized).Next() != nil
}

func (m *matcher) replace(text, mask string) string {
	matches := m.findAll(text)
	if len(matches) == 0 {
		return text
	}
	runes := []rune(text)
	var b strings.Builder
	b.Grow(len(text))
	pos := 0
	for _, mt := range matches {
		if mt.StartPos > pos {
			b.WriteString(string(runes[pos:mt.StartPos]))
		}
		end := mt.EndPos
		if end > len(runes) {
			end = len(runes)
		}
		for i := mt.StartPos; i < end; i++ {
			b.WriteString(mask)
		}
		pos = end
	}
	if pos < len(runes) {
		b.WriteString(string(runes[pos:]))
	}
	return b.String()
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

func selectNonOverlappingMatches(matches []Match) []Match {
	if len(matches) <= 1 {
		return matches
	}
	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].StartPos != matches[j].StartPos {
			return matches[i].StartPos < matches[j].StartPos
		}
		iLen := matches[i].EndPos - matches[i].StartPos
		jLen := matches[j].EndPos - matches[j].StartPos
		if iLen != jLen {
			return iLen > jLen
		}
		return matches[i].Word < matches[j].Word
	})

	selected := make([]Match, 0, len(matches))
	coveredEnd := -1
	for _, mt := range matches {
		if mt.StartPos < coveredEnd {
			continue
		}
		selected = append(selected, mt)
		if mt.EndPos > coveredEnd {
			coveredEnd = mt.EndPos
		}
	}
	return selected
}
