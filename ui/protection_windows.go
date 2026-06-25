//go:build windows

package ui

import (
	"unsafe"

	"fyne.io/fyne/v2"
	"golang.org/x/sys/windows"
)

const WDA_EXCLUDEFROMCAPTURE = 0x00000011

var user32 = windows.NewLazySystemDLL("user32.dll")
var setWindowDisplayAffinity = user32.NewProc("SetWindowDisplayAffinity")

func nativeWindowPtr(w fyne.Window) unsafe.Pointer {
	gw := getGLFWWindow(w)
	if gw == nil {
		return nil
	}
	return unsafe.Pointer(gw.GetWin32Window())
}

// SetContentProtection makes the Fyne window invisible to Windows screen capture.
func SetContentProtection(w fyne.Window) {
	ptr := nativeWindowPtr(w)
	if ptr == nil {
		return
	}
	hwnd := uintptr(ptr)
	setWindowDisplayAffinity.Call(hwnd, WDA_EXCLUDEFROMCAPTURE)
}

// SetWindowKeyFocusEnabled is a no-op on Windows; focus handling differs there.
func SetWindowKeyFocusEnabled(w fyne.Window, enabled bool) {}
