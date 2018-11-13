package main

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var interval = 1 * time.Second
var threshold = 3 // in seconds

func main() {
	var idle bool = false
	for {
		time.Sleep(interval)
		seconds, err := getIdleTime()
		if err != nil {
			log.Fatalf("%v", err)
		}

		if changeState(&idle, threshold, seconds) {
			if idle {
				log.Println("Idle start")
			} else {
				log.Println("Idle stop")
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
		return num, ErrParse
	}
	return num, err
}

// checks that xprintidle ran and output is what we expect, eg: 1234\n
func checkXprintidle(output []byte, inErr error) error {
	if inErr != nil {
		return ErrXprintidleRun
	}
	if len(output) < 1 {
		return ErrXprintidleResult
	}
	if output[len(output)-1] == 10 { // remove /n
		output = output[:len(output)-1]
	}
	if !bytesAreDigits(output) {
		return ErrXprintidleResult
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

var ErrXprintidleRun = errors.New("error running xprintidle (not installed?)")
var ErrXprintidleResult = errors.New("unexpected result from xprintidle")
var ErrParse = errors.New("parse error")
