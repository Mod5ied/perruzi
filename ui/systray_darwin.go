//go:build darwin

package ui

import "fyne.io/systray"

// StartSystray integrates the system tray with Fyne's existing NSApplication run loop.
// It returns a cleanup function to call before the application exits.
func StartSystray(show func(), quit func()) func() {
	start, end := systray.RunWithExternalLoop(
		func() { InitTray(show, quit) },
		func() {},
	)
	start()
	return end
}
