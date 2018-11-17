package main

import (
	"github.com/getlantern/systray"
	"github.com/sfkrystal/idle/icon"
)

func onReady() {
	systray.SetIcon(icon.Dino)
	mToggle := systray.AddMenuItem("Set Idle", "")
	go func() {
		idle := true
		for {
			select {
			case <-mToggle.ClickedCh:
				if idle {
					mToggle.SetTitle("Set Work")
					systray.SetIcon(icon.DinoSleep)
					idle = false
				} else {
					mToggle.SetTitle("Set Idle")
					systray.SetIcon(icon.Dino)
					idle = true
				}
			}
		}
	}()
}

func onExit() {
	systray.Quit()
}
