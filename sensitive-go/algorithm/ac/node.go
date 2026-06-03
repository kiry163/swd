package ac

import "github.com/Karrecy/sensitive-go/dict"

// Node represents a node in the Aho-Corasick trie
type Node struct {
	children map[rune]*Node // Child nodes
	fail     *Node          // Failure pointer for AC automation
	word     *dict.Word     // The word if this is a terminal node
	isEnd    bool           // Whether this node marks the end of a word
}

// newNode creates a new trie node
func newNode() *Node {
	return &Node{
		children: make(map[rune]*Node),
		fail:     nil,
		word:     nil,
		isEnd:    false,
	}
}

// addChild adds a child node for the given rune
func (n *Node) addChild(r rune) *Node {
	if child, exists := n.children[r]; exists {
		return child
	}
	child := newNode()
	n.children[r] = child
	return child
}

// getChild retrieves a child node for the given rune
func (n *Node) getChild(r rune) (*Node, bool) {
	child, exists := n.children[r]
	return child, exists
}

// setWord marks this node as a terminal node with the associated word
func (n *Node) setWord(word *dict.Word) {
	n.word = word
	n.isEnd = true
}

// hasChild checks if a child exists for the given rune
func (n *Node) hasChild(r rune) bool {
	_, exists := n.children[r]
	return exists
}


