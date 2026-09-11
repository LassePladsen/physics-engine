package physics

import (
	"testing"
)

func TestParticle2Step(t *testing.T) {
	// TODO: acceleration
	movedTo := func(p Particle2D, expectedPosition Vec2, dt float64) {
		oldP := p
		p.Step(dt)
		if !p.Position.EqualsWithTol(expectedPosition) {
			t.Fatalf("Particle2d.Step moved to the wrong position: expected %v, got %v. particle=%+v", expectedPosition, p.Position, oldP)
		}
	}

	movedTo(Particle2D{Position: Vec2{1, -0.5}, Velocity: Vec2{0, 0}}, Vec2{1, -0.5}, 1)
	movedTo(Particle2D{Position: Vec2{0, 0}, Velocity: Vec2{10, -1}}, Vec2{10, -1}, 1)
	movedTo(Particle2D{Position: Vec2{0, 0}, Velocity: Vec2{10, -1}}, Vec2{10*100, -1*100}, 100)
	movedTo(Particle2D{Position: Vec2{50, -100}, Velocity: Vec2{10, 50}}, Vec2{60, -50}, 1)

	p := Particle2D{Position: Vec2{100, 5}, Velocity: Vec2Left()}
	oldP := p
	iters := 1000
	for range(iters) {
		p.Step(1)
	}
	expectedPosition := Vec2{-900, 5}
	if !p.Position.EqualsWithTol(expectedPosition) {
		t.Fatalf("Particle2d.Step moved to the wrong position after %v iterations: expected %v, got %v. particle=%+v", iters, expectedPosition, p.Position, oldP)
	}
}
