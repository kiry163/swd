package variant

import "strings"

// TraditionalProcessor handles traditional/simplified Chinese conversion
type TraditionalProcessor struct {
	s2tMap map[rune]rune // Simplified to Traditional
	t2sMap map[rune]rune // Traditional to Simplified
}

// NewTraditionalProcessor creates a new traditional Chinese processor
func NewTraditionalProcessor() *TraditionalProcessor {
	s2t, t2s := buildConversionMaps()
	return &TraditionalProcessor{
		s2tMap: s2t,
		t2sMap: t2s,
	}
}

// Process converts traditional Chinese to simplified for normalization
func (p *TraditionalProcessor) Process(text string) string {
	return p.ToSimplified(text)
}

// Name returns the processor name
func (p *TraditionalProcessor) Name() string {
	return "traditional"
}

// ToSimplified converts traditional Chinese to simplified
func (p *TraditionalProcessor) ToSimplified(text string) string {
	var builder strings.Builder
	builder.Grow(len(text))

	for _, r := range text {
		if simplified, exists := p.t2sMap[r]; exists {
			builder.WriteRune(simplified)
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// ToTraditional converts simplified Chinese to traditional
func (p *TraditionalProcessor) ToTraditional(text string) string {
	var builder strings.Builder
	builder.Grow(len(text))

	for _, r := range text {
		if traditional, exists := p.s2tMap[r]; exists {
			builder.WriteRune(traditional)
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// buildConversionMaps builds simplified<->traditional conversion maps
// In production, use a complete library like github.com/liuzl/gocc
func buildConversionMaps() (map[rune]rune, map[rune]rune) {
	// Extended simplified to traditional mappings for common characters
	s2t := map[rune]rune{
		// Original examples
		'测': '測', '试': '試', '词': '詞', '敏': '敏', '感': '感',

		// Common conversions (100 most frequent)
		'个': '個', '为': '為', '国': '國', '来': '來', '对': '對',
		'们': '們', '时': '時', '会': '會', '过': '過', '发': '發',
		'后': '後', '学': '學', '当': '當', '样': '樣', '还': '還',
		'现': '現', '与': '與', '关': '關', '开': '開', '动': '動',
		'问': '問', '两': '兩', '应': '應', '电': '電', '体': '體',
		'实': '實', '无': '無', '业': '業', '东': '東', '听': '聽',
		'长': '長', '见': '見', '书': '書', '头': '頭', '车': '車',
		'门': '門', '马': '馬', '号': '號', '义': '義', '亲': '親',
		'记': '記', '师': '師', '岁': '歲', '区': '區', '变': '變',
		'压': '壓', '产': '產', '声': '聲', '议': '議', '处': '處',
		'卖': '賣', '买': '買', '战': '戰', '认': '認', '让': '讓',
		'从': '從', '结': '結', '给': '給', '节': '節', '独': '獨',
		'飞': '飛', '万': '萬', '风': '風', '办': '辦', '务': '務',
		'写': '寫', '观': '觀', '习': '習', '报': '報', '场': '場',
		'带': '帶', '队': '隊', '导': '導', '经': '經', '运': '運',
		'历': '歷', '类': '類', '总': '總', '医': '醫', '张': '張',
		'级': '級', '约': '約', '组': '組', '继': '繼', '断': '斷',
		'将': '將', '专': '專', '传': '傳', '达': '達', '亚': '亞',
		'连': '連', '选': '選', '价': '價', '则': '則', '较': '較',
		'尔': '爾', '转': '轉', '规': '規', '参': '參', '标': '標',

		// Sensitive words related
		'党': '黨', '权': '權', '湾': '灣',
		'统': '統', '复': '復', '兴': '興', '华': '華',
	}

	t2s := make(map[rune]rune)
	for k, v := range s2t {
		t2s[v] = k
	}

	return s2t, t2s
}
