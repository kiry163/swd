package variant

import "strings"

// SimilarProcessor handles similar character detection
type SimilarProcessor struct {
	similarMap map[rune][]rune
}

// NewSimilarProcessor creates a new similar character processor
func NewSimilarProcessor() *SimilarProcessor {
	return &SimilarProcessor{
		similarMap: buildSimilarMap(),
	}
}

// Process normalizes similar characters to their base form
func (p *SimilarProcessor) Process(text string) string {
	var builder strings.Builder
	builder.Grow(len(text))

	for _, r := range text {
		// Find the base character for similar variants
		base := p.findBase(r)
		builder.WriteRune(base)
	}

	return builder.String()
}

// Name returns the processor name
func (p *SimilarProcessor) Name() string {
	return "similar"
}

// findBase finds the base character for a given rune
func (p *SimilarProcessor) findBase(r rune) rune {
	// Check if this character is in any similarity group
	for base, similars := range p.similarMap {
		for _, similar := range similars {
			if r == similar {
				return base
			}
		}
		if r == base {
			return base
		}
	}
	return r
}

// AddSimilarRule adds a custom similar character rule
func (p *SimilarProcessor) AddSimilarRule(base rune, similars ...rune) {
	if p.similarMap == nil {
		p.similarMap = make(map[rune][]rune)
	}
	p.similarMap[base] = append(p.similarMap[base], similars...)
}

// buildSimilarMap builds a map of similar characters
// Base character -> similar variants
func buildSimilarMap() map[rune][]rune {
	return map[rune][]rune{
		// Original examples
		'傻': {'煞', '𫓧', '儍'},
		'比': {'笔', '币', '苝', '怭'},
		'测': {'側', '厕'},
		'草': {'艹', '屮', '荡'},

		// Numbers and letters
		'0': {'O', 'o', '〇', '零', '○'},
		'1': {'l', 'I', 'i', '|', '一'},
		'2': {'二', '貳'},
		'3': {'三', '叁'},
		'4': {'四', '肆'},
		'5': {'五', '伍'},
		'6': {'六', '陸'},
		'7': {'七', '柒'},
		'8': {'八', '捌'},
		'9': {'九', '玖'},

		// Common similar looking characters
		'a': {'@', 'α'},
		'o': {'0', 'O', '〇'},
		'i': {'l', '1', 'I', '|'},
		's': {'$', '§'},
		'e': {'3', 'ε'},
		'g': {'9', 'q'},

		// Chinese character similarities
		'日': {'曰', '目'},
		'土': {'士', '壬'},
		'刀': {'力', '刃'},
		'人': {'入', '八'},
		'千': {'干', '于'},
		'未': {'末', '朱'},
		'己': {'已', '巳'},
		'戊': {'戌', '戍'},
		'大': {'太', '犬'},
		'天': {'夫', '无'},

		// Sensitive word evasion characters
		'政': {'政', '正', '征'},
		'色': {'色', '涩', '瑟'},
		'赌': {'堵', '睹'},
		'毒': {'独', '督'},
		'黄': {'皇', '煌'},
		'暴': {'爆', '曝'},
	}
}

// IsSimilar checks if two characters are similar
func IsSimilar(r1, r2 rune) bool {
	processor := NewSimilarProcessor()
	base1 := processor.findBase(r1)
	base2 := processor.findBase(r2)
	return base1 == base2
}
