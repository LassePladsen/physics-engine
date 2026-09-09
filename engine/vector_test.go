package engine

import (
	"testing"
	"math"
)

func TestVec2EqualsWithTol(t *testing.T) {
	ensure := func(u, v Vec2, expected bool) {
		got := u.EqualsWithTol(v)
		if expected != got {
			t.Fatalf("Vec2.EqualsWithTol failed: expected %+v, got %+v. u=%+v, v=%+v", expected, got, u, v)
		}
	}
	ensure(Vec2{1.5, -10}, Vec2{1.499999999999999999999999999999999999999999999999999999, -10.000000000000000000000000000000000000000000000001}, true)
	ensure(Vec2{10, -100}.Add(Vec2{-20, 100}), Vec2{-10, 0}, true)
	ensure(Vec2{10, -100}.Add(Vec2{-20.1, 100}), Vec2{-10, 0}, false)
}

func TestVec2Add(t *testing.T) {
	equals := func(u, v, expected Vec2) {
		got := u.Add(v)
		if got != expected {
			t.Fatalf("Vec2.Add failed: expected %+v, got %+v. u=%+v, v=%+v", expected, got, u, v)
		}
	}

	equals(Vec2{0.5, 15.9}, Vec2{1, 100}, Vec2{1.5, 115.9})
	equals(Vec2{-50, 0}, Vec2{-100, 0}, Vec2{-150, 0})
}

func TestVec2Sub(t *testing.T) {
	equals := func(u, v, expected Vec2) {
		got := u.Sub(v)
		if got != expected {
			t.Fatalf("Vec2.Sub failed: expected %+v, got %+v. u=%+v, v=%+v", expected, got, u, v)
		}
	}

	equals(Vec2{0.5, 100}, Vec2{1, -50}, Vec2{-0.5, 150})
	equals(Vec2{-10, 0}, Vec2{-100, -5}, Vec2{90, 5})
}


func TestVec2Dot(t *testing.T) {
	equals := func(u, v Vec2, expected float64) {
		got := u.Dot(v)
		if got != expected {
			t.Fatalf("Vec2.Dot failed: expected %+v, got %+v. u=%+v, v=%+v", expected, got, u, v)
		}
	}

	equals(Vec2{0.5, 100}, Vec2{1, -50}, 0.5*1 + 100*(-50))
	equals(Vec2{-0.1, 0}, Vec2{-100, -5}, -0.1*(-100) + 0*(-5))
}

func TestVec2Mul(t *testing.T) {
	equals := func(u Vec2, scalar float64, expected Vec2) {
		got := u.Mul(scalar)
		if got != expected {
			t.Fatalf("Vec2.Mul failed: expected %+v, got %+v. u=%+v, scalar=%+v", expected, got, u, scalar)
		}
	}

	equals(Vec2{0.5, 100}, 5, Vec2{2.5, 500})
	equals(Vec2{-0.1, 0}, -50, Vec2{5, 0})
	equals(Vec2{50, -10}, 0, Vec2{0, 0})
	equals(Vec2{50, -10}, -0, Vec2{0, 0})
	equals(Vec2{50, -10}, 1, Vec2{50, -10})
	equals(Vec2{50, -10}, -1, Vec2{-50, 10})
}

func TestVec2Length(t *testing.T) {
	tol := 0.001
	equals := func(u Vec2, expected float64) {
		got := u.Length()
		if math.Abs(got - expected) >= tol {
			t.Fatalf("Vec2.Length failed: expected %+v, got %+v. u=%+v", expected, got, u)
		}
	}

	equals(Vec2{1, 1}, math.Sqrt(2))
	equals(Vec2{0, 0}, 0)
	equals(Vec2{100, 0}, 100)
	equals(Vec2{-10, -10}, 14.142)
	equals(Vec2{10, -10}, 14.142)
	equals(Vec2{-10, 10}, 14.142)
}

func TestVec2Normalize(t *testing.T) {
	equals := func(u Vec2, expected Vec2) {
		got := u.Normalize()
		if !expected.EqualsWithTol(got){
			t.Fatalf("Vec2.Normalize failed: expected %+v, got %+v. u=%+v", expected, got, u)
		}
	}

	equals(Vec2{1, 1}, Vec2{math.Sqrt(1), math.Sqrt(1)})
	equals(Vec2{-15150, 0}, Vec2Left())
	equals(Vec2{0, 100}, Vec2Up())
	equals(Vec2{-10, -10}, Vec2{-math.Sqrt(1), -math.Sqrt(1)})
}
