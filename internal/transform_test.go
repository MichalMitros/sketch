package internal

import (
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/require"
)

const (
	eps = 1e-9
)

func TestNewTransformationStack_HasInitialLayer(t *testing.T) {
	s := NewTransformationStack()
	m := s.GeometryMatrix()
	require.True(t, geoMIsIdentity(m), "expected identity GeoM, got %s", m.String())
}

func TestEnsureInitialLayer_WhenEmpty(t *testing.T) {
	var s TransformationStack
	s.EnsureInitialLayer()
	require.Len(t, s.layers, 1)
}

func TestEnsureInitialLayer_Idempotent(t *testing.T) {
	s := NewTransformationStack()
	before := len(s.layers)
	s.EnsureInitialLayer()
	s.EnsureInitialLayer()
	require.Len(t, s.layers, before)
}

func TestPushPull_Single(t *testing.T) {
	s := NewTransformationStack()
	s.Push()
	require.Len(t, s.layers, 2)
	s.Pull()
	require.Len(t, s.layers, 1)
}

func TestPushPull_Multiple(t *testing.T) {
	s := NewTransformationStack()
	for i := 0; i < 5; i++ {
		s.Push()
	}
	require.Len(t, s.layers, 6)
	for i := 0; i < 5; i++ {
		s.Pull()
	}
	require.Len(t, s.layers, 1)
}

func TestPull_DoesNotRemoveInitialLayer(t *testing.T) {
	s := NewTransformationStack()
	s.Pull()
	require.Len(t, s.layers, 1)
}

func TestPull_OnEmptyStack(t *testing.T) {
	var s TransformationStack
	s.Pull()
	s.EnsureInitialLayer()
	require.Len(t, s.layers, 1)
}

func TestScale(t *testing.T) {
	s := NewTransformationStack()
	s.Scale(2, 3)
	m := s.GeometryMatrix()

	x, y := applyGeoM(m, 1, 0)
	require.InDelta(t, 2.0, x, eps)
	require.InDelta(t, 0.0, y, eps)

	x, y = applyGeoM(m, 0, 1)
	require.InDelta(t, 0.0, x, eps)
	require.InDelta(t, 3.0, y, eps)
}

func TestRotate(t *testing.T) {
	s := NewTransformationStack()
	s.Rotate(math.Pi / 2)
	m := s.GeometryMatrix()

	x, y := applyGeoM(m, 1, 0)
	require.InDelta(t, 0.0, x, eps)
	require.InDelta(t, 1.0, y, eps)
}

func TestTranslate(t *testing.T) {
	s := NewTransformationStack()
	s.Translate(10, -5)
	m := s.GeometryMatrix()

	x, y := applyGeoM(m, 0, 0)
	require.InDelta(t, 10.0, x, eps)
	require.InDelta(t, -5.0, y, eps)
}

func TestAppend_CombinesAdjacentSameKind(t *testing.T) {
	s := NewTransformationStack()
	s.Scale(2, 2)
	s.Scale(3, 4)

	require.Len(t, s.layers[0], 1)
	op := s.layers[0][0]
	require.Equal(t, TransformScale, op.kind)
	require.InDelta(t, 6.0, op.x, eps)
	require.InDelta(t, 8.0, op.y, eps)
}

func TestAppend_CombinesAdjacentRotations(t *testing.T) {
	s := NewTransformationStack()
	s.Rotate(0.5)
	s.Rotate(0.3)

	require.Len(t, s.layers[0], 1)
	op := s.layers[0][0]
	require.Equal(t, TransformRotate, op.kind)
	require.InDelta(t, 0.8, op.x, eps)
}

func TestAppend_CombinesAdjacentTranslations(t *testing.T) {
	s := NewTransformationStack()
	s.Translate(1, 2)
	s.Translate(3, 4)

	require.Len(t, s.layers[0], 1)
	op := s.layers[0][0]
	require.Equal(t, TransformTranslate, op.kind)
	require.InDelta(t, 4.0, op.x, eps)
	require.InDelta(t, 6.0, op.y, eps)
}

func TestAppend_DoesNotCombineDifferentKinds(t *testing.T) {
	s := NewTransformationStack()
	s.Scale(2, 2)
	s.Translate(1, 1)
	s.Rotate(0.1)

	require.Len(t, s.layers[0], 3)
}

func TestGeometryMatrix_ScaleThenTranslate(t *testing.T) {
	s := NewTransformationStack()
	s.Scale(2, 3)
	s.Translate(10, 20)

	x, y := applyGeoM(s.GeometryMatrix(), 1, 1)
	require.InDelta(t, 12.0, x, eps)
	require.InDelta(t, 23.0, y, eps)
}

func TestGeometryMatrix_TranslateThenRotate(t *testing.T) {
	s := NewTransformationStack()
	s.Translate(10, 0)
	s.Rotate(math.Pi / 2)

	x, y := applyGeoM(s.GeometryMatrix(), 0, 0)
	require.InDelta(t, 0.0, x, eps)
	require.InDelta(t, 10.0, y, eps)
}

func TestGeometryMatrix_RotateThenTranslate(t *testing.T) {
	s := NewTransformationStack()
	s.Rotate(math.Pi / 2)
	s.Translate(10, 0)

	x, y := applyGeoM(s.GeometryMatrix(), 0, 0)
	require.InDelta(t, 10.0, x, eps)
	require.InDelta(t, 0.0, y, eps)
}

func TestGeometryMatrix_AllThreeKinds(t *testing.T) {
	s := NewTransformationStack()
	s.Scale(2, 2)
	s.Rotate(math.Pi / 2)
	s.Translate(5, 5)

	x, y := applyGeoM(s.GeometryMatrix(), 1, 0)
	require.InDelta(t, 5.0, x, eps)
	require.InDelta(t, 7.0, y, eps)
}

func TestPushPull_Isolation(t *testing.T) {
	s := NewTransformationStack()
	s.Scale(2, 2)

	s.Push()
	s.Translate(10, 0)
	s.Pull()

	x, y := applyGeoM(s.GeometryMatrix(), 1, 0)
	require.InDelta(t, 2.0, x, eps)
	require.InDelta(t, 0.0, y, eps)
}

func TestPushPull_Nested(t *testing.T) {
	s := NewTransformationStack()

	s.Translate(1, 0)

	s.Push()
	s.Scale(2, 2)

	s.Push()
	s.Rotate(math.Pi / 2)
	s.Pull()

	x, y := applyGeoM(s.GeometryMatrix(), 1, 0)
	require.InDelta(t, 4.0, x, eps)
	require.InDelta(t, 0.0, y, eps)

	s.Pull()
	x, y = applyGeoM(s.GeometryMatrix(), 1, 0)
	require.InDelta(t, 2.0, x, eps)
	require.InDelta(t, 0.0, y, eps)
}

func TestPushPull_EmptyLayerDoesNotAffectResult(t *testing.T) {
	s := NewTransformationStack()
	s.Translate(5, 5)
	s.Push()
	s.Pull()

	x, y := applyGeoM(s.GeometryMatrix(), 0, 0)
	require.InDelta(t, 5.0, x, eps)
	require.InDelta(t, 5.0, y, eps)
}

func TestGeometryMatrix_Identity(t *testing.T) {
	s := NewTransformationStack()
	m := s.GeometryMatrix()
	require.True(t, geoMIsIdentity(m), "empty stack should produce identity GeoM, got %s", m.String())

	x, y := applyGeoM(m, 7, -3)
	require.InDelta(t, 7.0, x, eps)
	require.InDelta(t, -3.0, y, eps)
}

func TestGeometryMatrix_LargeRotationAngle(t *testing.T) {
	s := NewTransformationStack()
	s.Rotate(5 * math.Pi)

	m := s.GeometryMatrix()
	x, y := applyGeoM(m, 1, 0)
	require.InDelta(t, -1.0, x, eps)
	require.InDelta(t, 0.0, y, eps)
}

func TestAppend_OnEmptyStack(t *testing.T) {
	var s TransformationStack
	s.Scale(3, 3)
	require.Len(t, s.layers, 1)

	x, y := applyGeoM(s.GeometryMatrix(), 1, 0)
	require.InDelta(t, 3.0, x, eps)
	require.InDelta(t, 0.0, y, eps)
}

func TestGeometryMatrix_CombinedAndSeparateGiveSameResult(t *testing.T) {
	s1 := NewTransformationStack()
	s1.Scale(2, 3)
	s1.Scale(4, 5)

	s2 := NewTransformationStack()
	s2.Scale(8, 15)

	require.True(t, geoMApproxEqual(s1.GeometryMatrix(), s2.GeometryMatrix(), 1e-12))
}

func TestGeometryMatrix_CombinedRotations(t *testing.T) {
	s1 := NewTransformationStack()
	s1.Rotate(0.2)
	s1.Rotate(0.3)

	s2 := NewTransformationStack()
	s2.Rotate(0.5)

	require.True(t, geoMApproxEqual(s1.GeometryMatrix(), s2.GeometryMatrix(), 1e-12))
}

func TestGeometryMatrix_CombinedTranslations(t *testing.T) {
	s1 := NewTransformationStack()
	s1.Translate(1, 2)
	s1.Translate(3, 4)

	s2 := NewTransformationStack()
	s2.Translate(4, 6)

	require.True(t, geoMApproxEqual(s1.GeometryMatrix(), s2.GeometryMatrix(), 1e-12))
}

func geoMApproxEqual(a, b ebiten.GeoM, eps float64) bool {
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(a.Element(i, j)-b.Element(i, j)) > eps {
				return false
			}
		}
	}
	return true
}

func geoMIsIdentity(m ebiten.GeoM) bool {
	return m.Element(0, 0) == 1 &&
		m.Element(0, 1) == 0 &&
		m.Element(0, 2) == 0 &&
		m.Element(1, 0) == 0 &&
		m.Element(1, 1) == 1 &&
		m.Element(1, 2) == 0
}

func applyGeoM(m ebiten.GeoM, x, y float64) (float64, float64) {
	return m.Apply(x, y)
}
