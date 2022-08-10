package linkchecker_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jasstkn/link-checker/internal/linkchecker"
	"github.com/google/go-cmp/cmp"
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/":
			w.WriteHeader(200)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	type args struct {
		url string
	}
	defer server.Close()

	tests := []struct {
		name          string
		args          args
		expectedNum   int
		expectedLinks []string
	}{
		{
			name: "0 broken links",
			args: args{
				url: server.URL,
			},
			expectedNum:   0,
			expectedLinks: nil,
		},
		{
			name: "1 broken link",
			args: args{
				url: server.URL + "/broken/url",
			},
			expectedNum:   1,
			expectedLinks: []string{server.URL + "/broken/url"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNum, gotLinks := linkchecker.ValidateLinks([]string{tt.args.url})

			if gotNum != tt.expectedNum || !cmp.Equal(gotLinks, tt.expectedLinks) {
				t.Errorf("BrokenLinks(%+v) = %+v, %+v; expected %+v, %+v.", tt.args.url, gotNum, gotLinks, tt.expectedNum, tt.expectedLinks)
			}
		})
	}
}

func TestLinkChecker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/":
			w.WriteHeader(200)
			fmt.Fprintln(w, "<a href=\"http://"+r.Host+"\">")
		case "/empty":
			w.WriteHeader(200)
		case "/broken":
			w.WriteHeader(200)
			fmt.Fprintln(w, "<a href=\"http://"+r.Host+"/broken-url\">")
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	defer server.Close()

	tests := []struct {
		name     string
		url      string
		expected string
		err      error
	}{
		{
			name:     "No links",
			url:      server.URL + "/empty",
			expected: "No links were found",
			err:      nil,
		},
		{
			name:     "0 broken",
			url:      server.URL + "/",
			expected: "1 links scanned, 0 broken links found",
			err:      nil,
		},
		{
			name:     "1 broken",
			url:      server.URL + "/broken",
			expected: "1 links scanned, 1 broken link found:\n" + server.URL + "/broken-url",
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := linkchecker.LinkChecker(tt.url)

			if got != tt.expected || err != tt.err {
				t.Errorf("LinkChecker(%+v) = %+v, %+v; expected %+v, %+v.", tt.url, got, err, tt.expected, tt.err)
			}
		})
	}
}
