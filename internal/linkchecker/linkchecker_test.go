package linkchecker_test

import (
	"github.com/Jasstkn/link-checker/internal/linkchecker"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseHtml(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected []string
	}{
		{
			name:     "Empty body",
			body:     "",
			expected: []string{},
		},
		{
			name:     "Found links",
			body:     "<a href=\"https://example.com\">",
			expected: []string{"https://example.com"},
		},
		{
			name:     "Links without schema",
			body:     "<a href=\"example.com\">",
			expected: []string{},
		},
		{
			name:     "Links to anchors",
			body:     "<a href=\"#head\">",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := linkchecker.ParseHtml(tt.body)

			if !cmp.Equal(got, tt.expected) {
				t.Errorf("ParseHtml(%+v) = %+v; expected %+v", tt.body, got, tt.expected)
			}
		})
	}
}
