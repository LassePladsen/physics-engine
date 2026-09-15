package physics

import "math"

const tolerance = 1e-9

func FloatEquals(first, second float64) bool {
	return math.Abs(first-second) < tolerance
}

type Vec2 struct {
	X, Y float64
}

// Zero vector (0, 0)
func Zero() Vec2 {
	return Vec2{0, 0}
}

// (1, 0)
func UnitRight() Vec2 {
	return Vec2{1, 0}
}

// (-1, 0)
func UnitLeft() Vec2 {
	return Vec2{-1, 0}
}

// (0, 1)
func UnitUp() Vec2 {
	return Vec2{0, 1}
}

// (0, -1)
func UnitDown() Vec2 {
	return Vec2{0, -1}
}

// Diagonal upwards to the right. (1/sqrt(2), 1/sqrt(2))
func UnitUpRight() Vec2 {
	return Vec2{1 / math.Sqrt2, 1 / math.Sqrt2}
}

// Diagonal upwards to the left. (-1/sqrt(2), 1/sqrt(2))
func UnitUpleft() Vec2 {
	return Vec2{-1 / math.Sqrt2, 1 / math.Sqrt2}
}

// Diagonal downwards to the right. (1/sqrt(2), -1/sqrt(2))
func UnitDownRight() Vec2 {
	return Vec2{1 / math.Sqrt2, -1 / math.Sqrt2}
}

// Diagonal downwards to the left. (-1/sqrt(2), -1/sqrt(2))
func UnitDownLeft() Vec2 {
	return Vec2{-1 / math.Sqrt2, -1 / math.Sqrt2}
}

// Unit vector
func UnitFromDegrees(degrees float64) Vec2 {
	return UnitFromRadians(DegreesToRadians(degrees))
}

func UnitFromRadians(radians float64) Vec2 {
	return Vec2{math.Cos(radians), math.Sin(radians)}
}

func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func RadiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
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

// Whether u has a positive component in the
// direction from the origin to point. Perpendicular and zero vectors do not
// point towards a point.
func (u Vec2) IsPointingTowards(point Vec2) bool {
	return u.Dot(point) > 0
}

// Multiply vector by a scalar
func (u Vec2) Mul(scalar float64) Vec2 {
	return Vec2{scalar * u.X, scalar * u.Y}
}

// Length or magnitude of vector u is sqrt(u.X^2 + u.Y^2)
func (u Vec2) Length() float64 {
	return math.Sqrt(u.X*u.X + u.Y*u.Y)
}

func (u Vec2) Normalize() Vec2 {
	// Check division by zero
	length := u.Length()
	if length == 0 {
		return Zero()
	}
	return u.Mul(1 / length)
}

// Euclidean distance to another vector
func (u Vec2) DistanceTo(other Vec2) float64 {
	return other.Sub(u).Length()
}

// Gets scalar length of vectur u in the direction of another vector
func (u Vec2) LengthInDirection(direction Vec2) float64 {
	unit := direction.Normalize()
	return u.Dot(unit)
}

// Sets vector u's scalar length in the direction of other vector
// u' = u - (u*v)v + newLength*v where v is the unit vector of 'other'
func (u Vec2) SetLengthInDirection(newLength float64, direction Vec2) Vec2 {
	unit := direction.Normalize()
	return u.Sub(unit.Mul(u.Dot(unit))).Add(unit.Mul(newLength))
}

// Vector from point start to point end
func (start Vec2) To(end Vec2) Vec2 {
	return end.Sub(start)
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

// Alias of u.Dot(v)
func Dot(u, v Vec2) float64 {
	return u.Dot(v)
}

// Alias of u.Mul(scalar)
func Mul(u Vec2, scalar float64) Vec2 {
	return u.Mul(scalar)
}

// Return new vector with equal X and Y: (val, val)
func NewDiagonal(val float64) Vec2 {
	return Vec2{val, val}
}
