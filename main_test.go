package main

import (
	"errors"
	"testing"
	"time"
)

// func TestCalcPreviousStateEnd(t *testing.T) {
// 	threshold := 5 * time.Minute
// 	interval := 1 * time.Second

// 	// start 18:20, worked for 10 mins
// 	// now it's 18:35 and we get idle flag after 5 mins* idle		* threshold
// 	// end of work was 5 mins* ago at 18:30							* threshold
// 	t.Run("previous state was work", func(t *testing.T) {
// 		idle := true
// 		now := time.Date(2018, time.November, 10, 18, 35, 0, 0, time.UTC)

// 		got := calcPreviousStateEnd(idle, now, threshold, interval)
// 		want := time.Date(2018, time.November, 10, 18, 30, 0, 0, time.UTC)
// 		assertTime(t, got, want)
// 	})
// 	// start 18:30:00, was idle for 10 mins
// 	// now it's 18:40 we get work flag
// 	// end of idle is up to but not including 1 sec* ago 			* interval
// 	t.Run("previous state was idle", func(t *testing.T) {
// 		idle := false
// 		now := time.Date(2018, time.November, 10, 18, 35, 0, 0, time.UTC)

// 		got := calcPreviousStateEnd(idle, now, threshold, interval)
// 		want := time.Date(2018, time.November, 10, 18, 34, 59, 0, time.UTC)
// 		assertTime(t, got, want)
// 	})
// }

func TestStateChangeMsg(t *testing.T) {
	t.Run("work to idle", func(t *testing.T) {
		idle := true
		start := time.Date(2018, time.November, 10, 18, 20, 40, 0, time.UTC)
		end := time.Date(2018, time.November, 10, 19, 25, 0, 0, time.UTC)

		got := stateChangeMsg(idle, start, end)
		want := "for 1h4m20s\nSat Nov 10 19:25:00 Idle "
		assertString(t, got, want)
	})
	t.Run("idle to work", func(t *testing.T) {
		idle := false
		start := time.Date(2018, time.November, 10, 18, 20, 40, 0, time.UTC)
		end := time.Date(2018, time.November, 10, 19, 25, 0, 0, time.UTC)

		got := stateChangeMsg(idle, start, end)
		want := "for 1h4m20s\nSat Nov 10 19:25:00 Work "
		assertString(t, got, want)
	})
}
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
