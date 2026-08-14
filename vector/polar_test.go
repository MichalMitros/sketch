package vector

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPolar(t *testing.T) {
	p := NewPolar(5.0, math.Pi)
	assert.InDelta(t, 5.0, p.Radius, eps)
	assert.InDelta(t, math.Pi, p.Angle, eps)
}

func TestPolarValues(t *testing.T) {
	p := NewPolar(3.0, math.Pi/2)
	r, a := p.Values()
	assert.InDelta(t, 3.0, r, eps)
	assert.InDelta(t, math.Pi/2, a, eps)
}

func TestPolarRotate(t *testing.T) {
	p := NewPolar(1.0, 0.0)
	result := p.Rotate(math.Pi)
	assert.InDelta(t, 1.0, result.Radius, eps)
	assert.InDelta(t, math.Pi, result.Angle, eps)
}

func TestPolarAdd(t *testing.T) {
	p1 := NewPolar(1.0, 0.0)
	p2 := NewPolar(1.0, math.Pi)
	result := p1.Add(p2)
	assert.InDelta(t, 0.0, result.Radius, eps)
}

func TestPolarDist(t *testing.T) {
	p1 := NewPolar(1.0, 0.0)
	p2 := NewPolar(2.0, 0.0)
	assert.InDelta(t, 1.0, p1.Dist(p2), eps)
}

func TestPolarDistSame(t *testing.T) {
	p := NewPolar(3.0, math.Pi/4)
	assert.InDelta(t, 0.0, p.Dist(p), eps)
}

func TestPolarVector(t *testing.T) {
	p := NewPolar(2.0, math.Pi)
	v := p.Vector()
	assert.InDelta(t, -2.0, v.X, eps)
	assert.InDelta(t, 0.0, v.Y, eps)
}

func TestPolarVectorRoundTrip(t *testing.T) {
	v := New(3.0, 4.0)
	p := v.Polar()
	v2 := p.Vector()
	assert.InDelta(t, v.X, v2.X, eps)
	assert.InDelta(t, v.Y, v2.Y, eps)
}
