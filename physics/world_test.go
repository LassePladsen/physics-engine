package physics

import "testing"

func TestWorldStepUpdatesParticlesInPlace(t *testing.T) {
	world := World{Particles: []Particle2D{{
		Position: Vec2{X: 1, Y: 2},
		Velocity: Vec2{X: 4, Y: -2},
	}}}

	world.Step(0.5)

	want := Vec2{X: 3, Y: 1}
	if got := world.Particles[0].Position; got != want {
		t.Errorf("particle position after Step = %v, want %v", got, want)
	}
}

func TestWorldEnsureParticlesInBoundsReflectsVelocity(t *testing.T) {
	world := World{Particles: []Particle2D{{
		Position: Vec2{X: 10.2, Y: -0.1},
		Velocity: Vec2{X: 2, Y: -3},
		Radius:   1,
	}}}

	world.EnsureParticlesInBounds(10, 10)

	particle := world.Particles[0]
	if particle.Velocity != (Vec2{X: -2, Y: 3}) {
		t.Errorf("velocity = %v, want {-2 3}", particle.Velocity)
	}
	if particle.Position != (Vec2{X: 9, Y: 1}) {
		t.Errorf("position = %v, want {9 1}", particle.Position)
	}
}
