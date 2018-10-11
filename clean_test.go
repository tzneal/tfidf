package tfidf_test

import (
	"testing"

	"github.com/tzneal/tfidf"
)

func TestClean(t *testing.T) {
	for _, tc := range []struct {
		Input    string
		Expected string
	}{
		{"‘test’", "'test'"},
		{`“test”`, `"test"`},
		{"te\u2014st", "te-st"},
	} {
		got := tfidf.ReplaceSmartQuotes(tc.Input)
		if got != tc.Expected {
			t.Errorf("expected '%s', got '%s'", tc.Expected, got)
		}
	}
}
