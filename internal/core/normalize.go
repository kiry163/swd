package core

import (
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/longbridgeapp/opencc"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
)

type mappedRune struct {
	r     rune
	index int
}

var (
	t2sOnce      sync.Once
	t2sConverter *opencc.OpenCC
	t2sErr       error
)

var similarMap = map[rune]rune{
	'@': 'a', '4': 'a',
	'3': 'e',
	'1': 'i', '!': 'i', '|': 'i',
	'0': 'o',
	'5': 's', '$': 's',
	'7': 't', '+': 't',
	'０': 'o', '１': 'i', '３': 'e', '４': 'a', '５': 's', '７': 't',
	'〇': 'o',
}

func normalizeText(text string, cfg Config) string {
	n, _ := normalizeWithMap(text, cfg)
	return n
}

func normalizeWithMap(text string, cfg Config) (string, []int) {
	n, err := newNormalizer(cfg)
	if err != nil {
		return text, nil
	}
	out := n.normalize(text)
	return out.text, out.mapping
}

type normalizedText struct {
	text    string
	mapping []int
}

type normalizer struct {
	cfg         Config
	traditional *opencc.OpenCC
}

func newNormalizer(cfg Config) (normalizer, error) {
	n := normalizer{cfg: cfg}
	if cfg.Traditional {
		cc, err := traditionalConverter()
		if err != nil {
			return normalizer{}, err
		}
		n.traditional = cc
	}
	return n, nil
}

func (n normalizer) normalize(text string) normalizedText {
	var b strings.Builder
	mapping := make([]int, 0, len(text))
	if n.traditional == nil && !n.cfg.IgnoreWidth {
		for i, r := range []rune(text) {
			n.appendNormalizedRune(&b, &mapping, r, i)
		}
		return normalizedText{text: b.String(), mapping: mapping}
	}

	runes := make([]mappedRune, 0, len(text))
	for i, r := range []rune(text) {
		runes = append(runes, mappedRune{r: r, index: i})
	}
	if n.traditional != nil {
		runes = n.convertTraditional(runes)
	}
	if n.cfg.IgnoreWidth {
		runes = n.foldWidth(runes)
	}
	for _, mr := range runes {
		n.appendNormalizedRune(&b, &mapping, mr.r, mr.index)
	}
	return normalizedText{text: b.String(), mapping: mapping}
}

func (n normalizer) normalizeTextOnly(text string) string {
	var b strings.Builder
	if n.traditional == nil && !n.cfg.IgnoreWidth {
		for _, r := range text {
			n.appendNormalizedRuneTextOnly(&b, r)
		}
		return b.String()
	}

	runes := make([]mappedRune, 0, len(text))
	for i, r := range []rune(text) {
		runes = append(runes, mappedRune{r: r, index: i})
	}
	if n.traditional != nil {
		runes = n.convertTraditional(runes)
	}
	if n.cfg.IgnoreWidth {
		runes = n.foldWidth(runes)
	}
	for _, mr := range runes {
		n.appendNormalizedRuneTextOnly(&b, mr.r)
	}
	return b.String()
}

func (n normalizer) appendNormalizedRune(b *strings.Builder, mapping *[]int, r rune, index int) {
	if n.cfg.IgnoreSymbol && isSymbolRune(r) {
		return
	}
	r = n.normalizeRune(r)
	if r == 0 {
		return
	}
	appendMappedRune(b, mapping, r, index)
}

func (n normalizer) appendNormalizedRuneTextOnly(b *strings.Builder, r rune) {
	if n.cfg.IgnoreSymbol && isSymbolRune(r) {
		return
	}
	r = n.normalizeRune(r)
	if r == 0 {
		return
	}
	b.WriteRune(r)
}

func (n normalizer) normalizeRune(r rune) rune {
	if n.cfg.IgnoreCase {
		r = unicode.ToLower(r)
	}
	if n.cfg.SimilarChar {
		if v, ok := similarMap[r]; ok {
			if v == 0 {
				return 0
			}
			r = v
		}
	}
	return r
}

func normalizeRune(r rune, cfg Config) rune {
	n, err := newNormalizer(cfg)
	if err != nil {
		return r
	}
	return n.normalizeRune(r)
}

func traditionalConverter() (*opencc.OpenCC, error) {
	t2sOnce.Do(func() {
		t2sConverter, t2sErr = opencc.New("t2s")
	})
	return t2sConverter, t2sErr
}

func (n normalizer) convertTraditional(input []mappedRune) []mappedRune {
	out := input
	for _, group := range n.traditional.DictChains {
		converted := make([]mappedRune, 0, len(out))
		for i := 0; i < len(out); {
			offset := len(out)
			if offset > i+10 {
				offset = i + 10
			}
			segment := mappedRunesString(out[i:offset])
			matchedKey := ""
			replacement := ""
			max := 0
			for _, dict := range group.Dicts {
				ret, err := dict.PrefixMatch(segment)
				if err != nil {
					return input
				}
				for k, v := range ret {
					if len(k) > max && len(v) > 0 {
						max = len(k)
						matchedKey = k
						replacement = v[0]
					}
				}
			}
			if max == 0 {
				converted = append(converted, out[i])
				i++
				continue
			}
			inLen := len([]rune(matchedKey))
			converted = append(converted, mapReplacementRunes([]rune(replacement), out[i:i+inLen])...)
			i += inLen
		}
		out = converted
	}
	return out
}

func (n normalizer) foldWidth(input []mappedRune) []mappedRune {
	out := make([]mappedRune, 0, len(input))
	for _, mr := range input {
		folded, _, err := transform.String(transform.Chain(width.Fold, norm.NFD), string(mr.r))
		if err != nil || folded == "" {
			out = append(out, mr)
			continue
		}
		for _, r := range folded {
			out = append(out, mappedRune{r: r, index: mr.index})
		}
	}
	return out
}

func mappedRunesString(runes []mappedRune) string {
	var b strings.Builder
	for _, r := range runes {
		b.WriteRune(r.r)
	}
	return b.String()
}

func mapReplacementRunes(replacement []rune, source []mappedRune) []mappedRune {
	if len(replacement) == 0 {
		return nil
	}
	out := make([]mappedRune, 0, len(replacement))
	for i, r := range replacement {
		sourceIndex := 0
		if len(replacement) == len(source) {
			sourceIndex = i
		} else if len(source) > 1 {
			sourceIndex = i * len(source) / len(replacement)
			if sourceIndex >= len(source) {
				sourceIndex = len(source) - 1
			}
		}
		out = append(out, mappedRune{r: r, index: source[sourceIndex].index})
	}
	return out
}

func appendMappedString(b *strings.Builder, mapping *[]int, text string, runeIndex int) {
	for _, r := range text {
		appendMappedRune(b, mapping, r, runeIndex)
	}
}

func appendMappedRune(b *strings.Builder, mapping *[]int, r rune, runeIndex int) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	b.Write(buf[:n])
	for j := 0; j < n; j++ {
		*mapping = append(*mapping, runeIndex)
	}
}

func isSymbolRune(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r) || unicode.IsSpace(r) || unicode.IsControl(r) || unicode.Is(unicode.Cf, r)
}
