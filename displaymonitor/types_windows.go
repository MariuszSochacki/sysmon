package displaymonitor

import (
	"golang.org/x/sys/windows"
)

type tWNDCLASSEXW struct {
	cbSize      uint32
	cbStyle     uint32
	pWndProc    uintptr
	cbClsExtra  int32
	cbWndExtra  int32
	hInstance   windows.Handle
	hIcon       windows.Handle
	hCursor     windows.Handle
	hBackground windows.Handle
	pMenuName   *uint16
	pClassName  *uint16
	hSmallIcon  windows.Handle
}
