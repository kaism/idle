package main

import (
	"fmt"
	"log"
	"time"
)

const interval time.Duration = 1 * time.Second
const threshold int = 2 * 60 // in seconds
const timeFormat = "Mon Jan 2 15:04:05"

func main() {
	var idle bool = false
	var start time.Time = time.Now()

	fmt.Printf("%v Work ", start.Format(timeFormat))
	for {
		time.Sleep(interval)
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
	}
}
