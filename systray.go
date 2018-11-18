package main

import (
	"github.com/getlantern/systray"
	"github.com/sfkrystal/idle/icon"
)

func onReady() {
	systray.SetIcon(icon.Dino)
	mQuit := systray.AddMenuItem("Quit", "")
	for {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
			abort <- struct{}{}
		}
	}
}

func onExit() {
}
