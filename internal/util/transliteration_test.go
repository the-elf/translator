package util

import "testing"

func TestToGeorgian(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"simple word", "gamarjoba", "გამარჯობა"},
		{"apostrophe digraph", "ch'ai", "ჭაი"},
		{"plain digraph before letter fallback", "shen", "შენ"},
		{"uppercase is lowercased", "SUP", "სუფ"},
		{"unmatched byte passes through", "gamarjoba!", "გამარჯობა!"},
		{"empty string", "", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ToGeorgian(c.input)
			if got != c.want {
				t.Errorf("ToGeorgian(%q) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}
