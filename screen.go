package sketch

import (
	"image/color"

	"github.com/MichalMitros/sketch/vec"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Screen is the screen used to draw to.
type Screen struct {
	img          *ebiten.Image
	antyaliasing bool
	background   color.Color
}

func newScreen(img *ebiten.Image, background color.Color, antyaliasing bool) *Screen {
	return &Screen{
		img:          img,
		background:   background,
		antyaliasing: antyaliasing,
	}
}

// Width returns the width of the screen.
func (s *Screen) Width() float64 {
	return float64(s.img.Bounds().Dx())
}

// Height returns the height of the screen.
func (s *Screen) Height() float64 {
	return float64(s.img.Bounds().Dy())
}

// ScreenSize returns the width and height of the screen.
func (s *Screen) ScreenSize() (float64, float64) {
	return s.Width(), s.Height()
}

// Clear clears the screen.
func (s *Screen) Clear() {
	s.img.Fill(s.background)
}

// Fill fills the screen with the given color.
func (s *Screen) Fill(c color.Color) {
	s.img.Fill(c)
}

// At returns the color at the given position.
func (s *Screen) At(x, y int) color.Color {
	return s.img.At(x, y)
}

// Line draws a line from v1 to v2 with the given stroke width and color.
func (s *Screen) Line(v1, v2 vec.Vector, strokeWidth float32, c color.Color) {
	vector.StrokeLine(
		s.img,
		float32(v1.X), float32(v1.Y),
		float32(v2.X), float32(v2.Y),
		strokeWidth, c,
		s.antyaliasing,
	)
}

// Circle draws a circle at the given position with the given radius, stroke width and color.
func (s *Screen) Circle(v vec.Vector, radius, strokeWidth float32, c color.Color) {
	vector.StrokeCircle(
		s.img,
		float32(v.X), float32(v.Y),
		radius, strokeWidth, c,
		s.antyaliasing,
	)
}

// Rectangle draws a rectangle at the given position with the given width, height, stroke width and color.
func (s *Screen) Rectangle(v vec.Vector, width, height, strokeWidth float32, c color.Color) {
	vector.StrokeRect(
		s.img,
		float32(v.X), float32(v.Y),
		width, height, strokeWidth, c,
		s.antyaliasing,
	)
}

// FillCircle draws a filled circle at the given position with the given radius, stroke width and color.
func (s *Screen) FillCircle(v vec.Vector, radius, strokeWidth float32, c color.Color) {
	vector.FillCircle(
		s.img,
		float32(v.X), float32(v.Y),
		radius, c,
		s.antyaliasing,
	)
}

// FillRectangle draws a filled rectangle at the given position with the given width, height, stroke width and color.
func (s *Screen) FillRectangle(v vec.Vector, width, height, strokeWidth float32, c color.Color) {
	vector.FillRect(
		s.img,
		float32(v.X), float32(v.Y),
		width, height, c,
		s.antyaliasing,
	)
}
