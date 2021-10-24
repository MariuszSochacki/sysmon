package displaymonitor

// New creates a new object of a display monitor and allows using it's interface.
func New() DisplayMonitor {
	return newImpl()
}

// DisplayMonitor is the interface that allows controlling the display monitoring service.
type DisplayMonitor interface {
	// Start starts the monitoring process and blocks until Stop is called or an error occurs.
	// As long as Start is running, any events will be passed to their respective handlers.
	Start() error
	// Stop stops monitoring process. Calling Stop will unblock the previous call to Start.
	Stop() error
	// SetResolutionChangeHandler sets the handler for the resolution change event.
	SetResolutionChangeHandler(func(ResolutionChangeEvent))
	// SetResolutionChangeHandler sets the handler for the session lock event.
	SetSessionLockHandler(func(SessionLockEvent))
}

// ResolutionChangeEvent is the structure holding information about a resolution change event.
// When a resolution change event occurs the resolution change handler will be called with
// an instance of this structure.
type ResolutionChangeEvent struct {
	// Width is the new width of the screen.
	Width int
	// Height is the new height of the screen.
	Height int
}

// SessionLockEvent is the structure holding information about a session lock event.
// When a session lock event occurs the resolution change handler will be called with
// an instance of this structure.
type SessionLockEvent struct {
	// ID is the ID of the session in which an event occured.
	ID int
	// Locked is true when the session has been locked and false if session has been unlocked.
	Locked bool
}

type displayMonitor struct {
	resolutionChangeHandler func(ResolutionChangeEvent)
	sessionLockHandler      func(SessionLockEvent)
}

func (dm *displayMonitor) SetResolutionChangeHandler(handler func(ResolutionChangeEvent)) {
	dm.resolutionChangeHandler = handler
}

func (dm *displayMonitor) SetSessionLockHandler(handler func(SessionLockEvent)) {
	dm.sessionLockHandler = handler
}
