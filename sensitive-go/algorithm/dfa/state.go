package dfa

import "github.com/Karrecy/sensitive-go/dict"

// State represents a state in the DFA with type-safe transitions
type State struct {
	transitions map[rune]*State // Character transitions to next states
	word        *dict.Word      // Word info if this is an end state
	isEnd       bool            // Whether this is an end state
}

// newState creates a new DFA state
func newState() *State {
	return &State{
		transitions: make(map[rune]*State),
		word:        nil,
		isEnd:       false,
	}
}

// isEndState checks if this state marks the end of a word
func (s *State) isEndState() bool {
	return s.isEnd
}

// getWord retrieves the word information from an end state
func (s *State) getWord() *dict.Word {
	return s.word
}

// setWord marks this state as an end state with word information
func (s *State) setWord(word *dict.Word) {
	s.word = word
	s.isEnd = true
}

// transition returns the next state for the given rune
func (s *State) transition(r rune) (*State, bool) {
	next, exists := s.transitions[r]
	return next, exists
}

// addTransition adds a transition to the next state
func (s *State) addTransition(r rune, next *State) {
	s.transitions[r] = next
}
