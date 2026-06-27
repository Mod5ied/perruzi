package ui

import (
	"bytes"
	"image"
	"image/color"
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
	// No title on macOS menu bar; only the icon is shown so Peruzzi stays
	// discreet during screen sharing.
	systray.SetTitle("")
	systray.SetTooltip("Peruzzi")

	iconBytes := assets.Icon().Content()
	resized := resizeAndWhitenIcon(iconBytes, 22, 22)
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

// resizeAndWhitenIcon scales the source icon to the requested size and turns
// every non-transparent pixel white. This makes the menu-bar icon blend in
// with the macOS light/dark menu bar during screen sharing.
func resizeAndWhitenIcon(data []byte, width, height int) []byte {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	xdraw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), xdraw.Over, nil)

	// The source app icon has a dark background and a bright (teal) glyph.
	// Rescale it, then keep only the bright glyph pixels (turn them white)
	// and make everything else transparent so the menu-bar icon is just the
	// white ghost shape.
	white := color.RGBA{0xff, 0xff, 0xff, 0xff}
	transparent := color.RGBA{0, 0, 0, 0}
	const threshold uint32 = 10000 // out of 65535; separates glyph from dark bg.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := dst.At(x, y).RGBA()
			// RGBA values are in the range [0, 65535] here.
			avg := (r + g + b) / 3
			if avg > threshold {
				dst.Set(x, y, white)
			} else {
				dst.Set(x, y, transparent)
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil
	}
	return buf.Bytes()
}

// resizeIcon is an alias for resizeAndWhitenIcon (kept for backward compatibility).
func resizeIcon(data []byte, width, height int) []byte {
	return resizeAndWhitenIcon(data, width, height)
}
