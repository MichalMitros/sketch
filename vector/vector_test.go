package vector

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

const eps = 1e-9

func TestNew(t *testing.T) {
	v := New(3.0, 4.0)
	assert.InDelta(t, 3.0, v.X, eps)
	assert.InDelta(t, 4.0, v.Y, eps)
}

func TestZero(t *testing.T) {
	v := Zero()
	assert.InDelta(t, 0.0, v.X, eps)
	assert.InDelta(t, 0.0, v.Y, eps)
}

func TestVectorValues(t *testing.T) {
	v := New(1.0, 2.0)
	x, y := v.Values()
	assert.InDelta(t, 1.0, x, eps)
	assert.InDelta(t, 2.0, y, eps)
}

func TestVectorAdd(t *testing.T) {
	v1 := New(1.0, 2.0)
	v2 := New(3.0, 4.0)
	result := v1.Add(v2)
	assert.InDelta(t, 4.0, result.X, eps)
	assert.InDelta(t, 6.0, result.Y, eps)
}

func TestVectorSub(t *testing.T) {
	v1 := New(5.0, 6.0)
	v2 := New(1.0, 2.0)
	result := v1.Sub(v2)
	assert.InDelta(t, 4.0, result.X, eps)
	assert.InDelta(t, 4.0, result.Y, eps)
}

func TestVectorMul(t *testing.T) {
	v1 := New(2.0, 3.0)
	v2 := New(4.0, 5.0)
	result := v1.Mul(v2)
	assert.InDelta(t, 8.0, result.X, eps)
	assert.InDelta(t, 15.0, result.Y, eps)
}

func TestVectorScale(t *testing.T) {
	v := New(2.0, 3.0)
	result := v.Scale(4.0)
	assert.InDelta(t, 8.0, result.X, eps)
	assert.InDelta(t, 12.0, result.Y, eps)
}

func TestVectorMag(t *testing.T) {
	v := New(3.0, 4.0)
	assert.InDelta(t, 5.0, v.Mag(), eps)
}

func TestVectorMagZero(t *testing.T) {
	v := New(0.0, 0.0)
	assert.InDelta(t, 0.0, v.Mag(), eps)
}

func TestVectorSetMag(t *testing.T) {
	v := New(3.0, 4.0)
	result := v.SetMag(10.0)
	assert.InDelta(t, 10.0, result.Mag(), eps)
	assert.InDelta(t, v.Angle(), result.Angle(), eps)
}

func TestVectorSetMagZeroVector(t *testing.T) {
	v := New(0.0, 0.0)
	result := v.SetMag(5.0)
	assert.InDelta(t, 0.0, result.X, eps)
	assert.InDelta(t, 0.0, result.Y, eps)
}

func TestVectorDist(t *testing.T) {
	v1 := New(0.0, 0.0)
	v2 := New(3.0, 4.0)
	assert.InDelta(t, 5.0, v1.Dist(v2), eps)
}

func TestVectorDistSame(t *testing.T) {
	v := New(1.0, 2.0)
	assert.InDelta(t, 0.0, v.Dist(v), eps)
}

func TestVectorRotate(t *testing.T) {
	v := New(1.0, 0.0)
	result := v.Rotate(math.Pi / 2)
	assert.InDelta(t, 0.0, result.X, eps)
	assert.InDelta(t, 1.0, result.Y, eps)
}

func TestVectorRotateFullCircle(t *testing.T) {
	v := New(1.0, 0.0)
	result := v.Rotate(2 * math.Pi)
	assert.InDelta(t, 1.0, result.X, eps)
	assert.InDelta(t, 0.0, result.Y, eps)
}

func TestVectorAngle(t *testing.T) {
	v := New(0.0, 1.0)
	assert.InDelta(t, math.Pi/2, v.Angle(), eps)
}

func TestVectorAngleZero(t *testing.T) {
	v := New(1.0, 0.0)
	assert.InDelta(t, 0.0, v.Angle(), eps)
}

func TestVectorPolar(t *testing.T) {
	v := New(1.0, 0.0)
	p := v.Polar()
	assert.InDelta(t, 1.0, p.Radius, eps)
	assert.InDelta(t, 0.0, p.Angle, eps)
}
