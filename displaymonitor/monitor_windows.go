package displaymonitor

import (
	"fmt"
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
	WM_USER              = 0x0400

	WM_STOP_DISPLAY_MONITOR = WM_USER
)

const (
	NOTIFY_FOR_ALL_SESSIONS = 0x0001
	WTS_SESSION_LOCK        = 0x0007
	WTS_SESSION_UNLOCK      = 0x0008
)

const (
	lowerWordMask  = 0xFFFF
	upperWordShift = 16
)

type displayMonitorWindows struct {
	displayMonitor
	windowHandle windows.Handle
}

func newImpl() DisplayMonitor {
	return &displayMonitorWindows{}
}

func (dm *displayMonitorWindows) Start() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var moduleHandle windows.Handle
	if err := windows.GetModuleHandleEx(0, nil, &moduleHandle); err != nil {
		return fmt.Errorf("failed to get module handle: %v", err)
	}

	_, err := registerClassEx(windowClassName, dm.windowProcedure, moduleHandle)
	if err != nil {
		return fmt.Errorf("failed to register for sessions change notifications: %v", err)
	}

	windowHandle, err := createWindowExW(windowClassName, windowName, moduleHandle)
	if err != nil {
		return fmt.Errorf("failed to register for sessions change notifications: %v", err)
	}
	dm.windowHandle = windowHandle

	if dm.sessionLockHandler != nil {
		wtsRegisterSessionNotification(windowHandle, NOTIFY_FOR_ALL_SESSIONS)
		if err != nil {
			return fmt.Errorf("failed to register for sessions change notifications: %v", err)
		}
	}

	return getMessageLoop()
}

func (dm *displayMonitorWindows) Stop() error {
	sendMessageW(dm.windowHandle, WM_STOP_DISPLAY_MONITOR, 0, 0)
	return nil
}

func (dm *displayMonitorWindows) windowProcedure(hwnd windows.Handle, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case WM_CLOSE, WM_STOP_DISPLAY_MONITOR:
		destroyWindow(dm.windowHandle)
	case WM_DESTROY:
		postQuitMessage(0)
	case WM_DISPLAYCHANGE:
		if dm.resolutionChangeHandler == nil {
			break
		}
		e := ResolutionChangeEvent{
			Width:  int(lparam) & lowerWordMask,
			Height: (int(lparam) >> upperWordShift) & lowerWordMask,
		}
		dm.resolutionChangeHandler(e)
	case WM_WTSSESSION_CHANGE:
		if dm.sessionLockHandler == nil {
			break
		}
		switch wparam {
		case WTS_SESSION_LOCK:
			dm.sessionLockHandler(SessionLockEvent{ID: int(lparam), Locked: true})
		case WTS_SESSION_UNLOCK:
			dm.sessionLockHandler(SessionLockEvent{ID: int(lparam), Locked: false})
		}
	}

	return defWindowProc(hwnd, msg, wparam, lparam)
}

func getMessageLoop() error {
	for {
		msg := tMSG{}

		ret, err := getMessageW(&msg, 0, 0, 0)
		if err != nil {
			return fmt.Errorf("failed to get message: %v", err)
		}
		if !ret {
			return nil
		}
		dispatchMessage(&msg)
	}
}
