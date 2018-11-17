package main

import (
	"github.com/getlantern/systray"
	"github.com/sfkrystal/idle/icon"
	"log"
)

func onReady() {
	log.Println("systray: Starting onReady")
	systray.SetIcon(icon.Dino)
	systray.Quit()
	log.Println("systray: Finished onReady")
}

func onExit() {
	log.Println("systray: Starting onExit")
	log.Println("systray: Finished onExit")
}
