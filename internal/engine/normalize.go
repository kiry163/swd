package engine

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var traditionalMap = map[rune]rune{
	'測': '测', '試': '试', '詞': '词', '政': '政', '黃': '黄', '賭': '赌', '毒': '毒', '網': '网', '絡': '络',
}

var similarMap = map[rune]rune{
	'@': 0,
	'０': '0', '１': '1', '２': '2', '３': '3', '４': '4', '５': '5', '６': '6', '７': '7', '８': '8', '９': '9',
	'〇': '0', '零': '0', '一': '1', '二': '2', '三': '3', '四': '4', '五': '5', '六': '6', '七': '7', '八': '8', '九': '9',
}

var pinyinMap = map[rune]string{
	'测': "ce", '試': "shi", '试': "shi", '敏': "min", '感': "gan", '词': "ci", '政': "zheng", '毒': "du", '赌': "du", '博': "bo",
}

func normalizeText(text string, cfg Config) string {
	n, _ := normalizeWithMap(text, cfg)
	return n
}

func normalizeWithMap(text string, cfg Config) (string, []int) {
	n := newNormalizer(cfg)
	out := n.normalize(text)
	return out.text, out.mapping
}

type normalizedText struct {
	text    string
	mapping []int
}

type normalizer struct {
	cfg Config
}

func newNormalizer(cfg Config) normalizer {
	return normalizer{cfg: cfg}
}

func (n normalizer) normalize(text string) normalizedText {
	var b strings.Builder
	mapping := make([]int, 0, len(text))
	runes := []rune(text)
	for i, r := range runes {
		r = n.normalizeRune(r)
		if r == 0 {
			continue
		}
		if n.cfg.IgnoreSymbol && isSymbolRune(r) {
			continue
		}
		if n.cfg.Pinyin {
			if py, ok := pinyinMap[r]; ok {
				appendMappedString(&b, &mapping, py, i)
				continue
			}
		}
		appendMappedRune(&b, &mapping, r, i)
	}
	return normalizedText{text: b.String(), mapping: mapping}
}

func (n normalizer) normalizeRune(r rune) rune {
	if n.cfg.IgnoreCase {
		r = unicode.ToLower(r)
	}
	if n.cfg.IgnoreWidth {
		r = toHalfWidth(r)
	}
	if n.cfg.Traditional {
		if v, ok := traditionalMap[r]; ok {
			r = v
		}
	}
	if n.cfg.SimilarChar || n.cfg.Homophone {
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
	return newNormalizer(cfg).normalizeRune(r)
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

func toHalfWidth(r rune) rune {
	if r == 12288 {
		return 32
	}
	if r >= 65281 && r <= 65374 {
		return r - 65248
	}
	return r
}

func isSymbolRune(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r) || unicode.IsSpace(r)
}
