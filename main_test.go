package main

import (
	"bytes"
	"testing"
)

func TestXprintidle(t *testing.T) {
	got := xprintidle()
	got = got[:len(got)-1] // remove /n
	if !bytesAreNumbers(t, got) {
		t.Errorf("unexpected result from xprintidle: '%s'", got)
	}
}

// check that each byte in a slice is a number
func bytesAreNumbers(t *testing.T, s []byte) bool {

	// utf8 decimal codes for 0-9
	numbers := []byte{48, 49, 50, 51, 52, 53, 54, 55, 56, 57}

	for i := 0; i < len(s); i++ {
		b := []byte{s[i]}
		if bytes.Contains(numbers, b) {
			continue
		} else {
			return false
		}
	}
	return true
}
