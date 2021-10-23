package main

import (
	"log"
	"sysmon/displaymonitor"
)

func main() {
	dm := displaymonitor.New()

	if err := dm.Start(); err != nil {
		log.Fatalf("Could not start display monitor: %v", err)
	}

Loop:
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
			break Loop
		}
	}
}
