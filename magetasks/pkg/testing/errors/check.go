package errors

import (
	"errors"
	"testing"
)

// Check if expected and actual error values are correct.
func Check(tb testing.TB, got, want error) {
	tb.Helper()
	if want != nil {
		if !errors.Is(got, want) {
			tb.Fatalf("want %v, got %v", want, got)
		}
	} else {
		if got != nil {
			tb.Fatal("got unexpected:", got)
		}
	}
}
