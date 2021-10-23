package displaymonitor

import "fmt"

type displayMonitor struct {
	events chan Event
}

func New() DisplayMonitor {
	return &displayMonitor{
		events: make(chan Event, 10),
	}
}

func (dm *displayMonitor) Start() error {
	return fmt.Errorf("not implemented")

}
func (dm *displayMonitor) Stop() error {
	return fmt.Errorf("not implemented")
}
func (dm *displayMonitor) GetEvent() (Event, error) {
	return nil, fmt.Errorf("not implemented")
}
