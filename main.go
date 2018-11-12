package main

import (
	"log"
	"os/exec"
)

func main() {
	output := xprintidle()
	log.Printf("%q", output)
}

func xprintidle() (output []byte) {
	output, _ = exec.Command("xprintidle").Output()
	return
}
