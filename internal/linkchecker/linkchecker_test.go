package linkchecker_test

import (
	"github.com/Jasstkn/link-checker/internal/linkchecker"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
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
	tests := []struct {
		name           string
		links          []string
		valid          bool
		expectedNum    int
		expectedBroken []string
	}{
		{
			name:           "1 link, 0 broken",
			links:          []string{"https://example.com"},
			valid:          true,
			expectedNum:    0,
			expectedBroken: nil,
		},
		{
			name:           "1 link, 1 broken",
			links:          []string{"https://example.com"},
			valid:          false,
			expectedNum:    1,
			expectedBroken: []string{"https://example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch tt.valid {
				case true:
					w.WriteHeader(200)
				case false:
					w.WriteHeader(404)
				}
			}))
			defer server.Close()

			gotNum, gotBroken := linkchecker.ValidateLinks(tt.links)

			if gotNum != tt.expectedNum && !cmp.Equal(gotBroken, tt.expectedBroken) {
				t.Errorf("BrokenLinks(%+v) = %+v, %+v; expected %+v, %+v.", tt.links, gotNum, gotBroken, tt.expectedNum, tt.expectedBroken)
			}
		})
	}
}
