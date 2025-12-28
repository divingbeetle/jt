package core

import (
	"encoding/json"
	"fmt"
)

// Convert transforms JSON data into MySQL JSON_TABLE format.
func Convert(jsonData []byte, opts Options) (string, error) {
	var data any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	if items, ok := data.([]any); ok && len(items) == 0 {
		return "", fmt.Errorf("empty input")
	}

	root := newNode()
	walk(root, data)

	jsonStr, _ := json.Marshal(data)
	return format(string(jsonStr), root, opts)
}
