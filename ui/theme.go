package ui

import (
	"image/color"

	"Peruzzi/assets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type peruzziTheme struct{}

// PeruzziTheme is the custom neubrutalist theme used by Peruzzi.
var PeruzziTheme fyne.Theme = &peruzziTheme{}

const (
	// Landing page palette.
	colorBg     = 0xFFFBF0
	colorBgAlt  = 0xF5F0E8
	colorBlack  = 0x0A0A0A
	colorTeal   = 0x00D4AA
	colorWhite  = 0xFFFFFF
	colorGray   = 0x6B7280
	colorYellow = 0xFFD60A
)

func c(rgb uint32) color.Color {
	return color.NRGBA{
		R: uint8(rgb >> 16),
		G: uint8(rgb >> 8),
		B: uint8(rgb),
		A: 0xff,
	}
}

var colours = map[fyne.ThemeColorName]color.Color{
	theme.ColorNameBackground:      c(colorBg),
	theme.ColorNameForeground:      c(colorBlack),
	theme.ColorNamePrimary:         c(colorTeal),
	theme.ColorNameButton:          c(colorTeal),
	theme.ColorNameInputBackground: c(colorWhite),
	theme.ColorNameDisabled:        c(colorGray),
	theme.ColorNamePlaceHolder:     c(colorGray),
	theme.ColorNameScrollBar:       c(colorBlack),
	theme.ColorNameShadow:          c(colorBlack),
	theme.ColorNameSelection:       c(colorTeal),
}

func (g *peruzziTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if c, ok := colours[n]; ok {
		return c
	}
	return theme.DefaultTheme().Color(n, v)
}

func (g *peruzziTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (g *peruzziTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return assets.FontMono()
	}
	return assets.FontBody()
}

func (g *peruzziTheme) Size(n fyne.ThemeSizeName) float32 {
	switch n {
	case theme.SizeNameText:
		return 13
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInnerPadding:
		return 6
	case theme.SizeNameInlineIcon:
		return 16
	}
	return theme.DefaultTheme().Size(n)
}
