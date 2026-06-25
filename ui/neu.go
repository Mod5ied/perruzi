package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	neuBorderWidth = 3
	neuShadowOff   = 4
	neuPadding     = 6
)

var (
	neuBlack = color.NRGBA{0x0A, 0x0A, 0x0A, 0xFF}
	neuWhite = color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}
	neuTeal  = color.NRGBA{0x00, 0xD4, 0xAA, 0xFF}
)

// NeuButton is a neubrutalist pressable button with a hard shadow and thick border.
type NeuButton struct {
	widget.BaseWidget
	Text     string
	Icon     fyne.Resource
	Fill     color.Color
	Shadow   color.Color
	OnTapped func()
	hovered  bool
	pressed  bool
	disabled bool
}

// NewNeuButton creates a neubrutalist button.
func NewNeuButton(text string, icon fyne.Resource, fill, shadow color.Color, tapped func()) *NeuButton {
	b := &NeuButton{Text: text, Icon: icon, Fill: fill, Shadow: shadow, OnTapped: tapped}
	b.ExtendBaseWidget(b)
	return b
}

// Tapped handles button clicks.
func (b *NeuButton) Tapped(_ *fyne.PointEvent) {
	if b.disabled {
		return
	}
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// Disable greys out the button and prevents taps.
func (b *NeuButton) Disable() {
	b.disabled = true
	b.hovered = false
	b.pressed = false
	b.Refresh()
}

// Enable restores normal button behaviour.
func (b *NeuButton) Enable() {
	b.disabled = false
	b.Refresh()
}

// MouseIn handles hover start.
func (b *NeuButton) MouseIn(_ *desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

// MouseOut handles hover end.
func (b *NeuButton) MouseOut() {
	b.hovered = false
	b.pressed = false
	b.Refresh()
}

// MouseDown handles press start.
func (b *NeuButton) MouseDown(_ *desktop.MouseEvent) {
	b.pressed = true
	b.Refresh()
}

// MouseUp handles press end.
func (b *NeuButton) MouseUp(_ *desktop.MouseEvent) {
	b.pressed = false
	b.Refresh()
}

// MouseMoved is required by desktop.Hoverable.
func (b *NeuButton) MouseMoved(_ *desktop.MouseEvent) {}

var _ fyne.Tappable = (*NeuButton)(nil)
var _ desktop.Hoverable = (*NeuButton)(nil)
var _ desktop.Mouseable = (*NeuButton)(nil)

// CreateRenderer builds the button renderer.
func (b *NeuButton) CreateRenderer() fyne.WidgetRenderer {
	shadow := canvas.NewRectangle(b.Shadow)
	border := canvas.NewRectangle(b.Fill)
	border.StrokeColor = neuBlack
	border.StrokeWidth = neuBorderWidth

	style := fyne.TextStyle{Bold: true}
	label := widget.NewLabelWithStyle(b.Text, fyne.TextAlignCenter, style)
	label.TextStyle = style

	var icon *canvas.Image
	if b.Icon != nil {
		icon = canvas.NewImageFromResource(b.Icon)
		icon.FillMode = canvas.ImageFillOriginal
	}

	return &neuButtonRenderer{
		b:      b,
		shadow: shadow,
		border: border,
		label:  label,
		icon:   icon,
	}
}

type neuButtonRenderer struct {
	b      *NeuButton
	shadow *canvas.Rectangle
	border *canvas.Rectangle
	label  *widget.Label
	icon   *canvas.Image
}

func (r *neuButtonRenderer) Objects() []fyne.CanvasObject {
	if r.icon != nil {
		return []fyne.CanvasObject{r.shadow, r.border, r.icon, r.label}
	}
	return []fyne.CanvasObject{r.shadow, r.border, r.label}
}

func (r *neuButtonRenderer) Layout(size fyne.Size) {
	off := float32(neuShadowOff)
	var pressX, pressY float32
	if r.b.disabled {
		pressX, pressY = 0, 0
		off = 0
	} else if r.b.pressed {
		pressX, pressY = 4, 4
		off = 0
	} else if r.b.hovered {
		pressX, pressY = 2, 2
		off = 2
	}

	boxW := size.Width - off
	boxH := size.Height - off

	r.shadow.Resize(fyne.NewSize(boxW, boxH))
	r.shadow.Move(fyne.NewPos(off, off))

	r.border.Resize(fyne.NewSize(boxW, boxH))
	r.border.Move(fyne.NewPos(pressX, pressY))

	labelSize := r.label.MinSize()
	if r.icon != nil {
		iconSize := r.icon.MinSize()
		if iconSize.Width == 0 || iconSize.Height == 0 {
			iconSize = fyne.NewSize(16, 16)
		}
		totalW := iconSize.Width + labelSize.Width + 4
		startX := pressX + (boxW-totalW)/2
		centerY := pressY + (boxH-iconSize.Height)/2
		r.icon.Resize(iconSize)
		r.icon.Move(fyne.NewPos(startX, centerY))
		r.label.Resize(labelSize)
		r.label.Move(fyne.NewPos(startX+iconSize.Width+4, pressY+(boxH-labelSize.Height)/2))
	} else {
		r.label.Resize(labelSize)
		r.label.Move(fyne.NewPos(pressX+(boxW-labelSize.Width)/2, pressY+(boxH-labelSize.Height)/2))
	}
}

func (r *neuButtonRenderer) MinSize() fyne.Size {
	labelSize := r.label.MinSize()
	if r.icon != nil {
		iconSize := r.icon.MinSize()
		if iconSize.Width == 0 || iconSize.Height == 0 {
			iconSize = fyne.NewSize(16, 16)
		}
		if r.b.Text == "" {
			return fyne.NewSize(max(iconSize.Width+12, 32), max(iconSize.Height+12, 32))
		}
		w := max(labelSize.Width+iconSize.Width+16, 64)
		h := max(max(labelSize.Height, iconSize.Height)+16, 36)
		return fyne.NewSize(w, h)
	}
	return fyne.NewSize(max(labelSize.Width+16, 64), max(labelSize.Height+16, 36))
}

func (r *neuButtonRenderer) Refresh() {
	if r.b.disabled {
		r.border.FillColor = color.NRGBA{0xE5, 0xE0, 0xD8, 0xFF}
		r.shadow.FillColor = color.NRGBA{0xE5, 0xE0, 0xD8, 0xFF}
	} else {
		r.border.FillColor = r.b.Fill
		r.shadow.FillColor = r.b.Shadow
	}
	r.border.Refresh()
	r.shadow.Refresh()
	r.label.Text = r.b.Text
	r.label.Refresh()
	r.Layout(r.b.Size())
}

func (r *neuButtonRenderer) Destroy() {}

// NeuBox is a neubrutalist container with a hard shadow, thick border and a child widget.
type NeuBox struct {
	widget.BaseWidget
	Child  fyne.CanvasObject
	Fill   color.Color
	Shadow color.Color
}

// NewNeuBox creates a neubrutalist container.
func NewNeuBox(child fyne.CanvasObject, fill, shadow color.Color) *NeuBox {
	b := &NeuBox{Child: child, Fill: fill, Shadow: shadow}
	b.ExtendBaseWidget(b)
	return b
}

// CreateRenderer builds the box renderer.
func (b *NeuBox) CreateRenderer() fyne.WidgetRenderer {
	shadow := canvas.NewRectangle(b.Shadow)
	border := canvas.NewRectangle(b.Fill)
	border.StrokeColor = neuBlack
	border.StrokeWidth = neuBorderWidth
	return &neuBoxRenderer{b: b, shadow: shadow, border: border, child: b.Child}
}

type neuBoxRenderer struct {
	b      *NeuBox
	shadow *canvas.Rectangle
	border *canvas.Rectangle
	child  fyne.CanvasObject
}

func (r *neuBoxRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.shadow, r.border, r.child}
}

func (r *neuBoxRenderer) Layout(size fyne.Size) {
	off := float32(neuShadowOff)
	pad := float32(neuPadding)

	boxW := size.Width - off
	boxH := size.Height - off

	r.shadow.Resize(fyne.NewSize(boxW, boxH))
	r.shadow.Move(fyne.NewPos(off, off))

	r.border.Resize(fyne.NewSize(boxW, boxH))
	r.border.Move(fyne.NewPos(0, 0))

	r.child.Resize(fyne.NewSize(boxW-2*pad, boxH-2*pad))
	r.child.Move(fyne.NewPos(pad, pad))
}

func (r *neuBoxRenderer) MinSize() fyne.Size {
	childSize := r.child.MinSize()
	pad := float32(neuPadding)
	off := float32(neuShadowOff)
	return fyne.NewSize(childSize.Width+2*pad+off, childSize.Height+2*pad+off)
}

func (r *neuBoxRenderer) Refresh() {
	r.border.FillColor = r.b.Fill
	r.shadow.FillColor = r.b.Shadow
	r.border.Refresh()
	r.shadow.Refresh()
	r.Layout(r.b.Size())
}

func (r *neuBoxRenderer) Destroy() {}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
