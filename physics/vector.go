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

// Diagonal upwards to the right. (1/sqrt(2), 1/sqrt(2))
func Vec2UpRight() Vec2 {
	return Vec2{1 / math.Sqrt2, 1 / math.Sqrt2}
}

// Diagonal upwards to the left. (-1/sqrt(2), 1/sqrt(2))
func Vec2UpLeft() Vec2 {
	return Vec2{-1 / math.Sqrt2, 1 / math.Sqrt2}
}

// Diagonal downwards to the right. (1/sqrt(2), -1/sqrt(2))
func Vec2DownRight() Vec2 {
	return Vec2{1 / math.Sqrt2, -1 / math.Sqrt2}
}

// Diagonal downwards to the left. (-1/sqrt(2), -1/sqrt(2))
func Vec2DownLeft() Vec2 {
	return Vec2{-1 / math.Sqrt2, -1 / math.Sqrt2}
}

// Checks if two vectors are the same within a given tolerance
func (u Vec2) EqualsWithTol(other Vec2, tolerance float64) bool {
	return math.Abs(u.X-other.X) < tolerance && math.Abs(u.Y-other.Y) < tolerance
}

// Checks if two vectors are the same within a small floating point error tolerance.
func (u Vec2) ApproxEquals(other Vec2) bool {
	return u.EqualsWithTol(other, tolerance)
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
	if u.Length() == 0 {
		return Vec2Zero()
	}
	return u.Mul(1 / u.Length())
}

// Euclidean distance to another vector
func (u Vec2) DistanceTo(other Vec2) float64 {
	return math.Sqrt(math.Pow(u.X-other.X, 2) + math.Pow(u.Y-other.Y, 2))
}

func AddVectors(vectors ...Vec2) Vec2 {
	var result Vec2
	for _, u := range vectors {
		result = result.Add(u)
	}
	return result
}

func SubtractVectors(vectors ...Vec2) Vec2 {
	result := vectors[0]
	for _, u := range vectors[1:] {
		result = result.Sub(u)
	}
	return result
}
