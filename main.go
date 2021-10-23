package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sysmon/displaymonitor"
)

func main() {
	notifySession := flag.Bool("notify", false, "If present will notify about session locks and unlocks")
	flag.Parse()

	dm := displaymonitor.New()

	if err := dm.Start(*notifySession); err != nil {
		log.Fatalf("Could not start display monitor: %v", err)
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig
		if err := dm.Stop(); err != nil {
			log.Fatalf("Failed to stop the DisplayMonitor: %v", err)
		}
	}()

EventLoop:
	for {
		e, err := dm.GetEvent()

		if err != nil {
			log.Fatalf("Failed reading event from display manager: %v\n", err)
		}

		switch v := e.(type) {
		case displaymonitor.ResolutionChangeEvent:
			log.Printf("width: %d\theight: %d", v.Width, v.Height)
		case displaymonitor.SessionLockEvent:
			log.Printf("session ID: %d\tchange: %t", v.ID, v.Locked)
		case displaymonitor.DisplayMonitorDone:
			log.Printf("DisplayMonitor finished")
			break EventLoop
		case displaymonitor.DisplayMonitorError:
			log.Fatalf("Received error event: %v", v)
		}
	}

}
