//go:build ignore

package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

func main() {
	makeAppIcon()
	makeShuffleIcon()
}

func makeAppIcon() {
	img := image.NewNRGBA(image.Rect(0, 0, 512, 512))
	bg := color.NRGBA{0x7c, 0x3a, 0xed, 0xff}
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// Head
	drawFilledCircle(img, 256, 200, 120, color.White)
	// Body
	drawFilledRect(img, 136, 200, 376, 400, color.White)
	// Bottom scallops
	drawFilledCircle(img, 176, 400, 40, bg)
	drawFilledCircle(img, 256, 400, 40, bg)
	drawFilledCircle(img, 336, 400, 40, bg)
	// Eyes
	drawFilledCircle(img, 216, 195, 22, bg)
	drawFilledCircle(img, 296, 195, 22, bg)

	f, _ := os.Create("assets/icon.png")
	defer f.Close()
	png.Encode(f, img)
}

func makeShuffleIcon() {
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	transparent := color.NRGBA{0, 0, 0, 0}
	draw.Draw(img, img.Bounds(), &image.Uniform{transparent}, image.Point{}, draw.Src)

	arrow := color.NRGBA{0xff, 0xff, 0xff, 0xff}

	// Two crossing arrows with arrowheads.
	// Arrow 1: from bottom-left to top-right.
	drawThickLine(img, 12, 48, 52, 16, 4, arrow)
	drawArrowHead(img, 52, 16, math.Pi/4, 8, 4, arrow)

	// Arrow 2: from top-left to bottom-right.
	drawThickLine(img, 12, 16, 52, 48, 4, arrow)
	drawArrowHead(img, 52, 48, -math.Pi/4, 8, 4, arrow)

	f, _ := os.Create("assets/shuffle.png")
	defer f.Close()
	png.Encode(f, img)
}

func drawFilledCircle(img *image.NRGBA, cx, cy, r int, c color.Color) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			dx, dy := x-cx, y-cy
			if dx*dx+dy*dy <= r*r {
				img.Set(x, y, c)
			}
		}
	}
}

func drawFilledRect(img *image.NRGBA, x0, y0, x1, y1 int, c color.Color) {
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			img.Set(x, y, c)
		}
	}
}

func drawThickLine(img *image.NRGBA, x0, y0, x1, y1, thickness int, c color.Color) {
	dx := x1 - x0
	dy := y1 - y0
	steps := int(math.Sqrt(float64(dx*dx+dy*dy))) * 2
	if steps < 1 {
		steps = 1
	}
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := int(float64(x0) + t*float64(dx))
		y := int(float64(y0) + t*float64(dy))
		drawFilledCircle(img, x, y, thickness/2, c)
	}
}

func drawArrowHead(img *image.NRGBA, x, y int, angle float64, length, thickness int, c color.Color) {
	// Two lines forming a V shape pointing in the given angle.
	for _, delta := range []float64{math.Pi / 6, -math.Pi / 6} {
		a := angle + delta
		x1 := int(float64(x) - float64(length)*math.Cos(a))
		y1 := int(float64(y) - float64(length)*math.Sin(a))
		drawThickLine(img, x, y, x1, y1, thickness, c)
	}
}
