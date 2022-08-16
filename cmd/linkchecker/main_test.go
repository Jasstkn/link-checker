package main

import (
	"testing"
)

func Test(t *testing.T) {
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
			err := ValidateURL(tt.url)

			if tt.wantedErr && err == nil {
				t.Errorf("LinkChecker(%+v) expected error but received %v", tt.url, err)
			}
		})
	}
}
