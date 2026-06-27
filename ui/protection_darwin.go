//go:build darwin

package ui

import (
	"fyne.io/fyne/v2"
)

/*
#cgo LDFLAGS: -framework Cocoa

void setWindowSharingType(void *windowPtr);
void hideWindowFromCapture(void *windowPtr);
void showWindowForCapture(void *windowPtr);
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

// HideWindowForCapture aggressively hides the window from screen-capture
// pickers (e.g., Google Meet in Chrome/Edge) during typing.
func HideWindowForCapture(w fyne.Window) {
	gw := getGLFWWindow(w)
	if gw == nil {
		return
	}
	C.hideWindowFromCapture(gw.GetCocoaWindow())
}

// ShowWindowForCapture restores the window after typing is finished.
func ShowWindowForCapture(w fyne.Window) {
	gw := getGLFWWindow(w)
	if gw == nil {
		return
	}
	C.showWindowForCapture(gw.GetCocoaWindow())
}

// SetWindowKeyFocusEnabled is a no-op on macOS in this build.
func SetWindowKeyFocusEnabled(w fyne.Window, enabled bool) {}
