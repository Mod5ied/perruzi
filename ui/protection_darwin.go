//go:build darwin

package ui

import (
	"fyne.io/fyne/v2"
)

/*
#cgo LDFLAGS: -framework Cocoa

void setWindowSharingType(void *windowPtr);
*/
import "C"

// SetContentProtection makes the Fyne window invisible to macOS screen capture
// (legacy window-capture APIs and most screensharing tools).
func SetContentProtection(w fyne.Window) {
	gw := getGLFWWindow(w)
	if gw == nil {
		return
	}
	C.setWindowSharingType(gw.GetCocoaWindow())
}

// SetWindowKeyFocusEnabled is a no-op on macOS in this build.
func SetWindowKeyFocusEnabled(w fyne.Window, enabled bool) {}
