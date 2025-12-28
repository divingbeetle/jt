package core

import (
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     Options
		contains []string
		wantErr  string
	}{
		// Structure
		{"array of objects", `[{"id": 1, "name": "Alice"}]`, Options{}, []string{
			"'$[*]' COLUMNS", "id INT PATH '$.id'", "name VARCHAR(255) PATH '$.name'",
		}, ""},
		{"single object", `{"id": 1}`, Options{}, []string{
			"'$' COLUMNS", "id INT PATH '$.id'",
		}, ""},
		{"primitive array", `[1, 2, 3]`, Options{}, []string{
			"value INT PATH '$'",
		}, ""},
		{"nested object", `{"user": {"profile": {"age": 25}}}`, Options{}, []string{
			"user_profile_age INT PATH '$.user.profile.age'",
		}, ""},
		{"nested array", `{"id": 1, "tags": [{"name": "go"}]}`, Options{}, []string{
			"id INT PATH '$.id'", "NESTED PATH '$.tags[*]' COLUMNS", "tags_name VARCHAR(255) PATH '$.name'",
		}, ""},
		{"deeply nested arrays", `{"depts": [{"teams": [{"lead": "Alice"}]}]}`, Options{}, []string{
			"NESTED PATH '$.depts[*]' COLUMNS", "NESTED PATH '$.teams[*]' COLUMNS",
		}, ""},
		{"primitive array in object", `{"scores": [85, 90]}`, Options{}, []string{
			"scores INT PATH '$'",
		}, ""},

		// Options
		{"collation", `[{"name": "test"}]`, Options{StringCollation: "utf8mb4_unicode_ci"}, []string{
			"VARCHAR(255) COLLATE utf8mb4_unicode_ci",
		}, ""},

		// Edge cases
		{"escapes quotes", `[{"name": "O'Brien"}]`, Options{}, []string{"O''Brien"}, ""},
		{"empty object", `{}`, Options{}, []string{"'$' COLUMNS"}, ""},
		{"nulls with objects", `[null, {"id": 1}]`, Options{}, []string{"id INT PATH '$.id'"}, ""},

		// Errors
		{"invalid JSON", `{invalid}`, Options{}, nil, "invalid JSON"},
		{"empty array", `[]`, Options{}, nil, "empty input"},
		{"mixed types", `[1, "hello"]`, Options{}, nil, "mixed types"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert([]byte(tt.input), tt.opts)

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("missing %q in:\n%s", want, result)
				}
			}
		})
	}
}

func TestToSQL(t *testing.T) {
	tests := []struct {
		name    string
		values  []any
		want    string
		wantErr bool
	}{
		{"empty", nil, "VARCHAR(255)", false},
		{"nulls", []any{nil, nil}, "VARCHAR(255)", false},
		{"int", []any{float64(1)}, "INT", false},
		{"float", []any{float64(1.5)}, "DOUBLE", false},
		{"int promoted to double", []any{float64(1), float64(1.5)}, "DOUBLE", false},
		{"bigint becomes double", []any{float64(2147483648)}, "DOUBLE", false},
		{"bool", []any{true}, "BOOLEAN", false},
		{"string", []any{"hello"}, "VARCHAR(255)", false},
		{"long string", []any{strings.Repeat("x", 300)}, "TEXT", false},
		{"mixed types", []any{float64(1), "text"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := newNode()
			for _, v := range tt.values {
				walk(n, v)
			}
			got, err := toSQL(n.types, "")
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWalk(t *testing.T) {
	t.Run("builds tree", func(t *testing.T) {
		root := newNode()
		walk(root, map[string]any{
			"name": "test",
			"tags": []any{map[string]any{"id": float64(1)}},
		})

		if _, ok := root.children["name"]; !ok {
			t.Error("missing 'name' child")
		}
		tags, ok := root.children["tags"]
		if !ok {
			t.Fatal("missing 'tags' child")
		}
		if !tags.types.has(array) {
			t.Error("tags should have array flag")
		}
		if _, ok := tags.children["id"]; !ok {
			t.Error("missing 'id' child in tags")
		}
	})
}
