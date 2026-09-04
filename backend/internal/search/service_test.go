package search

import "testing"

func TestNormalizeSearchQuery(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "   ", expected: ""},
		{name: "trim spaces", input: "  AAPL  ", expected: "AAPL"},
		{name: "keep case-sensitive text", input: "msft", expected: "msft"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeSearchQuery(tc.input); got != tc.expected {
				t.Fatalf("normalizeSearchQuery(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}
