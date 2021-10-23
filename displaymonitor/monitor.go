package displaymonitor

type DisplayMonitor interface {
	Start() error
	Stop() error
	GetEvent() (Event, error)
}

type Event interface{}
type ResolutionChangeEvent struct {
	Width, Height int
}
type SessionLockEvent struct {
	ID     int
	Locked bool
}
type DisplayMonitorDone struct{}
