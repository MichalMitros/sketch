package vector

import "math"

// Polar represents a 2D polar vector with Radius and Angle components stored as float64.
type Polar struct {
	Radius, Angle float64
}

// NewPolar returns a new polar vector with the given radius and angle components.
func NewPolar(radius, angle float64) Polar {
	return Polar{Radius: radius, Angle: angle}
}

// Values returns polar vector's radius and angle components.
func (p Polar) Values() (float64, float64) {
	return p.Radius, p.Angle
}

// Rotate returns a new polar vector rotated by angle radians (counter-clockwise).
func (p Polar) Rotate(angle float64) Polar {
	p.Angle += angle
	return p
}

// Add returns a new polar vector equal to the sum of p and other.
func (p Polar) Add(other Polar) Polar {
	x1, y1 := p.Vector().Values()
	x2, y2 := other.Vector().Values()
	return Polar{
		Radius: math.Sqrt((x1+x2)*(x1+x2) + (y1+y2)*(y1+y2)),
		Angle:  math.Atan2((y1 + y2), (x1 + x2)),
	}
}

// Dist returns distance between vector and provided other vector.
func (p Polar) Dist(other Polar) float64 {
	x1, y1 := p.Vector().Values()
	x2, y2 := other.Vector().Values()
	return math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
}

// Vector returns a new vector with the same radius and angle as the polar vector.
func (p Polar) Vector() Vector {
	return Vector{
		X: p.Radius * math.Cos(p.Angle),
		Y: p.Radius * math.Sin(p.Angle),
	}
}
