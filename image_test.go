package sketch

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"testing"
	"testing/fstest"

	"github.com/MichalMitros/sketch/vector"
	"github.com/stretchr/testify/assert"
)

func generateTestPNG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func TestImageFromFS(t *testing.T) {
	pngData := generateTestPNG(10, 20)
	fsys := fstest.MapFS{
		"test.png": &fstest.MapFile{Data: pngData},
	}
	img, err := ImageFromFS(fsys, "test.png")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.Equal(t, 10.0, img.Width())
	assert.Equal(t, 20.0, img.Height())
}

func TestImageFromStdImage(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 30, 40))
	img := ImageFromStdImage(src)
	assert.NotNil(t, img)
	assert.Equal(t, 30.0, img.Width())
	assert.Equal(t, 40.0, img.Height())
}

func TestNewBlankImage(t *testing.T) {
	img := NewBlankImage(vector.New(50, 60))
	assert.Equal(t, 50.0, img.Width())
	assert.Equal(t, 60.0, img.Height())
}

func TestNewFilledImage(t *testing.T) {
	scenarios := map[string]struct {
		dim   vector.Vector
		color color.Color
	}{
		"red":   {dim: vector.New(10, 20), color: color.RGBA{255, 0, 0, 255}},
		"green": {dim: vector.New(30, 40), color: color.RGBA{0, 255, 0, 255}},
		"blue":  {dim: vector.New(5, 5), color: color.RGBA{0, 0, 255, 255}},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			img := NewFilledImage(tc.dim, tc.color)
			assert.Equal(t, tc.dim.X, img.Width())
			assert.Equal(t, tc.dim.Y, img.Height())
		})
	}
}

func TestNewWhiteImage(t *testing.T) {
	img := NewWhiteImage(vector.New(10, 10))
	assert.Equal(t, 10.0, img.Width())
	assert.Equal(t, 10.0, img.Height())
}

func TestNewBlackImage(t *testing.T) {
	img := NewBlackImage(vector.New(10, 10))
	assert.Equal(t, 10.0, img.Width())
	assert.Equal(t, 10.0, img.Height())
}

func TestImageWidth(t *testing.T) {
	scenarios := map[string]struct {
		dim    vector.Vector
		expect float64
	}{
		"square":       {dim: vector.New(100, 100), expect: 100},
		"wide":         {dim: vector.New(200, 50), expect: 200},
		"tall":         {dim: vector.New(30, 120), expect: 30},
		"single pixel": {dim: vector.New(1, 1), expect: 1},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			img := NewBlankImage(tc.dim)
			assert.Equal(t, tc.expect, img.Width())
		})
	}
}

func TestImageHeight(t *testing.T) {
	scenarios := map[string]struct {
		dim    vector.Vector
		expect float64
	}{
		"square":       {dim: vector.New(100, 100), expect: 100},
		"wide":         {dim: vector.New(200, 50), expect: 50},
		"tall":         {dim: vector.New(30, 120), expect: 120},
		"single pixel": {dim: vector.New(1, 1), expect: 1},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			img := NewBlankImage(tc.dim)
			assert.Equal(t, tc.expect, img.Height())
		})
	}
}

func TestImageSize(t *testing.T) {
	scenarios := map[string]struct {
		dim    vector.Vector
		expect vector.Vector
	}{
		"square":    {dim: vector.New(100, 100), expect: vector.New(100, 100)},
		"rectangle": {dim: vector.New(200, 50), expect: vector.New(200, 50)},
		"one pixel": {dim: vector.New(1, 1), expect: vector.New(1, 1)},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			img := NewBlankImage(tc.dim)
			assert.Equal(t, tc.expect, img.Size())
		})
	}
}

func TestImageClone(t *testing.T) {
	original := NewFilledImage(vector.New(10, 10), color.RGBA{255, 0, 0, 255})
	clone := original.Clone()

	assert.Equal(t, original.Width(), clone.Width())
	assert.Equal(t, original.Height(), clone.Height())

	clone.Set(vector.New(0, 0), color.RGBA{0, 255, 0, 255})
}

func TestDefaultDrawOptions(t *testing.T) {
	opts := DefaultDrawOptions()
	assert.Equal(t, vector.New(0, 0), opts.Pos)
	assert.Equal(t, vector.New(0, 0), opts.Anchor)
	assert.Equal(t, vector.New(1, 1), opts.Scale)
	assert.Equal(t, 0.0, opts.Rotation)
	assert.Nil(t, opts.Tint)
	assert.Equal(t, 1.0, opts.Opacity)
}

func TestImageDraw(t *testing.T) {
	scenarios := map[string]struct {
		opts *DrawOptions
	}{
		"nil opts": {
			opts: nil,
		},
		"zero scale X": {
			opts: &DrawOptions{
				Scale:   vector.New(0, 1),
				Opacity: 1,
			},
		},
		"zero scale Y": {
			opts: &DrawOptions{
				Scale:   vector.New(1, 0),
				Opacity: 1,
			},
		},
		"zero opacity": {
			opts: &DrawOptions{
				Scale:   vector.New(1, 1),
				Opacity: 0,
			},
		},
		"position offset": {
			opts: &DrawOptions{
				Scale:   vector.New(1, 1),
				Pos:     vector.New(5, 5),
				Opacity: 1,
			},
		},
		"anchor center": {
			opts: &DrawOptions{
				Scale:   vector.New(1, 1),
				Pos:     vector.New(10, 10),
				Anchor:  vector.New(0.5, 0.5),
				Opacity: 1,
			},
		},
		"with tint": {
			opts: &DrawOptions{
				Scale:   vector.New(1, 1),
				Opacity: 1,
				Tint:    color.RGBA{0, 0, 255, 255},
			},
		},
		"with rotation": {
			opts: &DrawOptions{
				Scale:    vector.New(1, 1),
				Opacity:  1,
				Rotation: math.Pi,
			},
		},
		"with scale": {
			opts: &DrawOptions{
				Scale:   vector.New(2, 2),
				Opacity: 1,
			},
		},
		"negative scale": {
			opts: &DrawOptions{
				Scale:   vector.New(-1, -1),
				Opacity: 1,
			},
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			src := NewFilledImage(vector.New(10, 10), color.RGBA{255, 0, 0, 255})
			dst := NewFilledImage(vector.New(50, 50), color.White)
			src.Draw(dst, tc.opts)
		})
	}
}

func TestImageDrawAt(t *testing.T) {
	src := NewFilledImage(vector.New(10, 10), color.RGBA{255, 0, 0, 255})
	dst := NewFilledImage(vector.New(50, 50), color.White)
	src.DrawAt(dst, vector.New(20, 30))
}

func TestImageAt(t *testing.T) {
	img := NewBlankImage(vector.New(10, 10))
	img.Set(vector.New(0, 0), color.RGBA{255, 0, 0, 255})
	img.Set(vector.New(9, 9), color.RGBA{0, 255, 0, 255})
}

func TestImageSet(t *testing.T) {
	img := NewBlankImage(vector.New(10, 10))
	img.Set(vector.New(5, 5), color.RGBA{255, 0, 0, 255})
}

func TestImageFill(t *testing.T) {
	scenarios := map[string]struct {
		color color.Color
	}{
		"red":   {color: color.RGBA{255, 0, 0, 255}},
		"green": {color: color.RGBA{0, 255, 0, 255}},
		"blue":  {color: color.RGBA{0, 0, 255, 255}},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			img := NewBlankImage(vector.New(10, 10))
			img.Fill(tc.color)
		})
	}
}

func TestImageClear(t *testing.T) {
	img := NewFilledImage(vector.New(10, 10), color.RGBA{255, 0, 0, 255})
	img.Clear()
}

func TestImageSubImage(t *testing.T) {
	scenarios := map[string]struct {
		pos     vector.Vector
		size    vector.Vector
		expectW float64
		expectH float64
	}{
		"top left quarter": {
			pos:     vector.New(0, 0),
			size:    vector.New(5, 5),
			expectW: 5,
			expectH: 5,
		},
		"bottom right quarter": {
			pos:     vector.New(5, 5),
			size:    vector.New(5, 5),
			expectW: 5,
			expectH: 5,
		},
		"center region": {
			pos:     vector.New(2, 2),
			size:    vector.New(6, 6),
			expectW: 6,
			expectH: 6,
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			original := NewFilledImage(vector.New(10, 10), color.White)
			sub := original.SubImage(tc.pos, tc.size)

			assert.Equal(t, tc.expectW, sub.Width())
			assert.Equal(t, tc.expectH, sub.Height())
		})
	}
}

func TestImageCopyTo(t *testing.T) {
	src := NewFilledImage(vector.New(10, 10), color.RGBA{255, 0, 0, 255})
	dst := NewFilledImage(vector.New(50, 50), color.White)
	src.CopyTo(dst, vector.New(5, 5), vector.New(1, 1), vector.New(20, 30))
}

func TestImageTransformations(t *testing.T) {
	src := NewFilledImage(vector.New(10, 10), color.RGBA{255, 0, 0, 255})
	dst := NewFilledImage(vector.New(100, 100), color.White)

	opts := DrawOptions{
		Pos:      vector.New(50, 50),
		Anchor:   vector.New(0.5, 0.5),
		Scale:    vector.New(2, 3),
		Rotation: math.Pi / 4,
		Opacity:  0.5,
		Tint:     color.RGBA{255, 255, 255, 255},
	}
	src.Draw(dst, &opts)

	opts2 := DrawOptions{
		Pos:      vector.New(75, 75),
		Anchor:   vector.New(1, 1),
		Scale:    vector.New(0.5, 1.5),
		Rotation: -math.Pi / 6,
		Opacity:  0.8,
	}
	src.Draw(dst, &opts2)
}

func TestImageSubImageSharedData(t *testing.T) {
	original := NewFilledImage(vector.New(10, 10), color.White)
	sub := original.SubImage(vector.New(2, 2), vector.New(4, 4))
	sub.Set(vector.New(0, 0), color.RGBA{255, 0, 0, 255})
}

func TestImageLine(t *testing.T) {
	img := NewBlankImage(vector.New(100, 100))
	img.Line(vector.New(10, 50), vector.New(90, 50), 2, color.Black)
}

func TestImageRectangle(t *testing.T) {
	img := NewBlankImage(vector.New(100, 100))
	img.Rectangle(vector.New(20, 20), 40, 30, 2, color.Black)
}

func TestImageFillRectangle(t *testing.T) {
	img := NewBlankImage(vector.New(100, 100))
	img.FillRectangle(vector.New(10, 10), 20, 30, 0, color.Black)
}

func TestImageCircle(t *testing.T) {
	img := NewBlankImage(vector.New(100, 100))
	img.Circle(vector.New(50, 50), 20, 2, color.Black)
}

func TestImageFillCircle(t *testing.T) {
	img := NewBlankImage(vector.New(100, 100))
	img.FillCircle(vector.New(50, 50), 20, 0, color.Black)
}

func TestImageArc(t *testing.T) {
	scenarios := map[string]struct {
		startAngle float64
		endAngle   float64
	}{
		"full circle":     {startAngle: 0, endAngle: 2 * math.Pi},
		"half circle top": {startAngle: 0, endAngle: math.Pi},
		"quarter circle":  {startAngle: 0, endAngle: math.Pi / 2},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			img := NewBlankImage(vector.New(100, 100))
			img.Arc(vector.New(50, 50), 20, tc.startAngle, tc.endAngle, 2, color.Black)
		})
	}
}

func TestImageFillArc(t *testing.T) {
	scenarios := map[string]struct {
		startAngle float64
		endAngle   float64
	}{
		"full pie":     {startAngle: 0, endAngle: 2 * math.Pi},
		"half pie top": {startAngle: 0, endAngle: math.Pi},
		"quarter pie":  {startAngle: 0, endAngle: math.Pi / 2},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			img := NewBlankImage(vector.New(100, 100))
			img.FillArc(vector.New(50, 50), 20, tc.startAngle, tc.endAngle, color.Black)
		})
	}
}

func TestImageShape(t *testing.T) {
	scenarios := map[string]struct {
		points      []vector.Vector
		close       bool
		strokeWidth float64
	}{
		"less than 2 points": {
			points:      []vector.Vector{vector.New(10, 10)},
			close:       false,
			strokeWidth: 2,
		},
		"two points": {
			points:      []vector.Vector{vector.New(10, 50), vector.New(90, 50)},
			close:       false,
			strokeWidth: 2,
		},
		"open polygon": {
			points:      []vector.Vector{vector.New(10, 10), vector.New(50, 10), vector.New(50, 50)},
			close:       false,
			strokeWidth: 2,
		},
		"closed polygon": {
			points:      []vector.Vector{vector.New(10, 10), vector.New(50, 10), vector.New(50, 50)},
			close:       true,
			strokeWidth: 2,
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			img := NewBlankImage(vector.New(100, 100))
			img.Shape(tc.points, tc.close, tc.strokeWidth, color.Black)
		})
	}
}

func TestImageFillShape(t *testing.T) {
	scenarios := map[string]struct {
		points []vector.Vector
		close  bool
	}{
		"less than 2 points": {
			points: []vector.Vector{vector.New(10, 10)},
			close:  false,
		},
		"two points": {
			points: []vector.Vector{vector.New(10, 50), vector.New(90, 50)},
			close:  false,
		},
		"open polygon": {
			points: []vector.Vector{vector.New(10, 10), vector.New(50, 10), vector.New(50, 50)},
			close:  false,
		},
		"closed polygon": {
			points: []vector.Vector{vector.New(10, 10), vector.New(50, 10), vector.New(50, 50)},
			close:  true,
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			img := NewBlankImage(vector.New(100, 100))
			img.FillShape(tc.points, tc.close, color.Black)
		})
	}
}
