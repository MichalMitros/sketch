package sketch

import (
	"image/color"

	"github.com/MichalMitros/sketch/vec"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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

func (s *Screen) Width() float64 {
	return float64(s.img.Bounds().Dx())
}

func (s *Screen) Height() float64 {
	return float64(s.img.Bounds().Dy())
}

func (s *Screen) ScreenSize() (float64, float64) {
	return s.Width(), s.Height()
}

func (s *Screen) Clear() {
	s.img.Fill(s.background)
}

func (s *Screen) Fill(c color.Color) {
	s.img.Fill(c)
}

func (s *Screen) At(x, y int) color.Color {
	return s.img.At(x, y)
}

func (s *Screen) Line(v1, v2 vec.Vector, strokeWidth float32, c color.Color) {
	vector.StrokeLine(
		s.img,
		float32(v1.X), float32(v1.Y),
		float32(v2.X), float32(v2.Y),
		strokeWidth, c,
		s.antyaliasing,
	)
}

func (s *Screen) Circle(v vec.Vector, radius, strokeWidth float32, c color.Color) {
	vector.StrokeCircle(
		s.img,
		float32(v.X), float32(v.Y),
		radius, strokeWidth, c,
		s.antyaliasing,
	)
}

func (s *Screen) Rectangle(v vec.Vector, width, height, strokeWidth float32, c color.Color) {
	vector.StrokeRect(
		s.img,
		float32(v.X), float32(v.Y),
		width, height, strokeWidth, c,
		s.antyaliasing,
	)
}

func (s *Screen) FillCircle(v vec.Vector, radius, strokeWidth float32, c color.Color) {
	vector.FillCircle(
		s.img,
		float32(v.X), float32(v.Y),
		radius, c,
		s.antyaliasing,
	)
}

func (s *Screen) FillRectangle(v vec.Vector, width, height, strokeWidth float32, c color.Color) {
	vector.FillRect(
		s.img,
		float32(v.X), float32(v.Y),
		width, height, c,
		s.antyaliasing,
	)
}
