package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sysmon/displaymonitor"
)

func resolutionChangeHandler(e displaymonitor.ResolutionChangeEvent) {
	log.Printf("width: %d\theight: %d", e.Width, e.Height)
}

func sessionLockHandler(e displaymonitor.SessionLockEvent) {
	log.Printf("session ID: %d\tchange: %t", e.ID, e.Locked)
}

func main() {
	notifySession := flag.Bool("notify", false, "If present will notify about session locks and unlocks")
	flag.Parse()

	dm := displaymonitor.New()

	dm.SetResolutionChangeHandler(resolutionChangeHandler)

	if *notifySession {
		dm.SetSessionLockHandler(sessionLockHandler)
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig
		if err := dm.Stop(); err != nil {
			log.Fatalf("Failed to stop the DisplayMonitor: %v", err)
		}
	}()

	if err := dm.Start(); err != nil {
		log.Fatalf("Display monitor failed: %v", err)
	}
}
