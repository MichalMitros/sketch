package sketch

import (
	"image/color"
	"math"

	"github.com/MichalMitros/sketch/internal"
	"github.com/MichalMitros/sketch/vector"
	"github.com/hajimehoshi/ebiten/v2"
	ebiten_vec "github.com/hajimehoshi/ebiten/v2/vector"
)

// Screen is the screen used to draw to.
type Screen struct {
	img          *ebiten.Image
	antyaliasing bool
	background   color.Color
	transforms   internal.TransformationStack
}

func newScreen(img *ebiten.Image, background color.Color, antyaliasing bool) *Screen {
	return &Screen{
		img:          img,
		background:   background,
		antyaliasing: antyaliasing,
		transforms:   internal.NewTransformationStack(),
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

// Scale scales both coordinate axes by rate.
func (s *Screen) Scale(rate float64) {
	s.transforms.Scale(rate, rate)
}

// ScaleX scales the X coordinate axis by rate.
func (s *Screen) ScaleX(rate float64) {
	s.transforms.Scale(rate, 1)
}

// ScaleY scales the Y coordinate axis by rate.
func (s *Screen) ScaleY(rate float64) {
	s.transforms.Scale(1, rate)
}

// Rotate rotates the coordinate system by angle radians.
func (s *Screen) Rotate(angle float64) {
	s.transforms.Rotate(angle)
}

// Transform moves the origin of the coordinate system by v.
func (s *Screen) Transform(v vector.Vector) {
	s.transforms.Translate(v.X, v.Y)
}

// TransformX moves the origin of the coordinate system along the X axis.
func (s *Screen) TransformX(dx float64) {
	s.transforms.Translate(dx, 0)
}

// TransformY moves the origin of the coordinate system along the Y axis.
func (s *Screen) TransformY(dy float64) {
	s.transforms.Translate(0, dy)
}

// Push adds an identity transformation layer to the stack.
func (s *Screen) Push() {
	s.transforms.Push()
}

// Pull removes the most recently pushed transformation layer.
// Pulling the initial layer has no effect.
func (s *Screen) Pull() {
	s.transforms.Pull()
}

// Line draws a line from v1 to v2 with the given stroke width and color.
func (s *Screen) Line(v1, v2 vector.Vector, strokeWidth float64, c color.Color) {
	var path ebiten_vec.Path
	path.MoveTo(float32(v1.X), float32(v1.Y))
	path.LineTo(float32(v2.X), float32(v2.Y))
	s.strokePath(&path, strokeWidth, c)
}

// Circle draws a circle at the given position with the given radius, stroke width and color.
func (s *Screen) Circle(v vector.Vector, radius, strokeWidth float64, c color.Color) {
	path := circlePath(v, radius)
	s.strokePath(&path, strokeWidth, c)
}

// Rectangle draws a rectangle at the given position with the given width, height, stroke width and color.
func (s *Screen) Rectangle(v vector.Vector, width, height, strokeWidth float64, c color.Color) {
	path := rectanglePath(v, width, height)
	s.strokePath(&path, strokeWidth, c)
}

// FillCircle draws a filled circle at the given position with the given radius, stroke width and color.
func (s *Screen) FillCircle(v vector.Vector, radius, strokeWidth float64, c color.Color) {
	path := circlePath(v, radius)
	s.fillPath(&path, c)
}

// FillRectangle draws a filled rectangle at the given position with the given width, height, stroke width and color.
func (s *Screen) FillRectangle(v vector.Vector, width, height, strokeWidth float64, c color.Color) {
	path := rectanglePath(v, width, height)
	s.fillPath(&path, c)
}

func circlePath(v vector.Vector, radius float64) ebiten_vec.Path {
	var path ebiten_vec.Path
	path.Arc(
		float32(v.X), float32(v.Y), float32(radius),
		0, 2*math.Pi, ebiten_vec.Clockwise,
	)
	path.Close()
	return path
}

func rectanglePath(v vector.Vector, width, height float64) ebiten_vec.Path {
	var path ebiten_vec.Path
	path.MoveTo(float32(v.X), float32(v.Y))
	path.LineTo(float32(v.X+width), float32(v.Y))
	path.LineTo(float32(v.X+width), float32(v.Y+height))
	path.LineTo(float32(v.X), float32(v.Y+height))
	path.Close()
	return path
}

func (s *Screen) transformedPath(path *ebiten_vec.Path) ebiten_vec.Path {
	var transformed ebiten_vec.Path
	transformed.AddPath(path, &ebiten_vec.AddPathOptions{
		GeoM: s.transforms.GeometryMatrix(),
	})
	return transformed
}

func (s *Screen) drawPathOptions(c color.Color) ebiten_vec.DrawPathOptions {
	var options ebiten_vec.DrawPathOptions
	options.AntiAlias = s.antyaliasing
	options.ColorScale.ScaleWithColor(c)
	return options
}

func (s *Screen) fillPath(path *ebiten_vec.Path, c color.Color) {
	transformed := s.transformedPath(path)
	options := s.drawPathOptions(c)
	ebiten_vec.FillPath(s.img, &transformed, nil, &options)
}

func (s *Screen) strokePath(path *ebiten_vec.Path, strokeWidth float64, c color.Color) {
	var transformed ebiten_vec.Path
	transformed.AddStroke(path, &ebiten_vec.AddStrokeOptions{
		StrokeOptions: ebiten_vec.StrokeOptions{Width: float32(strokeWidth)},
		GeoM:          s.transforms.GeometryMatrix(),
	})
	options := s.drawPathOptions(c)
	ebiten_vec.FillPath(s.img, &transformed, nil, &options)
}
