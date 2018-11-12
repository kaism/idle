package main

import (
	"errors"
	"testing"
)

func TestBytesAreDigits(t *testing.T) {
	t.Run("all digits", func(t *testing.T) {
		slice := []byte{48, 49, 50, 51, 52, 53, 54, 55, 56, 57}
		got := bytesAreDigits(slice)
		want := true
		assertBool(t, got, want)
	})
	t.Run("contains non-digit", func(t *testing.T) {
		slice := []byte{10}
		got := bytesAreDigits(slice)
		want := false
		assertBool(t, got, want)
	})
}
func TestCheckXprintidle(t *testing.T) {
	t.Run("not installed error", func(t *testing.T) {
		err := errors.New("exec: \"xprintidle\": executable file not found in $PATH")
		got := checkXprintidle([]byte{}, err)
		assertError(t, got, ErrXprintidleRun)
	})
	t.Run("unexpected result error", func(t *testing.T) {
		output := []byte{}
		got := checkXprintidle(output, nil)
		assertError(t, got, ErrXprintidleResult)
	})
}

func assertBool(t *testing.T, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("got '%t' want '%t'", got, want)
	}
}
func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got == nil {
		t.Errorf("wanted error but didn't get one")
	}
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Errorf("got an error but didn't want one: %v", got)
	}
}
