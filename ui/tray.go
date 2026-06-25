package ui

import (
	"bytes"
	"image"
	"image/png"
	"os"

	"Peruzzi/assets"

	"fyne.io/systray"
	xdraw "golang.org/x/image/draw"
)

// exitFunc is swapped at runtime so tests or callers can override it.
var exitFunc = os.Exit

// InitTray sets up the system tray / menu bar icon and menu.
func InitTray(show func(), quit func()) {
	systray.SetTitle("Peruzzi")
	systray.SetTooltip("Peruzzi")

	iconBytes := assets.Icon().Content()
	resized := resizeIcon(iconBytes, 22, 22)
	if resized != nil {
		systray.SetIcon(resized)
	} else {
		systray.SetIcon(iconBytes)
	}

	showItem := systray.AddMenuItem("Open", "Open Peruzzi")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "Quit Peruzzi")

	go func() {
		for {
			select {
			case <-showItem.ClickedCh:
				show()
			case <-quitItem.ClickedCh:
				quit()
				systray.Quit()
				exitFunc(0)
				return
			}
		}
	}()
}

func resizeIcon(data []byte, width, height int) []byte {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	xdraw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), xdraw.Over, nil)

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil
	}
	return buf.Bytes()
}
