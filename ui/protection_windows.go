//go:build windows

package ui

import (
	"unsafe"

	"fyne.io/fyne/v2"
	"golang.org/x/sys/windows"
)

const WDA_EXCLUDEFROMCAPTURE = 0x00000011
const WDA_NONE = 0x00000000

var user32 = windows.NewLazySystemDLL("user32.dll")
var setWindowDisplayAffinity = user32.NewProc("SetWindowDisplayAffinity")
var showWindow = user32.NewProc("ShowWindow")

const (
	SW_HIDE = 0
	SW_SHOW = 5
)

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

// HideWindowForCapture hides the window from screen capture and the taskbar
// switcher during typing.
func HideWindowForCapture(w fyne.Window) {
	ptr := nativeWindowPtr(w)
	if ptr == nil {
		return
	}
	hwnd := uintptr(ptr)
	// Keep display affinity active (window excluded from capture even if visible),
	// then hide the actual window.
	setWindowDisplayAffinity.Call(hwnd, WDA_EXCLUDEFROMCAPTURE)
	showWindow.Call(hwnd, SW_HIDE)
}

// ShowWindowForCapture restores the window after typing is finished.
func ShowWindowForCapture(w fyne.Window) {
	ptr := nativeWindowPtr(w)
	if ptr == nil {
		return
	}
	hwnd := uintptr(ptr)
	showWindow.Call(hwnd, SW_SHOW)
	setWindowDisplayAffinity.Call(hwnd, WDA_EXCLUDEFROMCAPTURE)
}

// SetWindowKeyFocusEnabled is a no-op on Windows; focus handling differs there.
func SetWindowKeyFocusEnabled(w fyne.Window, enabled bool) {}
