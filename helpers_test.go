package main

import (
	"testing"
	"time"
)

func assertTime(t *testing.T, got, want time.Time) {
	t.Helper()
	if got != want {
		t.Errorf("got '%v' want '%v'", got, want)
	}
}
func assertString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
func assertInt(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
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
	} else if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Errorf("got an error but didn't want one: %v", got)
	}
}
