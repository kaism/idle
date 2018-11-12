package main

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
)

func main() {
	output, err := xprintidle()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("%q", output)
}

// tries to run xprintidle and checks output
func xprintidle() (output []byte, err error) {
	output, err = exec.Command("xprintidle").Output()
	err = checkXprintidle(output, err)
	return
}
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
