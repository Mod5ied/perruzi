package ui

import (
	"image/color"

	"Peruzzi/assets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// HumaniseCheck is a square neubrutalist checkbox used for the Humanise toggle.
type HumaniseCheck struct {
	widget.BaseWidget
	active   bool
	disabled bool
	OnToggle func(active bool)
}

// NewHumaniseCheck creates a new HumaniseCheck widget.
func NewHumaniseCheck(active bool, onToggle func(bool)) *HumaniseCheck {
	h := &HumaniseCheck{active: active, OnToggle: onToggle}
	h.ExtendBaseWidget(h)
	return h
}

// Tapped toggles the checkbox state.
func (h *HumaniseCheck) Tapped(_ *fyne.PointEvent) {
	if h.disabled {
		return
	}
	h.active = !h.active
	if h.OnToggle != nil {
		h.OnToggle(h.active)
	}
	h.Refresh()
}

// Disable prevents the checkbox from being toggled.
func (h *HumaniseCheck) Disable() {
	h.disabled = true
	h.Refresh()
}

// Enable allows the checkbox to be toggled.
func (h *HumaniseCheck) Enable() {
	h.disabled = false
	h.Refresh()
}

// IsActive returns the current toggle state.
func (h *HumaniseCheck) IsActive() bool {
	return h.active
}

// CreateRenderer builds the checkbox renderer.
func (h *HumaniseCheck) CreateRenderer() fyne.WidgetRenderer {
	border := canvas.NewRectangle(neuWhite)
	border.StrokeColor = neuBlack
	border.StrokeWidth = 2

	check := canvas.NewImageFromResource(assets.CheckIcon())
	check.FillMode = canvas.ImageFillOriginal

	return &humaniseCheckRenderer{
		h:     h,
		border: border,
		check:  check,
	}
}

type humaniseCheckRenderer struct {
	h      *HumaniseCheck
	border *canvas.Rectangle
	check  *canvas.Image
}

func (r *humaniseCheckRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.border, r.check}
}

func (r *humaniseCheckRenderer) Layout(size fyne.Size) {
	r.border.Resize(size)
	r.border.Move(fyne.NewPos(0, 0))

	iconSize := r.check.MinSize()
	if iconSize.Width == 0 || iconSize.Height == 0 {
		iconSize = fyne.NewSize(size.Width*0.55, size.Height*0.55)
	}
	// Scale icon down slightly and keep it centred.
	iconSize.Width *= 0.8
	iconSize.Height *= 0.8
	r.check.Resize(iconSize)
	r.check.Move(fyne.NewPos(
		(size.Width-iconSize.Width)/2,
		(size.Height-iconSize.Height)/2,
	))
}

func (r *humaniseCheckRenderer) MinSize() fyne.Size {
	return fyne.NewSize(16, 16)
}

func (r *humaniseCheckRenderer) Refresh() {
	if r.h.disabled {
		r.border.FillColor = color.NRGBA{0xE5, 0xE0, 0xD8, 0xFF}
	} else if r.h.active {
		r.border.FillColor = neuTeal
	} else {
		r.border.FillColor = neuWhite
	}
	if r.h.active && !r.h.disabled {
		r.check.Show()
	} else {
		r.check.Hide()
	}
	r.border.Refresh()
	r.check.Refresh()
	r.Layout(r.h.Size())
}

func (r *humaniseCheckRenderer) Destroy() {}
