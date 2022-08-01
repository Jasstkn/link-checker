package linkchecker_test

import (
	"github.com/Jasstkn/link-checker/internal/linkchecker"
	"github.com/google/go-cmp/cmp"
	"reflect"
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

func TestValidateLinks(t *testing.T) {
	type BrokenLinks struct {
		num    int
		broken []string
	}

	tests := []struct {
		name     string
		links    []string
		expected BrokenLinks
	}{
		{
			name:     "1 link, 0 broken",
			links:    []string{"https://example.com"},
			expected: BrokenLinks{0, nil},
		},
		{
			name:     "1 broken",
			links:    []string{"https://github.com/Jasstkn/link-checker"},
			expected: BrokenLinks{1, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BrokenLinks{}
			got.num, got.broken = linkchecker.ValidateLinks(tt.links)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("BrokenLinks(%+v) = %+v; expected %+v.", tt.links, got, tt.expected)
			}
		})
	}
}
