package charset

import "unicode"

// IsChineseChar checks if a rune is a Chinese character
func IsChineseChar(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

// IsCJK checks if a rune is a CJK (Chinese, Japanese, Korean) character
func IsCJK(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
		(r >= 0x3400 && r <= 0x4DBF) || // CJK Extension A
		(r >= 0x20000 && r <= 0x2A6DF) || // CJK Extension B
		(r >= 0x2A700 && r <= 0x2B73F) || // CJK Extension C
		(r >= 0x2B740 && r <= 0x2B81F) || // CJK Extension D
		(r >= 0x2B820 && r <= 0x2CEAF) || // CJK Extension E
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility Ideographs
		(r >= 0x2F800 && r <= 0x2FA1F) // CJK Compatibility Ideographs Supplement
}

// IsHangul checks if a rune is a Korean Hangul character
func IsHangul(r rune) bool {
	return (r >= 0xAC00 && r <= 0xD7AF) || // Hangul Syllables
		(r >= 0x1100 && r <= 0x11FF) || // Hangul Jamo
		(r >= 0x3130 && r <= 0x318F) || // Hangul Compatibility Jamo
		(r >= 0xA960 && r <= 0xA97F) || // Hangul Jamo Extended-A
		(r >= 0xD7B0 && r <= 0xD7FF) // Hangul Jamo Extended-B
}

// IsJapanese checks if a rune is a Japanese character (Hiragana or Katakana)
func IsJapanese(r rune) bool {
	return (r >= 0x3040 && r <= 0x309F) || // Hiragana
		(r >= 0x30A0 && r <= 0x30FF) // Katakana
}

// ToLower converts a rune to lowercase (handles ASCII and extended characters)
func ToLower(r rune) rune {
	return unicode.ToLower(r)
}

// ToUpper converts a rune to uppercase (handles ASCII and extended characters)
func ToUpper(r rune) rune {
	return unicode.ToUpper(r)
}

// IsAlphanumeric checks if a rune is alphanumeric
func IsAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// IsSymbol checks if a rune is a symbol or punctuation
func IsSymbol(r rune) bool {
	return unicode.IsSymbol(r) || unicode.IsPunct(r)
}

// CharType returns the type of character
type CharType int

const (
	CharTypeUnknown CharType = iota
	CharTypeChinese
	CharTypeEnglish
	CharTypeDigit
	CharTypeSymbol
	CharTypeSpace
)

// GetCharType returns the type of a character
func GetCharType(r rune) CharType {
	if unicode.IsSpace(r) {
		return CharTypeSpace
	}
	if unicode.IsDigit(r) {
		return CharTypeDigit
	}
	if IsChineseChar(r) {
		return CharTypeChinese
	}
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return CharTypeEnglish
	}
	if IsSymbol(r) {
		return CharTypeSymbol
	}
	return CharTypeUnknown
}


