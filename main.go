package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"log"
	"os"
	"time"
)

const interval time.Duration = 1 * time.Second
const threshold int = 2 * 60 // in seconds
const timeFormat = "Mon Jan 2 15:04:05"

var abort = make(chan struct{})

func main() {
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()

	go systray.Run(onReady, onExit)

	var idle bool = false
	var start time.Time = time.Now()

	fmt.Printf("%v Work ", start.Format(timeFormat))
	tick := time.Tick(interval)
loop:
	for {
		seconds, err := getIdleTime()
		if err != nil {
			log.Fatalf("%v", err)
		}
		if changeState(&idle, threshold, seconds) {
			var end time.Time
			thresholdDuration := time.Duration(threshold) * time.Second
			if idle {
				end = time.Now().Add(-thresholdDuration)
			} else {
				end = time.Now().Add(-interval)
			}
			fmt.Printf(stateChangeMsg(idle, start, end))
			start = end
		}

		select {
		case <-tick:
			// log.Printf("tick")
		case <-abort:
			// log.Printf("abort")
			break loop
		}
	}
}
