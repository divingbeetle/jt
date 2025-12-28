package core

import (
	"fmt"
	"sort"
	"strings"
)

// Options configures SQL generation.
type Options struct {
	StringCollation string
}

func format(jsonStr string, root *node, opts Options) (string, error) {
	cols, err := formatNode(root, "$", "", opts, 2)
	if err != nil {
		return "", err
	}

	jsonPath := "$"
	if root.types.has(array) {
		jsonPath = "$[*]"
	}

	return fmt.Sprintf(`JSON_TABLE(
  '%s',
  '%s' COLUMNS (
%s
  )
) AS jt`, escape(jsonStr), jsonPath, cols), nil
}

func formatNode(node *node, path, name string, opts Options, depth int) (string, error) {
	isRoot := path == "$" && name == ""
	currPath := path
	if node.types.has(array) {
		currPath = "$"
	}
	currName := name

	// Iterate children
	var lines []string
	for _, key := range sorted(node.children) {
		child := node.children[key]
		childPath := currPath + "." + key
		childName := key
		if currName != "" {
			childName = currName + "_" + key
		}
		line, err := formatNode(child, childPath, childName, opts, depth)
		if err != nil {
			return "", err
		}
		lines = append(lines, line)
	}

	// Leaf node: primitive value
	if len(node.children) == 0 {
		if currName == "" && isRoot {
			currName = "value" // for arrays of primitives
		}
		sqlType, err := toSQL(node.types, opts.StringCollation)
		if err != nil {
			return "", err
		}
		lines = append(lines, fmt.Sprintf("%s%s %s PATH '%s'", indent(depth), currName, sqlType, currPath))
	}

	inner := strings.Join(lines, ",\n")

	// If this is a nested array, wrap it with NESTED PATH
	isNestedArray := node.types.has(array) && !isRoot
	if isNestedArray {
		if inner == "" {
			return "", nil
		}
		inner = addIndent(inner, "  ")
		return fmt.Sprintf("%sNESTED PATH '%s[*]' COLUMNS (\n%s\n%s)", indent(depth), path, inner, indent(depth)), nil
	}

	return inner, nil
}

func addIndent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

var errMixedTypes = fmt.Errorf("mixed types in values")

func toSQL(t typeFlags, collation string) (string, error) {
	t &^= array | null

	var base string
	switch {
	case t == 0:
		base = "VARCHAR(255)"
	case t.only(numeric):
		if t.has(decimal | bigInt) {
			base = "DOUBLE"
		} else {
			base = "INT"
		}
	case t.only(textual):
		if t.has(longText) {
			base = "TEXT"
		} else {
			base = "VARCHAR(255)"
		}
	case t == boolean:
		base = "BOOLEAN"
	case t.has(numeric) && t.has(textual),
		t.has(numeric) && t.has(boolean),
		t.has(textual) && t.has(boolean):
		return "", errMixedTypes
	}

	if collation != "" && (base == "VARCHAR(255)" || base == "TEXT") {
		base += " COLLATE " + collation
	}

	return base, nil
}

func indent(n int) string {
	return strings.Repeat("  ", n)
}

// Sort map keys to have deterministic output
func sorted(m map[string]*node) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Escape single quotes for SQL string literals
func escape(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
