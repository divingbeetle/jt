// Package core converts JSON data to MySQL JSON_TABLE expressions.
//
// It analyzes JSON structure and automatically infers appropriate SQL types:
//   - Integers → INT
//   - Floats → DOUBLE
//   - Booleans → BOOLEAN
//   - Strings ≤255 chars → VARCHAR(255)
//   - Strings >255 chars → TEXT
//
// Nested objects are flattened with underscore-separated column names.
// Nested arrays become NESTED PATH expressions.
//
// Example usage:
//
//	result, err := core.Convert([]byte(`[{"id": 1, "name": "Alice"}]`), core.Options{})
//	// result: JSON_TABLE('[{"id":1,"name":"Alice"}]', '$[*]' COLUMNS (...)) AS jt
package core
