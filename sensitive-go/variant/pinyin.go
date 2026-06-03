package variant

import (
	"strings"
)

// PinyinProcessor handles pinyin variant detection
type PinyinProcessor struct {
	// pinyinMap maps characters to their pinyin representation
	pinyinMap map[rune]string
}

// NewPinyinProcessor creates a new pinyin processor
func NewPinyinProcessor() *PinyinProcessor {
	return &PinyinProcessor{
		pinyinMap: buildBasicPinyinMap(),
	}
}

// Process converts Chinese characters to pinyin for matching
func (p *PinyinProcessor) Process(text string) string {
	var builder strings.Builder
	builder.Grow(len(text))

	for _, r := range text {
		if py, exists := p.pinyinMap[r]; exists {
			builder.WriteString(py)
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// Name returns the processor name
func (p *PinyinProcessor) Name() string {
	return "pinyin"
}

// buildBasicPinyinMap builds a basic pinyin mapping
// In a production implementation, this would use a complete pinyin library
func buildBasicPinyinMap() map[rune]string {
	// Extended basic pinyin mappings for common characters
	return map[rune]string{
		// Common words
		'傻': "sha", '比': "bi", '测': "ce", '试': "shi", '敏': "min", '感': "gan", '词': "ci",

		// Numbers
		'一': "yi", '二': "er", '三': "san", '四': "si", '五': "wu",
		'六': "liu", '七': "qi", '八': "ba", '九': "jiu", '十': "shi",

		// Common characters (100 most frequent)
		'的': "de", '了': "le", '是': "shi", '我': "wo", '不': "bu",
		'在': "zai", '人': "ren", '有': "you", '他': "ta", '这': "zhe",
		'个': "ge", '们': "men", '中': "zhong", '来': "lai", '上': "shang",
		'大': "da", '为': "wei", '和': "he", '国': "guo", '地': "di",
		'到': "dao", '以': "yi", '说': "shuo", '时': "shi", '要': "yao",
		'就': "jiu", '出': "chu", '会': "hui", '可': "ke", '也': "ye",
		'你': "ni", '对': "dui", '生': "sheng", '能': "neng", '而': "er",
		'子': "zi", '那': "na", '得': "de", '于': "yu", '着': "zhe",
		'下': "xia", '自': "zi", '之': "zhi", '年': "nian", '过': "guo",
		'发': "fa", '后': "hou", '作': "zuo", '里': "li", '用': "yong",
		'道': "dao", '行': "xing", '所': "suo", '然': "ran", '家': "jia",
		'种': "zhong", '事': "shi", '成': "cheng", '方': "fang", '多': "duo",
		'经': "jing", '么': "me", '去': "qu", '法': "fa", '学': "xue",
		'如': "ru", '她': "ta", '只': "zhi", '现': "xian", '当': "dang",
		'样': "yang", '看': "kan", '文': "wen", '无': "wu", '开': "kai",
		'手': "shou", '主': "zhu",
		'又': "you", '高': "gao", '小': "xiao", '动': "dong",
		'部': "bu", '机': "ji", '问': "wen", '分': "fen",

		// Sensitive word related
		'政': "zheng", '治': "zhi", '色': "se", '情': "qing", '暴': "bao",
		'毒': "du", '品': "pin", '赌': "du", '博': "bo",
		'枪': "qiang", '支': "zhi", '弹': "dan", '药': "yao", '死': "si",
		'杀': "sha", '血': "xue", '腥': "xing", '恐': "kong", '怖': "bu",
	}
}

// ToPinyinInitial converts text to pinyin initials
func ToPinyinInitial(text string) string {
	processor := NewPinyinProcessor()
	fullPinyin := processor.Process(text)

	// Extract first letter of each pinyin syllable
	var result strings.Builder
	words := strings.Fields(fullPinyin)
	for _, word := range words {
		if len(word) > 0 {
			result.WriteByte(word[0])
		}
	}

	return result.String()
}
