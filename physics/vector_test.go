package physics

import (
	"math"
	"testing"
)

func TestFloatEqualsAndVec2EqualityTolerance(t *testing.T) {
	tests := []struct {
		name      string
		got, want bool
	}{
		{"equal values", FloatEquals(12.5, 12.5), true}, {"within tolerance", FloatEquals(0, tolerance/2), true},
		{"at tolerance", FloatEquals(0, tolerance), false}, {"outside tolerance", FloatEquals(0, tolerance*2), false},
		{"vector within tolerance", Vec2{1, -2}.EqualsWithTol(Vec2{1.5, -2.5}, 1), true},
		{"vector at tolerance", Vec2{1, -2}.EqualsWithTol(Vec2{2, -2}, 1), false},
		{"almost equal", Vec2Zero().ApproxEquals(Vec2{tolerance / 2, -tolerance / 2}), true},
		{"almost equal boundary", Vec2Zero().ApproxEquals(Vec2{tolerance, 0}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("result = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestVec2Arithmetic(t *testing.T) {
	tests := []struct {
		name      string
		got, want Vec2
	}{
		{"add", Vec2{0.5, 15.9}.Add(Vec2{1, 100}), Vec2{1.5, 115.9}},
		{"subtract", Vec2{0.5, 100}.Sub(Vec2{1, -50}), Vec2{-0.5, 150}},
		{"negative scalar", Vec2{-0.1, 0}.Mul(-50), Vec2{5, 0}}, {"zero scalar", Vec2{50, -10}.Mul(0), Vec2Zero()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("result = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestVectorHelpers(t *testing.T) {
	if got := AddVectors(Vec2{1, 2}, Vec2{-3, 4}, Vec2{5, -6}); got != (Vec2{3, 0}) {
		t.Errorf("AddVectors() = %v, want {3 0}", got)
	}
	if got := SubtractVectors(Vec2{1, 2}, Vec2{-3, 4}, Vec2{5, -6}); got != (Vec2{-1, 4}) {
		t.Errorf("SubtractVectors() = %v, want {-1 4}", got)
	}
}

func TestVec2DotAndLength(t *testing.T) {
	dotTests := []struct {
		name string
		u, v Vec2
		want float64
	}{
		{"perpendicular", Vec2{3, 0}, Vec2{0, 4}, 0}, {"mixed signs", Vec2{0.5, 100}, Vec2{1, -50}, -4999.5},
	}
	for _, tt := range dotTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Dot(tt.v); got != tt.want {
				t.Fatalf("dot = %v, want %v", got, tt.want)
			}
		})
	}

	lengthTests := []struct {
		name string
		u    Vec2
		want float64
	}{
		{"zero", Vec2Zero(), 0}, {"axis aligned", Vec2{-10, 0}, 10}, {"three four five", Vec2{3, 4}, 5}, {"diagonal", Vec2{1, 1}, math.Sqrt2},
	}
	for _, tt := range lengthTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Length(); !FloatEquals(got, tt.want) {
				t.Fatalf("length = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec2DistanceTo(t *testing.T) {
	tests := []struct {
		name     string
		u, other Vec2
		want     float64
	}{
		{"same point", Vec2{2, -3}, Vec2{2, -3}, 0},
		{"horizontal", Vec2{-2, 4}, Vec2{5, 4}, 7},
		{"vertical", Vec2{1, -6}, Vec2{1, 2}, 8},
		{"three four five", Vec2{-1, 2}, Vec2{2, 6}, 5},
		{"fractional mixed signs", Vec2{-0.5, 1.25}, Vec2{1, -0.75}, 2.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.DistanceTo(tt.other); !FloatEquals(got, tt.want) {
				t.Fatalf("distance = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec2Normalize(t *testing.T) {
	tests := []struct {
		name    string
		u, want Vec2
	}{
		{"zero", Vec2Zero(), Vec2Zero()}, {"right", Vec2{7, 0}, UnitRight()}, {"up", Vec2{0, 100}, UnitUp()}, {"diagonal", Vec2{-1, -1}, Vec2{-1 / math.Sqrt2, -1 / math.Sqrt2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.u.Normalize()
			if !got.ApproxEquals(tt.want) {
				t.Fatalf("normalized vector = %v, want %v", got, tt.want)
			}
			if tt.want != Vec2Zero() && !FloatEquals(got.Length(), 1) {
				t.Fatalf("normalized length = %v, want 1", got.Length())
			}
		})
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		input float64
		want  int
	}{{0, 0}, {1.49, 1}, {1.5, 2}, {-1.49, -1}, {-1.5, -2}}
	for _, tt := range tests {
		if got := Round(tt.input); got != tt.want {
			t.Errorf("Round(%v) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
