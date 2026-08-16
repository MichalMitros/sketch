package sketch

import (
	"image"
	"image/color"
	"io/fs"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	ebiten_vec "github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/michalmitros/sketch/vector"
)

// Image represents a drawable 2D image.
type Image struct {
	img *ebiten.Image
}

// ImageFromFile reads an image file from the given path and returns it as an *Image.
// Supports GIF, JPEG, PNG.
func ImageFromFile(path string) (*Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}
	return &Image{img: img}, nil
}

// ImageFromFS reads an image file from the given filesystem path and returns it as an *Image.
// Supports GIF, JPEG, PNG.
func ImageFromFS(fs fs.FS, path string) (*Image, error) {
	img, _, err := ebitenutil.NewImageFromFileSystem(fs, path)
	if err != nil {
		return nil, err
	}
	return &Image{img: img}, nil
}

// ImageFromStdImage reads an image from the given io.Reader and returns it as an *Image.
// Supports GIF, JPEG, PNG.
func ImageFromStdImage(src image.Image) *Image {
	return &Image{img: ebiten.NewImageFromImage(src)}
}

func fromEbitenImage(img *ebiten.Image) *Image {
	return &Image{img: img}
}

func (i *Image) toEbitenImage() *ebiten.Image {
	return i.img
}

// NewBlankImage creates a new blank (transparent) image with the given dimensions.
// Width and height must be positive.
func NewBlankImage(dim vector.Vector) *Image {
	return &Image{
		img: ebiten.NewImage(int(dim.X), int(dim.Y)),
	}
}

// NewFilledImage creates a new image with the given dimensions, filled with the given color.
// Width and height must be positive.
func NewFilledImage(dim vector.Vector, c color.Color) *Image {
	img := NewBlankImage(dim)
	img.Fill(c)
	return img
}

// NewWhiteImage creates a new image with the given dimensions, filled with white.
// Width and height must be positive.
func NewWhiteImage(dim vector.Vector) *Image {
	return NewFilledImage(dim, color.White)
}

// NewBlackImage creates a new image with the given dimensions, filled with black.
// Width and height must be positive.
func NewBlackImage(dim vector.Vector) *Image {
	return NewFilledImage(dim, color.Black)
}

// Width returns the width of the image in pixels.
func (i *Image) Width() float64 {
	return float64(i.img.Bounds().Dx())
}

// Height returns the height of the image in pixels.
func (i *Image) Height() float64 {
	return float64(i.img.Bounds().Dy())
}

// Size returns the dimensions of the image as a vector (width, height).
func (i *Image) Size() vector.Vector {
	return vector.New(float64(i.Width()), float64(i.Height()))
}

// Clone returns a new Image that is an independent pixel-by-pixel copy
// of this image. Changes to the clone do not affect the original.
func (i *Image) Clone() *Image {
	clone := NewBlankImage(i.Size())
	clone.img.DrawImage(i.img, nil)
	return clone
}

// DrawOptions controls how an image is drawn.
// Zero values result in sensible defaults: position at (0,0), original scale,
// top-left anchor, no rotation, full opacity, no tint.
type DrawOptions struct {
	// Pos is the position where the image is drawn (defaults to origin).
	Pos vector.Vector
	// Anchor is the anchor point within the image.
	// (0,0) means top-left, (0.5,0.5) means center, (1,1) means bottom-right.
	// Defaults to (0,0) – top-left.
	Anchor vector.Vector
	// Scale factors for X and Y axes. (1,1) = original size.
	// Negative values flip the image. Defaults to (1,1).
	Scale vector.Vector
	// Rotation in radians around the anchor point (counter-clockwise).
	// Defaults to 0.
	Rotation float64
	// Tint is an optional color multiplier. A value of white means no tint
	// (original colors preserved). nil means no tint.
	Tint color.Color
	// Opacity controls transparency: 0 = fully transparent, 1 = fully opaque.
	// Defaults to 1.
	Opacity float64
}

// DefaultDrawOptions returns draw options with sensible defaults:
// position at origin, original scale, top-left anchor, no rotation, full opacity, no tint.
func DefaultDrawOptions() DrawOptions {
	return DrawOptions{
		Scale:   vector.New(1, 1),
		Opacity: 1,
	}
}

// Draw draws the image onto the destination image using the given options.
func (i *Image) Draw(dst *Image, opts *DrawOptions) {
	if opts == nil {
		o := DefaultDrawOptions()
		opts = &o
	}

	// Skip drawing if invisible.
	if opts.Scale.X == 0 || opts.Scale.Y == 0 || opts.Opacity <= 0 {
		return
	}

	var eo ebiten.DrawImageOptions

	w, h := i.Size().Values()
	eo.GeoM.Translate(-w*opts.Anchor.X, -h*opts.Anchor.Y)
	eo.GeoM.Scale(opts.Scale.X, opts.Scale.Y)

	if opts.Rotation != 0 {
		eo.GeoM.Rotate(opts.Rotation)
	}

	eo.GeoM.Translate(opts.Pos.X, opts.Pos.Y)

	if opts.Tint != nil {
		eo.ColorScale.ScaleWithColor(opts.Tint)
	}

	if opts.Opacity < 1 {
		eo.ColorScale.ScaleAlpha(float32(opts.Opacity))
	}

	dst.img.DrawImage(i.img, &eo)
}

// DrawAt draws the image at the given position of the destination image.
func (i *Image) DrawAt(dst *Image, pos vector.Vector) {
	opts := DefaultDrawOptions()
	opts.Pos = pos
	i.Draw(dst, &opts)
}

// At returns the color of the pixel at the given coordinates.
func (i *Image) At(v vector.Vector) color.Color {
	return i.img.At(int(v.X), int(v.Y))
}

// Set sets the color of the pixel at the given coordinates.
func (i *Image) Set(v vector.Vector, c color.Color) {
	i.img.Set(int(v.X), int(v.Y), c)
}

// Fill fills the entire image with the given color.
func (i *Image) Fill(c color.Color) {
	i.img.Fill(c)
}

// Clear clears the image to fully transparent black.
func (i *Image) Clear() {
	i.img.Clear()
}

// SubImage returns a new Image that represents a rectangular sub-region
// of the original image. The returned image shares the same underlying
// pixel data - modifying pixels in one affects the other.
// Use Clone() if you need an independent copy.
func (i *Image) SubImage(pos, size vector.Vector) *Image {
	r := image.Rect(int(pos.X), int(pos.Y), int(pos.X+size.X), int(pos.Y+size.Y))
	subimage, _ := i.img.SubImage(r).(*ebiten.Image)
	return &Image{
		img: subimage,
	}
}

// CopyTo copies a rectangular region from this image onto another image.
// srcX, srcY specify the top-left corner of the source rectangle in this image.
// dstX, dstY specify where to place it on the destination image.
func (i *Image) CopyTo(dst *Image, srcPos, size, dstPos vector.Vector) {
	src := i.SubImage(srcPos, size)
	opts := DefaultDrawOptions()
	opts.Pos = dstPos
	src.Draw(dst, &opts)
}

// Line draws a line from v1 to v2 with the given stroke width and color.
func (i *Image) Line(v1, v2 vector.Vector, strokeWidth float64, c color.Color) {
	var path ebiten_vec.Path
	path.MoveTo(float32(v1.X), float32(v1.Y))
	path.LineTo(float32(v2.X), float32(v2.Y))
	i.strokePath(&path, strokeWidth, c)
}

// Rectangle draws a rectangle at the given position with the given width, height, stroke width and color.
func (i *Image) Rectangle(v vector.Vector, width, height, strokeWidth float64, c color.Color) {
	path := rectanglePath(v, width, height)
	i.strokePath(&path, strokeWidth, c)
}

// FillRectangle draws a filled rectangle at the given position with the given width, height, stroke width and color.
func (i *Image) FillRectangle(v vector.Vector, width, height, strokeWidth float64, c color.Color) {
	path := rectanglePath(v, width, height)
	i.fillPath(&path, c)
}

// Circle draws a circle at the given position with the given radius, stroke width and color.
func (i *Image) Circle(v vector.Vector, radius, strokeWidth float64, c color.Color) {
	path := circlePath(v, radius)
	i.strokePath(&path, strokeWidth, c)
}

// FillCircle draws a filled circle at the given position with the given radius, stroke width and color.
func (i *Image) FillCircle(v vector.Vector, radius, strokeWidth float64, c color.Color) {
	path := circlePath(v, radius)
	i.fillPath(&path, c)
}

// Arc draws an arc at the given position with the given radius, start angle, end angle, stroke width and color.
func (i *Image) Arc(v vector.Vector, radius, startAngle, endAngle, strokeWidth float64, c color.Color) {
	path := arcPath(v, radius, startAngle, endAngle)
	i.strokePath(&path, strokeWidth, c)
}

// FillArc draws a filled arc (pie slice) at the given position with the given radius, start angle, end angle and color.
func (i *Image) FillArc(v vector.Vector, radius, startAngle, endAngle float64, c color.Color) {
	path := fillArcPath(v, radius, startAngle, endAngle)
	i.fillPath(&path, c)
}

// Shape draws a polygon through the given points with the given stroke width and color.
// If close is true, the shape is closed. If len(points) < 2, nothing is drawn.
// If len(points) == 2, a line is drawn.
func (i *Image) Shape(points []vector.Vector, close bool, strokeWidth float64, c color.Color) {
	if len(points) < 2 {
		return
	}
	path := polygonPath(points, close)
	i.strokePath(&path, strokeWidth, c)
}

// FillShape draws a filled polygon through the given points with the given color.
// If close is true, the shape is closed. If len(points) < 2, nothing is drawn.
// If len(points) == 2, a line is drawn.
func (i *Image) FillShape(points []vector.Vector, close bool, c color.Color) {
	if len(points) < 2 {
		return
	}
	path := polygonPath(points, close)
	i.fillPath(&path, c)
}

func (i *Image) fillPath(path *ebiten_vec.Path, c color.Color) {
	var options ebiten_vec.DrawPathOptions
	options.AntiAlias = true
	options.ColorScale.ScaleWithColor(c)
	ebiten_vec.FillPath(i.img, path, nil, &options)
}

func (i *Image) strokePath(path *ebiten_vec.Path, strokeWidth float64, c color.Color) {
	var stroked ebiten_vec.Path
	stroked.AddStroke(path, &ebiten_vec.AddStrokeOptions{
		StrokeOptions: ebiten_vec.StrokeOptions{Width: float32(strokeWidth)},
	})
	var options ebiten_vec.DrawPathOptions
	options.AntiAlias = true
	options.ColorScale.ScaleWithColor(c)
	ebiten_vec.FillPath(i.img, &stroked, nil, &options)
}
