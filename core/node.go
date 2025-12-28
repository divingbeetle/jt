package core

import "math"

// Tracks value type information observed during JSON traversal.
// A node can have multiple type flags set if the JSON data is heterogeneous.
type typeFlags uint16

func (t typeFlags) has(flag typeFlags) bool {
	return t&flag != 0
}

func (t typeFlags) only(flags typeFlags) bool {
	return t&^flags == 0
}

const (
	array typeFlags = 1 << iota
	integer
	decimal
	bigInt
	shortText
	longText
	boolean
	null

	numeric = integer | decimal | bigInt
	textual = shortText | longText
)

const longTextLen = 255

// Represents a node in the JSON structure tree.
type node struct {
	types    typeFlags
	children map[string]*node
}

func newNode() *node {
	return &node{children: make(map[string]*node)}
}

func (node *node) getOrCreateChild(key string) *node {
	child, found := node.children[key]
	if !found {
		child = newNode()
		node.children[key] = child
	}
	return child
}

// Traverse and record type information
func walk(node *node, v any) {
	switch val := v.(type) {
	case []any:
		node.types |= array
		for _, item := range val {
			walk(node, item)
		}
	case map[string]any:
		for key, item := range val {
			walk(node.getOrCreateChild(key), item)
		}
	case float64:
		switch {
		case val != math.Trunc(val):
			node.types |= decimal
		case val < math.MinInt32 || val > math.MaxInt32:
			node.types |= bigInt
		default:
			node.types |= integer
		}
	case string:
		if len(val) > longTextLen {
			node.types |= longText
		} else {
			node.types |= shortText
		}
	case bool:
		node.types |= boolean
	case nil:
		node.types |= null
	}
}
