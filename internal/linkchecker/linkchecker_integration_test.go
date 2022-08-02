//go:build integration
// +build integration

package linkchecker_test

import (
	"github.com/Jasstkn/link-checker/internal/linkchecker"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestIntegrationValidateLinks(t *testing.T) {
	tests := []struct {
		name           string
		links          []string
		expectedNum    int
		expectedBroken []string
	}{
		{
			name:           "1 link, 0 broken",
			links:          []string{"https://example.com"},
			expectedNum:    0,
			expectedBroken: nil,
		},
		{
			name:           "Some broken",
			links:          []string{"https://github.com/Jasstkn/test-repo.git"},
			expectedNum:    1,
			expectedBroken: []string{"https://github.com/Jasstkn/test-repo.git"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNum, gotBroken := linkchecker.ValidateLinks(tt.links)

			if gotNum != tt.expectedNum && !cmp.Equal(gotBroken, tt.expectedBroken) {
				t.Errorf("BrokenLinks(%+v) = %+v, %+v; expected %+v, %+v.", tt.links, gotNum, gotBroken, tt.expectedNum, tt.expectedBroken)
			}
		})
	}
}

func TestIntegration(t *testing.T) {

}
