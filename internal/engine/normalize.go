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
	var b strings.Builder
	mapping := make([]int, 0, len(text))
	runes := []rune(text)
	for i, r := range runes {
		r = normalizeRune(r, cfg)
		if r == 0 {
			continue
		}
		if cfg.IgnoreSymbol && isSymbolRune(r) {
			continue
		}
		if cfg.Pinyin {
			if py, ok := pinyinMap[r]; ok {
				for _, pr := range py {
					var buf [utf8.UTFMax]byte
					n := utf8.EncodeRune(buf[:], pr)
					b.Write(buf[:n])
					for j := 0; j < n; j++ {
						mapping = append(mapping, i)
					}
				}
				continue
			}
		}
		var buf [utf8.UTFMax]byte
		n := utf8.EncodeRune(buf[:], r)
		b.Write(buf[:n])
		for j := 0; j < n; j++ {
			mapping = append(mapping, i)
		}
	}
	return b.String(), mapping
}

func normalizeRune(r rune, cfg Config) rune {
	if cfg.IgnoreCase {
		r = unicode.ToLower(r)
	}
	if cfg.IgnoreWidth {
		r = toHalfWidth(r)
	}
	if cfg.Traditional {
		if v, ok := traditionalMap[r]; ok {
			r = v
		}
	}
	if cfg.SimilarChar || cfg.Homophone {
		if v, ok := similarMap[r]; ok {
			if v == 0 {
				return 0
			}
			r = v
		}
	}
	return r
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
