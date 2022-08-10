package linkchecker_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jasstkn/link-checker/pkg/linkchecker"
	"github.com/google/go-cmp/cmp"
)

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
		case "/partial":
			w.WriteHeader(200)
			fmt.Fprintln(w, "<a href=\"http://"+r.Host+"/\">\n<a href=\"http://"+r.Host+"/broken-url1\">")
		case "/broken2":
			w.WriteHeader(200)
			fmt.Fprintln(w, "<a href=\"http://"+r.Host+"/broken-url1\">\n<a href=\"http://"+r.Host+"/broken-url2\">")
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	defer server.Close()

	tests := []struct {
		name     string
		url      string
		expected string
		err      bool
	}{
		{
			name:     "0 total",
			url:      server.URL + "/empty",
			expected: "No links were found",
			err:      false,
		},
		{
			name:     "1 total: 0 broken",
			url:      server.URL + "/",
			expected: "1 link scanned, 0 broken links found",
			err:      false,
		},
		{
			name:     "1 total: 1 broken",
			url:      server.URL + "/broken",
			expected: "1 link scanned, 1 broken link found:\n" + server.URL + "/broken-url",
			err:      false,
		},
		{
			name:     "2 total: 2 broken",
			url:      server.URL + "/broken2",
			expected: "2 links scanned, 2 broken links found:\n" + server.URL + "/broken-url1;\n" + server.URL + "/broken-url2",
			err:      false,
		},
		{
			name:     "2 total: 1 broken",
			url:      server.URL + "/partial",
			expected: "2 links scanned, 1 broken link found:\n" + server.URL + "/broken-url1",
			err:      false,
		},
		{
			name:     "wrong protocol: hhttp",
			url:      "h" + server.URL,
			expected: "",
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := linkchecker.LinkChecker(tt.url)

			if tt.err && err == nil {
				t.Errorf("LinkChecker(%+v) expected error but received %v", tt.url, err)
			}

			if !cmp.Equal(got, tt.expected) {
				t.Errorf("LinkChecker(%+v) = %+v; expected %+v", tt.url, got, tt.expected)
			}
		})
	}
}
