package displaymonitor

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32 = windows.NewLazySystemDLL("user32.dll")

	pCreateWindowExW  = user32.NewProc("CreateWindowExW")
	pDefWindowProcW   = user32.NewProc("DefWindowProcW")
	pDestroyWindow    = user32.NewProc("DestroyWindow")
	pDispatchMessageW = user32.NewProc("DispatchMessageW")
	pGetMessageW      = user32.NewProc("GetMessageW")
	pSendMessageW     = user32.NewProc("SendMessageW")
	pRegisterClassExW = user32.NewProc("RegisterClassExW")
	pGetDesktopWindow = user32.NewProc("GetDesktopWindow")
	pPostQuitMessage  = user32.NewProc("PostQuitMessage")
)

var (
	wtsapi32                        = windows.NewLazySystemDLL("Wtsapi32.dll")
	pWTSRegisterSessionNotification = wtsapi32.NewProc("WTSRegisterSessionNotification")
)

func createWindowExW(clsName, wndName string, instance windows.Handle) (windows.Handle, error) {
	ret, _, err := pCreateWindowExW.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(clsName))),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(wndName))),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(getDesktopWindow()),
		uintptr(0),
		uintptr(instance),
		uintptr(0),
	)

	if ret == 0 {
		return 0, err
	}

	return windows.Handle(ret), nil
}

func registerClassEx(clsName string, wndProc interface{}, instance windows.Handle) (uint16, error) {
	wcx := tWNDCLASSEXW{
		pWndProc:   windows.NewCallback(wndProc),
		hInstance:  instance,
		pClassName: windows.StringToUTF16Ptr(clsName),
	}
	wcx.cbSize = uint32(unsafe.Sizeof(wcx))

	ret, _, err := pRegisterClassExW.Call(
		uintptr(unsafe.Pointer(&wcx)),
	)

	if ret == 0 {
		return 0, err
	}

	return uint16(ret), nil
}

func destroyWindow(hwnd windows.Handle) error {
	ret, _, err := pDestroyWindow.Call(
		uintptr(hwnd))

	if ret == 0 {
		return err
	}

	return nil
}

func defWindowProc(hwnd windows.Handle, msg uint32, wparam, lparam uintptr) uintptr {
	ret, _, _ := pDefWindowProcW.Call(
		uintptr(hwnd),
		uintptr(msg),
		uintptr(wparam),
		uintptr(lparam),
	)

	return uintptr(ret)
}

func getMessageW(msg *tMSG, hwnd windows.Handle, msgFilterMin, msgFilterMax uint32) (bool, error) {
	ret, _, err := pGetMessageW.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)

	if int32(ret) == -1 {
		return false, err
	}

	return int32(ret) != 0, nil
}

func dispatchMessage(msg *tMSG) {
	pDispatchMessageW.Call(
		uintptr(unsafe.Pointer(msg)),
	)
}

func getDesktopWindow() windows.Handle {
	ret, _, _ := pGetDesktopWindow.Call()
	return windows.Handle(ret)
}

func wtsRegisterSessionNotification(hwnd windows.Handle, dwFlags uint32) error {
	ret, _, err := pWTSRegisterSessionNotification.Call(
		uintptr(hwnd),
		uintptr(dwFlags),
	)

	if ret == 0 {
		return err
	}

	return nil
}

func postQuitMessage(exitCode int32) {
	pPostQuitMessage.Call(
		uintptr(exitCode),
	)
}

func sendMessageW(hwnd windows.Handle, msg, lparam, wparam uint32) uintptr {
	ret, _, _ := pSendMessageW.Call(
		uintptr(hwnd),
		uintptr(msg),
		uintptr(lparam),
		uintptr(wparam),
	)

	return ret
}
