package displaymonitor

import (
	"fmt"
	"log"
	"runtime"

	"golang.org/x/sys/windows"
)

const (
	windowClassName = "DisplayMonitor"
	windowName      = "Display Monitor"
)

const (
	WM_CLOSE             = 0x0010
	WM_DESTROY           = 0x0002
	WM_DISPLAYCHANGE     = 0x007e
	WM_WTSSESSION_CHANGE = 0x02b1
)

const (
	NOTIFY_FOR_ALL_SESSIONS = 0x0001
	WTS_SESSION_LOCK        = 0x0007
	WTS_SESSION_UNLOCK      = 0x0008
)

const (
	eventChanSize = 10
)

const (
	lowerWordMask  = 0xFFFF
	upperWordShift = 16
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

func (dm *displayMonitor) Start(notifySession bool) error {
	go func() {
		runtime.LockOSThread()
		var moduleHandle windows.Handle
		if err := windows.GetModuleHandleEx(0, nil, &moduleHandle); err != nil {
			dm.events <- fmt.Errorf("failed to get module handle: %v", err)
			return
		}

		_, err := registerClassEx(windowClassName, dm.windowProcedure, moduleHandle)
		if err != nil {
			dm.events <- fmt.Errorf("failed to register window class: %v", err)
			return
		}

		windowHandle, err := createWindowExW(windowClassName, windowName, moduleHandle)
		if err != nil {
			dm.events <- fmt.Errorf("failed to create window: %v", err)
			return
		}
		dm.windowHandle = windowHandle

		if notifySession {
			wtsRegisterSessionNotification(windowHandle, NOTIFY_FOR_ALL_SESSIONS)
			if err != nil {
				dm.events <- fmt.Errorf("failed to register for sessions change notifications: %v", err)
				return
			}
		}
		getMessageLoop(windowHandle)
		destroyWindow(windowHandle)
		runtime.UnlockOSThread()
	}()

	return nil
}

func (dm *displayMonitor) Stop() error {
	closeWindow(dm.windowHandle)
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
	switch msg {
	case WM_CLOSE:
		destroyWindow(dm.windowHandle)
	case WM_DESTROY:
		dm.events <- DisplayMonitorDone{}
		postQuitMessage(0)
	case WM_DISPLAYCHANGE:
		dm.events <- ResolutionChangeEvent{
			Width:  int(lparam) & lowerWordMask,
			Height: (int(lparam) >> upperWordShift) & lowerWordMask,
		}
	case WM_WTSSESSION_CHANGE:
		switch wparam {
		case WTS_SESSION_LOCK:
			dm.events <- SessionLockEvent{ID: int(lparam), Locked: true}
		case WTS_SESSION_UNLOCK:
			dm.events <- SessionLockEvent{ID: int(lparam), Locked: false}
		}

	}
	return defWindowProc(hwnd, msg, wparam, lparam)
}

func getMessageLoop(windowHandle windows.Handle) {
	for {
		msg := tMSG{}
		m, err := getMessage(&msg, 0, 0, 0)
		if err != nil {
			log.Println(err)
			break
		}

		if !m {
			break
		}
		dispatchMessage(&msg)

	}
}
