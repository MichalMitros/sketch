package sketch

import (
	"image/color"

	"github.com/MichalMitros/sketch/vector"
	"github.com/hajimehoshi/ebiten/v2"
	ebiten_vec "github.com/hajimehoshi/ebiten/v2/vector"
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
func (s *Screen) Line(v1, v2 vector.Vector, strokeWidth float64, c color.Color) {
	ebiten_vec.StrokeLine(
		s.img,
		float32(v1.X), float32(v1.Y),
		float32(v2.X), float32(v2.Y),
		float32(strokeWidth), c,
		s.antyaliasing,
	)
}

// Circle draws a circle at the given position with the given radius, stroke width and color.
func (s *Screen) Circle(v vector.Vector, radius, strokeWidth float64, c color.Color) {
	ebiten_vec.StrokeCircle(
		s.img,
		float32(v.X), float32(v.Y),
		float32(radius), float32(strokeWidth), c,
		s.antyaliasing,
	)
}

// Rectangle draws a rectangle at the given position with the given width, height, stroke width and color.
func (s *Screen) Rectangle(v vector.Vector, width, height, strokeWidth float64, c color.Color) {
	ebiten_vec.StrokeRect(
		s.img,
		float32(v.X), float32(v.Y),
		float32(width), float32(height), float32(strokeWidth), c,
		s.antyaliasing,
	)
}

// FillCircle draws a filled circle at the given position with the given radius, stroke width and color.
func (s *Screen) FillCircle(v vector.Vector, radius, strokeWidth float64, c color.Color) {
	ebiten_vec.FillCircle(
		s.img,
		float32(v.X), float32(v.Y),
		float32(radius), c,
		s.antyaliasing,
	)
}

// FillRectangle draws a filled rectangle at the given position with the given width, height, stroke width and color.
func (s *Screen) FillRectangle(v vector.Vector, width, height, strokeWidth float64, c color.Color) {
	ebiten_vec.FillRect(
		s.img,
		float32(v.X), float32(v.Y),
		float32(width), float32(height), c,
		s.antyaliasing,
	)
}
