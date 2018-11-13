package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const interval = 1 * time.Second
const threshold = 5 * 60 // in seconds
const timeFormat = "Mon Jan 2 15:04:05 MST 2006"

func main() {
	var idle bool = false
	var startIdle time.Time
	var startWork time.Time = time.Now()

	fmt.Printf("%v Work", startWork.Format(timeFormat))
	for {
		time.Sleep(interval)
		seconds, err := getIdleTime()
		if err != nil {
			log.Fatalf("%v", err)
		}
		if changeState(&idle, threshold, seconds) {
			if idle {
				// complete work statement
				duration := time.Since(startWork).Truncate(time.Second).String()
				fmt.Printf(" for %s\n", duration)
				// begin idle statement
				startIdle = time.Now().Add(-time.Duration(seconds) * time.Second)
				fmt.Printf("%v Idle", startIdle.Format(timeFormat))
			} else {
				// complete idle statement
				duration := time.Since(startIdle).Truncate(time.Second).String()
				fmt.Printf(" for %s\n", duration)
				// begin work statement
				startWork = time.Now()
				fmt.Printf("%v Work", startIdle.Format(timeFormat))

			}
		}
	}
}

// returns true if the state was changed
func changeState(idle *bool, threshold int, seconds int) bool {
	if seconds >= threshold && !*idle {
		*idle = true
		return true
	}
	if seconds < threshold && *idle {
		*idle = false
		return true
	}
	return false
}

// returns idle time in seconds
func getIdleTime() (seconds int, err error) {
	output, err := xprintidle()
	if err != nil {
		return
	}
	msecs, err := parseXprintidleOutput(output)
	if err != nil {
		return
	}
	seconds = msecs / 1000 // always rounds down
	return
}

// tries to run xprintidle and checks output
func xprintidle() (output []byte, err error) {
	output, err = exec.Command("xprintidle").Output()
	err = checkXprintidle(output, err)
	return
}

// parse xprintidle output to int, eg: 1234\n to 1234 (int)
func parseXprintidleOutput(bytes []byte) (int, error) {
	str := strings.TrimSpace(string(bytes))
	num, err := strconv.Atoi(str)
	if err != nil {
		return num, errParse
	}
	return num, err
}

// checks that xprintidle ran and output is what we expect, eg: 1234\n
func checkXprintidle(output []byte, inErr error) error {
	if inErr != nil {
		return errXprintidleRun
	}
	if len(output) < 1 {
		return errXprintidleResult
	}
	if output[len(output)-1] == 10 { // remove /n
		output = output[:len(output)-1]
	}
	if !bytesAreDigits(output) {
		return errXprintidleResult
	}
	return nil
}

// checks that all bytes in a slice are digits
func bytesAreDigits(s []byte) bool {
	digits := []byte{48, 49, 50, 51, 52, 53, 54, 55, 56, 57} // utf8 decimal codes for 0-9
	for i := 0; i < len(s); i++ {
		b := []byte{s[i]}
		if !bytes.Contains(digits, b) {
			return false
		}
	}
	return true
}

var errXprintidleRun = errors.New("error running xprintidle (not installed?)")
var errXprintidleResult = errors.New("unexpected result from xprintidle")
var errParse = errors.New("parse error")
