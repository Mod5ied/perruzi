//go:build windows

package ui

import "fyne.io/systray"

// StartSystray runs the system tray in its own goroutine.
// It returns a cleanup function to call before the application exits.
func StartSystray(show func(), quit func()) func() {
	go systray.Run(
		func() { InitTray(show, quit) },
		func() {},
	)
	return func() { systray.Quit() }
}
