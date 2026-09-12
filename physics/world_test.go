package physics

import (
	"math"
	"testing"
)

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

func TestWorldStepResolvesHeadOnCollision(t *testing.T) {
	world := World{ParticleRestitution: 0.5, Particles: []Particle2D{
		{Mass: 1, Position: Vec2{X: 0, Y: 0}, Velocity: Vec2{X: 5, Y: 9.3195}, Radius: 0.2},
		{Mass: 1, Position: Vec2{X: 0.4, Y: 0}, Velocity: Vec2{X: -5, Y: 9.3195}, Radius: 0.2},
	}}

	world.Step(0)

	if got := world.Particles[0].Velocity; !got.ApproxEquals(Vec2{X: -2.5, Y: 9.3195}) {
		t.Errorf("first particle velocity = %v, want {-2.5 9.3195}", got)
	}
	if got := world.Particles[1].Velocity; !got.ApproxEquals(Vec2{X: 2.5, Y: 9.3195}) {
		t.Errorf("second particle velocity = %v, want {2.5 9.3195}", got)
	}
}

func TestWorldDoCollisionsSkipsPairsThatAreNotApproachingOrTouching(t *testing.T) {
	tests := []struct {
		name      string
		particles []Particle2D
	}{
		{
			name: "touching but moving apart",
			particles: []Particle2D{
				{Mass: 1, Position: Vec2{0, 0}, Velocity: Vec2{-1, 0}, Radius: 1},
				{Mass: 1, Position: Vec2{2, 0}, Velocity: Vec2{1, 0}, Radius: 1},
			},
		},
		{
			name: "approaching but separated",
			particles: []Particle2D{
				{Mass: 1, Position: Vec2{0, 0}, Velocity: Vec2{1, 0}, Radius: 1},
				{Mass: 1, Position: Vec2{3, 0}, Velocity: Vec2{-1, 0}, Radius: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := World{Particles: tt.particles}
			want := append([]Particle2D(nil), world.Particles...)
			world.DoCollisions(1)
			if world.Particles[0] != want[0] || world.Particles[1] != want[1] {
				t.Errorf("particles = %+v, want %+v", world.Particles, want)
			}
		})
	}
}

func TestWorldDoCollisionsPanicsForInvalidRestitution(t *testing.T) {
	for _, restitution := range []float64{-0.01, 1.01, math.NaN()} {
		t.Run("invalid restitution", func(t *testing.T) {
			world := World{}
			defer func() {
				if recover() == nil {
					t.Fatal("DoCollisions did not panic")
				}
			}()
			world.DoCollisions(restitution)
		})
	}
}

func TestWorldEnsureParticlesInBoundsReflectsVelocity(t *testing.T) {
	world := World{WindowBoundsRestitution: 1, Particles: []Particle2D{{
		Position: Vec2{X: 10.2, Y: 10.1},
		Velocity: Vec2{X: 2, Y: 3},
		Radius:   1,
	}}}

	world.EnsureParticlesInBounds(10, 10)

	particle := world.Particles[0]
	if particle.Velocity != (Vec2{X: -2, Y: -3}) {
		t.Errorf("velocity = %v, want {-2 -3}", particle.Velocity)
	}
	if particle.Position != (Vec2{X: 9, Y: 9}) {
		t.Errorf("position = %v, want {9 9}", particle.Position)
	}
}

func TestWorldEnsureParticlesInBoundsStopsParticleOnBottomBoundary(t *testing.T) {
	world := World{WindowBoundsRestitution: 0.1, Particles: []Particle2D{{
		Position:     Vec2{X: 5, Y: 1},
		Velocity:     Vec2{Y: 0},
		Acceleration: Vec2{Y: -9.81},
		Radius:       1,
	}}}

	for range 10 {
		world.Step(1.0 / 60.0)
		world.EnsureParticlesInBounds(10, 10)
	}

	particle := world.Particles[0]
	if particle.Velocity != (Vec2{Y: 0}) {
		t.Errorf("velocity = %v, want {0 0}", particle.Velocity)
	}
	if particle.Position != (Vec2{X: 5, Y: 1}) {
		t.Errorf("position = %v, want {5 1}", particle.Position)
	}
}

func TestWorldEnsureParticlesInBoundsBouncesBeforeSettlingOnGround(t *testing.T) {
	world := World{WindowBoundsRestitution: 0.1, Particles: []Particle2D{{
		Position:     Vec2{X: 5, Y: 0.9},
		Velocity:     Vec2{Y: -3},
		Acceleration: Vec2{Y: -9.81},
		Radius:       1,
	}}}

	world.Step(1.0 / 60.0)
	world.EnsureParticlesInBounds(10, 10)
	if got := world.Particles[0].Velocity.Y; !FloatEquals(got, 0.31635) {
		t.Fatalf("velocity after first ground impact = %v, want 0.31635", got)
	}

	for range 300 {
		world.Step(1.0 / 60.0)
		world.EnsureParticlesInBounds(10, 10)
	}

	if got := world.Particles[0].Velocity.Y; got != 0 {
		t.Errorf("velocity after settling = %v, want 0", got)
	}
}
