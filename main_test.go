package main

import (
	"errors"
	"testing"
)

func TestChangeState(t *testing.T) {
	threshold := 5 * 60
	t.Run("returns true when changed", func(t *testing.T) {
		t.Run("idle to not idle", func(t *testing.T) {
			idle := true
			seconds := 0
			got := changeState(&idle, threshold, seconds)
			assertBool(t, got, true)
		})
		t.Run("not idle to idle", func(t *testing.T) {
			idle := false
			seconds := (5 * 60) + 1
			got := changeState(&idle, threshold, seconds)
			assertBool(t, got, true)
		})
	})
	t.Run("returns false when not changed", func(t *testing.T) {
		t.Run("idle to idle", func(t *testing.T) {
			idle := true
			seconds := (5 * 60) + 2
			got := changeState(&idle, threshold, seconds)
			assertBool(t, got, false)
		})
		t.Run("not idle to not idle", func(t *testing.T) {
			idle := false
			seconds := 0
			got := changeState(&idle, threshold, seconds)
			assertBool(t, got, false)
		})
	})
	t.Run("changes idle when it should", func(t *testing.T) {
		t.Run("idle to not idle", func(t *testing.T) {
			idle := true
			seconds := 0
			_ = changeState(&idle, threshold, seconds)
			assertBool(t, idle, false)
		})
		t.Run("not idle to idle", func(t *testing.T) {
			idle := false
			seconds := (5 * 60) + 1
			_ = changeState(&idle, threshold, seconds)
			assertBool(t, idle, true)
		})
	})
	t.Run("doesn't change idle when it shouldn't", func(t *testing.T) {
		t.Run("idle to idle", func(t *testing.T) {
			idle := true
			seconds := (5 * 60) + 2
			_ = changeState(&idle, threshold, seconds)
			assertBool(t, idle, true)
		})
		t.Run("not idle to not idle", func(t *testing.T) {
			idle := false
			seconds := 0
			_ = changeState(&idle, threshold, seconds)
			assertBool(t, idle, false)
		})
	})

}
func TestParseXprintidleOutput(t *testing.T) {
	t.Run("parse byte slice to int", func(t *testing.T) {
		got, _ := parseXprintidleOutput([]byte{49, 48, 50, 51, 52, 53, 54, 55, 56, 57})
		want := 1023456789
		assertInt(t, got, want)
	})
	t.Run("strconv.Atoi parse error", func(t *testing.T) {
		_, got := parseXprintidleOutput([]byte{10})
		want := errParse
		assertError(t, got, want)
	})
}
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
		assertError(t, got, errXprintidleRun)
	})
	t.Run("unexpected result error", func(t *testing.T) {
		output := []byte{}
		got := checkXprintidle(output, nil)
		assertError(t, got, errXprintidleResult)
	})
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
