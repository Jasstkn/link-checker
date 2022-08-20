package utils_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jasstkn/link-checker/internal/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestParseHTML(t *testing.T) {
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
			got := utils.ParseHTML(tt.body)

			if !cmp.Equal(got, tt.expected) {
				t.Errorf("ParseHTML(%+v) = %+v; expected %+v", tt.body, got, tt.expected)
			}
		})
	}
}

func TestValidateLinks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/":
			w.WriteHeader(200)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	defer server.Close()

	tests := []struct {
		name          string
		urls          []string
		expectedNum   int
		expectedLinks []string
	}{
		{
			name:          "0 broken links",
			urls:          []string{server.URL},
			expectedNum:   0,
			expectedLinks: nil,
		},
		{
			name:          "1 broken link",
			urls:          []string{server.URL + "/broken/url"},
			expectedNum:   1,
			expectedLinks: []string{server.URL + "/broken/url"},
		},
		{
			name:          "2 broken links",
			urls:          []string{server.URL + "/broken/url1", server.URL + "/broken/url2"},
			expectedNum:   2,
			expectedLinks: []string{server.URL + "/broken/url1", server.URL + "/broken/url2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNum, gotLinks := utils.ValidateLinks(tt.urls)

			if gotNum != tt.expectedNum || !assert.ElementsMatch(t, gotLinks, tt.expectedLinks) {
				t.Errorf("BrokenLinks(%+v) = %+v, %+v; expected %+v, %+v.", tt.urls, gotNum, gotLinks, tt.expectedNum, tt.expectedLinks)
			}
		})
	}
}

func TestValidateUrl(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantedErr bool
	}{
		{
			name:      "Valid URL",
			url:       "https://example.com",
			wantedErr: false,
		},
		{
			name:      "Invalid URL",
			url:       "htps://example.com",
			wantedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateURL(tt.url)

			if tt.wantedErr && err == nil {
				t.Errorf("LinkChecker(%+v) expected error but received %v", tt.url, err)
			}
		})
	}
}
