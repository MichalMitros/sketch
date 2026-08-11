package internal

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// TransformOperationKind is the kind of a transformation operation.
type TransformOperationKind uint8

const (
	// TransformScale scales the coordinate axes.
	TransformScale TransformOperationKind = iota
	// TransformRotate rotates the coordinate system.
	TransformRotate
	// TransformTranslate moves the origin of the coordinate system.
	TransformTranslate
)

// TransformOperation is a single transformation operation.
type TransformOperation struct {
	kind TransformOperationKind
	x    float64
	y    float64
}

// TransformationStack stores operations rather than transformed coordinates.
// Replaying them only when something is drawn avoids compounding rounding errors
// between drawing calls.
type TransformationStack struct {
	layers [][]TransformOperation
}

// NewTransformationStack creates a new TransformationStack.
func NewTransformationStack() TransformationStack {
	return TransformationStack{layers: [][]TransformOperation{nil}}
}

// EnsureInitialLayer ensures that the initial layer exists.
func (s *TransformationStack) EnsureInitialLayer() {
	if len(s.layers) == 0 {
		s.layers = append(s.layers, nil)
	}
}

// Push adds transformation layer to the stack.
func (s *TransformationStack) Push() {
	s.EnsureInitialLayer()
	s.layers = append(s.layers, nil)
}

// Pull removes the most recently pushed transformation layer.
func (s *TransformationStack) Pull() {
	s.EnsureInitialLayer()
	if len(s.layers) > 1 {
		s.layers = s.layers[:len(s.layers)-1]
	}
}

// Append appends a transformation operation to the most recent layer.
func (s *TransformationStack) Append(op TransformOperation) {
	s.EnsureInitialLayer()
	layer := len(s.layers) - 1
	operations := s.layers[layer]

	// Combining adjacent operations reduces the number of floating-point and
	// trigonometric calculations without changing operation order.
	if len(operations) > 0 && operations[len(operations)-1].kind == op.kind {
		last := &operations[len(operations)-1]
		switch op.kind {
		case TransformScale:
			last.x *= op.x
			last.y *= op.y
		case TransformRotate:
			last.x += op.x
		case TransformTranslate:
			last.x += op.x
			last.y += op.y
		}
		s.layers[layer] = operations
		return
	}

	s.layers[layer] = append(operations, op)
}

// Scale scales the coordinate axes.
func (s *TransformationStack) Scale(x, y float64) {
	s.Append(TransformOperation{kind: TransformScale, x: x, y: y})
}

// Rotate rotates the coordinate system.
func (s *TransformationStack) Rotate(angle float64) {
	s.Append(TransformOperation{kind: TransformRotate, x: angle})
}

// Translate moves the origin of the coordinate system.
func (s *TransformationStack) Translate(x, y float64) {
	s.Append(TransformOperation{kind: TransformTranslate, x: x, y: y})
}

// GeometryMatrix returns the transformation matrix.
func (s *TransformationStack) GeometryMatrix() ebiten.GeoM {
	s.EnsureInitialLayer()
	var matrix ebiten.GeoM
	for _, layer := range s.layers {
		for _, op := range layer {
			switch op.kind {
			case TransformScale:
				matrix.Scale(op.x, op.y)
			case TransformRotate:
				// Keeping the accumulated angle small improves argument
				// reduction in trigonometric functions.
				matrix.Rotate(math.Remainder(op.x, 2*math.Pi))
			case TransformTranslate:
				matrix.Translate(op.x, op.y)
			}
		}
	}
	return matrix
}
