package physics

import "math"

const tolerance = 1e-9

func FloatEquals(first, second float64) bool {
	return math.Abs(first-second) < tolerance
}

type Vec2 struct {
	X, Y float64
}

// (0, 0)
func Vec2Zero() Vec2 {
	return Vec2{0, 0}
}

// (1, 0)
func Vec2Right() Vec2 {
	return Vec2{1, 0}
}

// (-1, 0)
func Vec2Left() Vec2 {
	return Vec2{-1, 0}
}

// (0, 1)
func Vec2Up() Vec2 {
	return Vec2{0, 1}
}

// (0, -1)
func Vec2Down() Vec2 {
	return Vec2{0, -1}
}

// Checks if two vectors are the same within a small floating point error tolerance.
func (u Vec2) EqualsWithTol(other Vec2) bool {
	return math.Abs(u.X-other.X) < tolerance && math.Abs(u.Y-other.Y) < tolerance
}

func (u Vec2) Add(other Vec2) Vec2 {
	u.X += other.X
	u.Y += other.Y
	return u
}

func (u Vec2) Sub(other Vec2) Vec2 {
	u.X -= other.X
	u.Y -= other.Y
	return u
}

func (u Vec2) Dot(other Vec2) float64 {
	return u.X*other.X + u.Y*other.Y
}

// Scalar multiplication: 5v = {5v.X, 5v.Y}
func (u Vec2) Mul(scalar float64) Vec2 {
	return Vec2{scalar * u.X, scalar * u.Y}
}

// Length or magnitude of vector u is sqrt(u.X^2 + u.Y^2)
func (u Vec2) Length() float64 {
	return math.Sqrt(u.X*u.X + u.Y*u.Y)
}

func (u Vec2) Normalize() Vec2 {
	// Check division by zero
	if math.Abs(u.X) < tolerance && math.Abs(u.Y) < tolerance {
		return Vec2Zero()
	}
	if math.Abs(u.X) < tolerance {
		return Vec2{0, math.Abs(u.Y) / u.Y}
	}
	if math.Abs(u.Y) < tolerance {
		return Vec2{math.Abs(u.X) / u.X, 0}
	}
	return Vec2{math.Abs(u.X) / u.X, math.Abs(u.Y) / u.Y}
}
