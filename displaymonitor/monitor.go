package displaymonitor

type DisplayMonitor interface {
	Start(notifySession bool) error
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
type DisplayMonitorError error
