package ui

import (
	"Peruzzi/assets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// HumaniseButton is a square toggle button with a Material icon and a hover tooltip.
type HumaniseButton struct {
	widget.Button
	active  bool
	tooltip *widget.PopUp
}

// NewHumaniseButton creates the Humanize toggle button.
func NewHumaniseButton() *HumaniseButton {
	h := &HumaniseButton{active: true}
	h.ExtendBaseWidget(h)
	h.SetIcon(assets.AutoFixHighIcon())
	h.updateStyle()
	return h
}

func (h *HumaniseButton) Tapped(_ *fyne.PointEvent) {
	h.active = !h.active
	h.updateStyle()
}

func (h *HumaniseButton) updateStyle() {
	if h.active {
		h.Importance = widget.HighImportance
	} else {
		h.Importance = widget.LowImportance
	}
	h.Refresh()
}

// IsActive returns the toggle state.
func (h *HumaniseButton) IsActive() bool {
	return h.active
}

// MouseIn shows the tooltip when the cursor enters the button.
func (h *HumaniseButton) MouseIn(e *desktop.MouseEvent) {
	h.Button.MouseIn(e)
	if h.tooltip == nil {
		label := widget.NewLabel("Humanize mode")
		label.TextStyle = fyne.TextStyle{Monospace: true}
		h.tooltip = widget.NewPopUp(label, fyne.CurrentApp().Driver().CanvasForObject(h))
	}

	btnPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(h)
	size := h.tooltip.MinSize()
	x := btnPos.X + h.Size().Width/2 - size.Width/2
	y := btnPos.Y - size.Height - 4
	h.tooltip.Move(fyne.NewPos(x, y))
	h.tooltip.Show()
}

// MouseOut hides the tooltip.
func (h *HumaniseButton) MouseOut() {
	h.Button.MouseOut()
	if h.tooltip != nil {
		h.tooltip.Hide()
	}
}

// MouseMoved is required by the desktop.Hoverable interface.
func (h *HumaniseButton) MouseMoved(_ *desktop.MouseEvent) {}

// Ensure HumaniseButton implements desktop.Hoverable.
var _ desktop.Hoverable = (*HumaniseButton)(nil)
