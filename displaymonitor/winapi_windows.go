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
	pRegisterClassExW = user32.NewProc("RegisterClassExW")
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
		uintptr(0),
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

	ret, _, err := pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wcx)))

	if ret == 0 {
		return 0, err
	}
	return uint16(ret), nil
}

func destroyWindow(hwnd windows.Handle) error {
	ret, _, err := pDestroyWindow.Call(uintptr(hwnd))
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
