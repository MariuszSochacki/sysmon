package displaymonitor

import (
	"fmt"

	"golang.org/x/sys/windows"
)

const (
	windowClassName = "DisplayMonitor"
	windowName      = "Display Monitor"
)

const (
	eventChanSize = 10
)

type displayMonitor struct {
	events       chan Event
	windowHandle windows.Handle
}

func New() DisplayMonitor {
	return &displayMonitor{
		events: make(chan Event, eventChanSize),
	}
}

func (dm *displayMonitor) Start() error {
	var moduleHandle windows.Handle
	if err := windows.GetModuleHandleEx(0, nil, &moduleHandle); err != nil {
		return fmt.Errorf("failed to get module handle: %v", err)
	}

	_, err := registerClassEx(windowClassName, dm.windowProcedure, moduleHandle)
	if err != nil {
		return fmt.Errorf("failed to register window class: %v", err)
	}

	windowHandle, err := createWindowExW(windowClassName, windowName, moduleHandle)
	if err != nil {
		return fmt.Errorf("failed to create window: %v", err)
	}
	dm.windowHandle = windowHandle

	return nil
}

func (dm *displayMonitor) Stop() error {
	dm.events <- DisplayMonitorDone{}
	if err := destroyWindow(dm.windowHandle); err != nil {
		return fmt.Errorf("failed to destroy window: %v", err)
	}
	return nil
}

func (dm *displayMonitor) GetEvent() (Event, error) {
	e, ok := <-dm.events
	if !ok {
		return nil, fmt.Errorf("could not read events")
	}
	switch e.(type) {
	case DisplayMonitorDone:
		close(dm.events)
	}
	return e, nil
}

func (dm *displayMonitor) windowProcedure(hwnd windows.Handle, msg uint32, wparam, lparam uintptr) uintptr {
	return defWindowProc(hwnd, msg, wparam, lparam)
}
