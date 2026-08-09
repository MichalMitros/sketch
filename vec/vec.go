// Package vec provides a simple 2D vector type.
package vec

import "math"

// Vector represents a 2D vector with X and Y components stored as float64.
type Vector struct {
	X, Y float64
}

func New(x, y float64) Vector {
	return Vector{X: x, Y: y}
}

// Add returns a new vector equal to the sum of v and other.
func (v Vector) Add(other Vector) Vector {
	return Vector{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub returns a new vector equal to the difference of v and other.
func (v Vector) Sub(other Vector) Vector {
	return Vector{X: v.X - other.X, Y: v.Y - other.Y}
}

// Mul returns a new vector equal to the product of v and other.
func (v Vector) Mul(other Vector) Vector {
	return Vector{X: v.X * other.X, Y: v.Y * other.Y}
}

// Scale returns a new vector with both components multiplied by s.
func (v Vector) Scale(s float64) Vector {
	return Vector{X: v.X * s, Y: v.Y * s}
}

// Mag returns the magnitude (length) of the vector.
func (v Vector) Mag() float64 {
	return float64(math.Sqrt(v.X*v.X + v.Y*v.Y))
}

// SetMag sets the magnitude (length) of the vector without changing its direction.
func (v Vector) SetMag(mag float64) Vector {
	current := v.Mag()
	if current == 0 {
		return v
	}
	scale := mag / current
	return New(v.X*scale, v.Y*scale)
}

// Dist returns distance between vector and provided other vector.
func (v Vector) Dist(other Vector) float64 {
	return float64(math.Sqrt((other.X-v.X)*(other.X-v.X) + (other.Y-v.Y)*(other.Y-v.Y)))
}

// Rotate returns a new vector rotated by angle radians (counter-clockwise).
func (v Vector) Rotate(angle float64) Vector {
	cos := float64(math.Cos(angle))
	sin := float64(math.Sin(angle))
	return Vector{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

// Angle returns the angle of the vector in radians (counter-clockwise).
func (v Vector) Angle() float64 {
	return float64(math.Atan2(v.Y, v.X))
}
