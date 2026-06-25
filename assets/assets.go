package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed icon.png
var iconPNG []byte

//go:embed auto_fix_high.svg
var autoFixHighSVG []byte

//go:embed check.svg
var checkSVG []byte

//go:embed fonts/SpaceGrotesk-Variable.ttf
var spaceGrotesk []byte

//go:embed fonts/JetBrainsMono-Variable.ttf
var jetBrainsMono []byte

// Icon returns the application tray / menu bar icon.
func Icon() fyne.Resource {
	return fyne.NewStaticResource("icon.png", iconPNG)
}

// AutoFixHighIcon returns the Material auto_fix_high SVG icon.
func AutoFixHighIcon() fyne.Resource {
	return fyne.NewStaticResource("auto_fix_high.svg", autoFixHighSVG)
}

// CheckIcon returns a simple checkmark SVG icon.
func CheckIcon() fyne.Resource {
	return fyne.NewStaticResource("check.svg", checkSVG)
}

// FontBody returns the Space Grotesk variable font.
func FontBody() fyne.Resource {
	return fyne.NewStaticResource("SpaceGrotesk-Variable.ttf", spaceGrotesk)
}

// FontMono returns the JetBrains Mono variable font.
func FontMono() fyne.Resource {
	return fyne.NewStaticResource("JetBrainsMono-Variable.ttf", jetBrainsMono)
}
