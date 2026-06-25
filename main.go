package main

import (
	"os/exec"
	"runtime"

	"Peruzzi/keyboard"
	"Peruzzi/ui"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

func main() {
	runtime.LockOSThread()

	a := app.NewWithID("com.Peruzzi.app")
	a.Settings().SetTheme(ui.PeruzziTheme)

	w := ui.NewMainWindow(a)

	if runtime.GOOS == "darwin" {
		if !keyboard.IsAccessibilityGranted() {
			dialog.ShowInformation(
				"Permission Required",
				"Peruzzi needs Accessibility access to inject keystrokes.\n\nIf Peruzzi is already listed in System Settings > Privacy & Security > Accessibility, try this:\n1. Remove Peruzzi from the list (click the minus button).\n2. Click the plus button and add /Applications/Peruzzi.app.\n3. Make sure the toggle is ON.\n4. Quit and relaunch Peruzzi.\n\nIf you launched Peruzzi from Downloads or a zip file, move it to /Applications first.",
				w.Window,
			)
			exec.Command("open", "x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility").Start()
		}
	}

	keyboard.StartEscListener()
	defer keyboard.StopEscListener()

	endSystray := ui.StartSystray(w.Window.Show, a.Quit)
	defer endSystray()

	w.Window.ShowAndRun()
}
