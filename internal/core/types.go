package core

// Word describes a user-supplied sensitive word.
type Word struct {
	Text string
	Type string
}

// Match is a detection result in rune positions.
type Match struct {
	Word     string
	Type     string
	Text     string
	StartPos int
	EndPos   int
}
