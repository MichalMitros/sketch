package sketch

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/MichalMitros/sketch/vector"
	"github.com/stretchr/testify/assert"
)

func TestNewScene(t *testing.T) {
	img := NewBlankImage(vector.New(100, 200))
	s := newScene(img, color.Black, true)

	assert.Equal(t, img, s.img)
	assert.Equal(t, color.Black, s.background)
	assert.True(t, s.antyaliasing)
}

func TestSceneWidth(t *testing.T) {
	s := newTestScene(100, 200)
	assert.Equal(t, 100.0, s.Width())
}

func TestSceneHeight(t *testing.T) {
	s := newTestScene(100, 200)
	assert.Equal(t, 200.0, s.Height())
}

func TestSceneSize(t *testing.T) {
	s := newTestScene(100, 200)
	assert.Equal(t, vector.New(100, 200), s.Size())
}

func TestSceneClear(t *testing.T) {
	s := newTestScene(10, 10)
	s.Fill(color.Black)
	s.Clear()
}

func TestSceneFill(t *testing.T) {
	s := newTestScene(10, 10)
	s.Fill(color.Black)
}

func TestSceneScale(t *testing.T) {
	s := newTestScene(100, 100)
	s.Scale(2)
	s.FillRectangle(vector.New(10, 10), 10, 10, 0, color.Black)
}

func TestSceneScaleX(t *testing.T) {
	s := newTestScene(100, 100)
	s.ScaleX(2)
	s.FillRectangle(vector.New(10, 10), 10, 10, 0, color.Black)
}

func TestSceneScaleY(t *testing.T) {
	s := newTestScene(100, 100)
	s.ScaleY(2)
	s.FillRectangle(vector.New(10, 10), 10, 10, 0, color.Black)
}

func TestSceneRotate(t *testing.T) {
	s := newTestScene(100, 100)
	s.Rotate(math.Pi / 2)
	s.FillRectangle(vector.New(10, 10), 10, 10, 0, color.Black)
}

func TestSceneTransform(t *testing.T) {
	s := newTestScene(100, 100)
	s.Transform(vector.New(20, 30))
	s.FillRectangle(vector.New(0, 0), 10, 10, 0, color.Black)
}

func TestSceneTransformX(t *testing.T) {
	s := newTestScene(100, 100)
	s.TransformX(20)
	s.FillRectangle(vector.New(0, 0), 10, 10, 0, color.Black)
}

func TestSceneTransformY(t *testing.T) {
	s := newTestScene(100, 100)
	s.TransformY(30)
	s.FillRectangle(vector.New(0, 0), 10, 10, 0, color.Black)
}

func TestScenePush(t *testing.T) {
	s := newTestScene(100, 100)
	s.FillRectangle(vector.New(0, 0), 10, 10, 0, color.Black)

	s.Push()
	s.Transform(vector.New(30, 30))
	s.FillRectangle(vector.New(0, 0), 10, 10, 0, color.Black)
}

func TestScenePull(t *testing.T) {
	s := newTestScene(100, 100)

	s.Push()
	s.Transform(vector.New(30, 30))
	s.FillRectangle(vector.New(0, 0), 10, 10, 0, color.Black)

	s.Pull()
	s.FillRectangle(vector.New(0, 0), 10, 10, 0, color.Black)
}

func TestSceneLine(t *testing.T) {
	s := newTestScene(100, 100)
	s.Line(vector.New(10, 50), vector.New(90, 50), 2, color.Black)
}

func TestSceneRectangle(t *testing.T) {
	s := newTestScene(100, 100)
	s.Rectangle(vector.New(20, 20), 40, 30, 2, color.Black)
}

func TestSceneFillRectangle(t *testing.T) {
	s := newTestScene(100, 100)
	s.FillRectangle(vector.New(20, 20), 40, 30, 0, color.Black)
}

func TestSceneCircle(t *testing.T) {
	s := newTestScene(100, 100)
	s.Circle(vector.New(50, 50), 20, 2, color.Black)
}

func TestSceneFillCircle(t *testing.T) {
	s := newTestScene(100, 100)
	s.FillCircle(vector.New(50, 50), 20, 0, color.Black)
}

func TestSceneArc(t *testing.T) {
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
			s := newTestScene(100, 100)
			s.Arc(vector.New(50, 50), 20, tc.startAngle, tc.endAngle, 2, color.Black)
		})
	}
}

func TestSceneFillArc(t *testing.T) {
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
			s := newTestScene(100, 100)
			s.FillArc(vector.New(50, 50), 20, tc.startAngle, tc.endAngle, color.Black)
		})
	}
}

func TestSceneShape(t *testing.T) {
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
			s := newTestScene(100, 100)
			s.Shape(tc.points, tc.close, tc.strokeWidth, color.Black)
		})
	}
}

func TestSceneFillShape(t *testing.T) {
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
			s := newTestScene(100, 100)
			s.FillShape(tc.points, tc.close, color.Black)
		})
	}
}

func TestSceneDrawImage(t *testing.T) {
	scenarios := map[string]struct {
		setupImg func() *Image
		pos      vector.Vector
		size     vector.Vector
	}{
		"nil image": {
			setupImg: func() *Image { return nil },
			pos:      vector.New(10, 10),
			size:     vector.New(20, 20),
		},
		"ok": {
			setupImg: func() *Image { return NewFilledImage(vector.New(10, 10), color.Black) },
			pos:      vector.New(20, 30),
			size:     vector.New(40, 50),
		},
		"ok with transform": {
			setupImg: func() *Image { return NewFilledImage(vector.New(10, 10), color.Black) },
			pos:      vector.New(0, 0),
			size:     vector.New(20, 20),
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			s := newTestScene(100, 100)

			if name == "with transform" {
				s.Transform(vector.New(30, 30))
			}

			img := tc.setupImg()
			s.DrawImage(img, tc.pos, tc.size)
		})
	}
}

func TestCirclePath(t *testing.T) {
	path := circlePath(vector.New(50, 50), 20)
	bounds := path.Bounds()

	assert.InDelta(t, 30.0, float64(bounds.Min.X), 1.0)
	assert.InDelta(t, 30.0, float64(bounds.Min.Y), 1.0)
	assert.InDelta(t, 70.0, float64(bounds.Max.X), 1.0)
	assert.InDelta(t, 70.0, float64(bounds.Max.Y), 1.0)
}

func TestArcPath(t *testing.T) {
	scenarios := map[string]struct {
		startAngle float64
		endAngle   float64
		minX       float64
		minY       float64
		maxX       float64
		maxY       float64
	}{
		"half circle top": {
			startAngle: 0,
			endAngle:   math.Pi,
			minX:       30,
			minY:       50,
			maxX:       70,
			maxY:       70,
		},
		"quarter circle": {
			startAngle: 0,
			endAngle:   math.Pi / 2,
			minX:       50,
			minY:       50,
			maxX:       70,
			maxY:       70,
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			path := arcPath(vector.New(50, 50), 20, tc.startAngle, tc.endAngle)
			bounds := path.Bounds()

			assert.InDelta(t, tc.minX, float64(bounds.Min.X), 1.0)
			assert.InDelta(t, tc.minY, float64(bounds.Min.Y), 1.0)
			assert.InDelta(t, tc.maxX, float64(bounds.Max.X), 1.0)
			assert.InDelta(t, tc.maxY, float64(bounds.Max.Y), 1.0)
		})
	}
}

func TestFillArcPath(t *testing.T) {
	scenarios := map[string]struct {
		startAngle float64
		endAngle   float64
		contains   vector.Vector
	}{
		"half pie top": {
			startAngle: 0,
			endAngle:   math.Pi,
			contains:   vector.New(50, 50),
		},
		"quarter pie": {
			startAngle: 0,
			endAngle:   math.Pi / 2,
			contains:   vector.New(50, 50),
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			path := fillArcPath(vector.New(50, 50), 20, tc.startAngle, tc.endAngle)
			bounds := path.Bounds()

			assert.True(t, bounds.Min.X <= int(tc.contains.X) && int(tc.contains.X) <= bounds.Max.X)
			assert.True(t, bounds.Min.Y <= int(tc.contains.Y) && int(tc.contains.Y) <= bounds.Max.Y)
		})
	}
}

func TestRectanglePath(t *testing.T) {
	path := rectanglePath(vector.New(20, 20), 40, 30)
	bounds := path.Bounds()

	assert.Equal(t, 20, bounds.Min.X)
	assert.Equal(t, 20, bounds.Min.Y)
	assert.Equal(t, 60, bounds.Max.X)
	assert.Equal(t, 50, bounds.Max.Y)
}

func TestPolygonPath(t *testing.T) {
	scenarios := map[string]struct {
		points       []vector.Vector
		close        bool
		expectBounds image.Rectangle
	}{
		"open": {
			points:       []vector.Vector{vector.New(10, 10), vector.New(50, 10), vector.New(50, 50)},
			close:        false,
			expectBounds: image.Rect(10, 10, 50, 50),
		},
		"closed": {
			points:       []vector.Vector{vector.New(10, 10), vector.New(50, 10), vector.New(50, 50)},
			close:        true,
			expectBounds: image.Rect(10, 10, 50, 50),
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			path := polygonPath(tc.points, tc.close)
			bounds := path.Bounds()

			assert.Equal(t, tc.expectBounds.Min.X, bounds.Min.X)
			assert.Equal(t, tc.expectBounds.Min.Y, bounds.Min.Y)
			assert.Equal(t, tc.expectBounds.Max.X, bounds.Max.X)
			assert.Equal(t, tc.expectBounds.Max.Y, bounds.Max.Y)
		})
	}
}

func TestTransformedPath(t *testing.T) {
	s := newTestScene(100, 100)
	s.Transform(vector.New(10, 20))

	path := rectanglePath(vector.New(0, 0), 30, 40)
	transformed := s.transformedPath(&path)
	bounds := transformed.Bounds()

	assert.Equal(t, 10, bounds.Min.X)
	assert.Equal(t, 20, bounds.Min.Y)
	assert.Equal(t, 40, bounds.Max.X)
	assert.Equal(t, 60, bounds.Max.Y)
}

func TestTransformedPathWithScale(t *testing.T) {
	s := newTestScene(100, 100)
	s.Scale(2)

	path := rectanglePath(vector.New(10, 10), 20, 30)
	transformed := s.transformedPath(&path)
	bounds := transformed.Bounds()

	assert.Equal(t, 20, bounds.Min.X)
	assert.Equal(t, 20, bounds.Min.Y)
	assert.Equal(t, 60, bounds.Max.X)
	assert.Equal(t, 80, bounds.Max.Y)
}

func TestDrawPathOptions(t *testing.T) {
	scenarios := map[string]struct {
		antyaliasing bool
		c            color.Color
	}{
		"with antyaliasing": {
			antyaliasing: true,
			c:            color.Black,
		},
		"without antyaliasing": {
			antyaliasing: false,
			c:            color.White,
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			img := NewBlankImage(vector.New(10, 10))
			s := newScene(img, color.White, tc.antyaliasing)
			opts := s.drawPathOptions(tc.c)

			assert.Equal(t, tc.antyaliasing, opts.AntiAlias)
		})
	}
}

func TestFillPath(t *testing.T) {
	s := newTestScene(100, 100)
	path := rectanglePath(vector.New(10, 10), 20, 30)
	s.fillPath(&path, color.Black)
}

func TestStrokePath(t *testing.T) {
	s := newTestScene(100, 100)
	path := rectanglePath(vector.New(10, 10), 20, 30)
	s.strokePath(&path, 2, color.Black)
}

func TestSceneTransformStackIsolation(t *testing.T) {
	s := newTestScene(100, 100)

	s.Scale(2)
	s.Push()
	s.Transform(vector.New(10, 10))
	s.Pull()

	path := rectanglePath(vector.New(5, 5), 10, 10)
	transformed := s.transformedPath(&path)
	bounds := transformed.Bounds()

	assert.Equal(t, 10, bounds.Min.X)
	assert.Equal(t, 10, bounds.Min.Y)
}

func newTestScene(w, h int) *Scene {
	img := NewBlankImage(vector.New(float64(w), float64(h)))
	img.Fill(color.White)
	return newScene(img, color.White, true)
}
