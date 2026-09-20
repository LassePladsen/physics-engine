package physics

import (
	"math"
	"testing"
)

func TestGenerateRandomCreatesParticlesWithinConfiguredRanges(t *testing.T) {
	config := WorldGenConfig{
		NumParticles:  100,
		MinMass:       2,
		MaxMass:       4,
		MinPositions:  Vec2{X: -10, Y: 20},
		MaxPositions:  Vec2{X: 30, Y: 60},
		MinVelocities: Vec2{X: -5, Y: -3},
		MaxVelocities: Vec2{X: 7, Y: 9},
		MinRadius:     1,
		MaxRadius:     3,
	}

	world := GenerateRandom(config)
	if got := len(world.Particles); got != config.NumParticles {
		t.Fatalf("number of particles = %d, want %d", got, config.NumParticles)
	}

	for i, particle := range world.Particles {
		if particle.Mass < config.MinMass || particle.Mass > config.MaxMass {
			t.Errorf("particle %d mass = %v, want within [%v, %v]", i, particle.Mass, config.MinMass, config.MaxMass)
		}
		if particle.Radius < config.MinRadius || particle.Radius > config.MaxRadius {
			t.Errorf("particle %d radius = %v, want within [%v, %v]", i, particle.Radius, config.MinRadius, config.MaxRadius)
		}
		if particle.Velocity.X < config.MinVelocities.X || particle.Velocity.X > config.MaxVelocities.X ||
			particle.Velocity.Y < config.MinVelocities.Y || particle.Velocity.Y > config.MaxVelocities.Y {
			t.Errorf("particle %d velocity = %v, want within %v to %v", i, particle.Velocity, config.MinVelocities, config.MaxVelocities)
		}

		minX, maxX := config.MinPositions.X+particle.Radius*1.1, config.MaxPositions.X-particle.Radius*1.1
		minY, maxY := config.MinPositions.Y+particle.Radius*1.1, config.MaxPositions.Y-particle.Radius*1.1
		if particle.Position.X < minX || particle.Position.X > maxX || particle.Position.Y < minY || particle.Position.Y > maxY {
			t.Errorf("particle %d position = %v, want within {%v %v} to {%v %v}", i, particle.Position, minX, minY, maxX, maxY)
		}
	}
}

func TestGenerateRandomWithZeroParticles(t *testing.T) {
	world := GenerateRandom(WorldGenConfig{})
	if len(world.Particles) != 0 {
		t.Errorf("number of particles = %d, want 0", len(world.Particles))
	}
}

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
		t.Errorf("first particle velocity = %v, want {-2.5 9.3195}. Position = %v", got, world.Particles[0].Position)
	}
	if got := world.Particles[1].Velocity; !got.ApproxEquals(Vec2{X: 2.5, Y: 9.3195}) {
		t.Errorf("second particle velocity = %v, want {2.5 9.3195}. Position = %v", got, world.Particles[1].Position)
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

func TestWorldEnsureParticlesInBoundsStopsParticleOnBottomBoundary(t *testing.T) {
	world := World{Gravity: Vec2{Y: -9.81}, WindowBoundsRestitution: 0.1, Particles: []Particle2D{{
		Position: Vec2{X: 5, Y: 1},
		Velocity: Vec2{Y: 0},
		Radius:   1,
	}}}

	for range 10 {
		world.Step(1.0 / 60.0)
		world.EnsureParticlesInBoundary(10, 10)
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
	world := World{Gravity: Vec2{Y: -9.81}, WindowBoundsRestitution: 0.1, Particles: []Particle2D{{
		Position: Vec2{X: 5, Y: 0.9},
		Velocity: Vec2{Y: -3},
		Radius:   1,
	}}}

	world.Step(1.0 / 60.0)
	world.EnsureParticlesInBoundary(10, 10)
	if got := world.Particles[0].Velocity.Y; !FloatEquals(got, 0.31635) {
		t.Fatalf("velocity after first ground impact = %v, want 0.31635", got)
	}

	for range 300 {
		world.Step(1.0 / 60.0)
		world.EnsureParticlesInBoundary(10, 10)
	}

	if got := world.Particles[0].Velocity.Y; got != 0 {
		t.Errorf("velocity after settling = %v, want 0", got)
	}
}

func TestWorldEnsureParticlesInBoundsDoesNotCancelAirborneGravity(t *testing.T) {
	world := World{
		Gravity: Vec2{Y: -9.81},
		Particles: []Particle2D{{
			Position: Vec2{X: 5, Y: 5},
			Radius:   1,
		}},
	}

	world.Step(1.0 / 60.0)
	world.EnsureParticlesInBoundary(10, 10)

	if got, want := world.Particles[0].Velocity.Y, -9.81/60.0; !FloatEquals(got, want) {
		t.Errorf("airborne vertical velocity = %v, want %v", got, want)
	}
}

func TestWorldEnsureParticlesInBoundsAppliesGroundFrictionWithoutReversing(t *testing.T) {
	tests := []struct {
		name     string
		gravity  Vec2
		position Vec2
		velocity Vec2
		want     Vec2
	}{
		{
			name:     "slows velocity tangent to downward gravity",
			gravity:  Vec2{Y: -10},
			position: Vec2{X: 5, Y: 0.9},
			velocity: Vec2{X: 3, Y: 2},
			want:     Vec2{X: 0.5, Y: 0},
		},
		{
			name:     "slows velocity tangent to horizontal gravity at a wall",
			gravity:  Vec2{X: 10},
			position: Vec2{X: 9.1, Y: 5},
			velocity: Vec2{X: 2, Y: 3},
			want:     Vec2{X: 0, Y: 0.5},
		},
		{
			name:     "stops instead of reversing",
			gravity:  Vec2{Y: -10},
			position: Vec2{X: 5, Y: 0.9},
			velocity: Vec2{X: 1},
			want:     Vec2{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := World{
				Friction: 0.5,
				Gravity:  tt.gravity,
				Particles: []Particle2D{{
					Position: tt.position,
					Velocity: tt.velocity,
					Radius:   1,
				}},
				lastDeltaTime: 0.5,
			}

			world.EnsureParticlesInBoundary(10, 10)
			if got := world.Particles[0].Velocity; !got.ApproxEquals(tt.want) {
				t.Errorf("velocity = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO:
// func TestWorldAssimilate(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		particles     []Particle2D
// 		wantParticles []Particle2D
// 	}{
// 		{
// 			name: "assimilate two equal particles",
// 			particles: []Particle2D{
// 				{Mass: 1, Position: Vec2{X: 0, Y: 0}, Velocity: Vec2{X: 5, Y: 9.3195}, Radius: 0.2},
// 				{Mass: 1, Position: Vec2{X: 0.1, Y: 0}, Velocity: Vec2{X: -5, Y: 9.3195}, Radius: 0.2},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			world := World{
// 				Particles:     tt.particles,
// 				lastDeltaTime: 0.5,
// 			}

// 			world.EnsureParticlesInBoundary(10, 10)
// 			if got := world.Particles[0].Velocity; !got.ApproxEquals(tt.wantNumParticles) {
// 				t.Errorf("velocity = %v, want %v", got, tt.wantNumParticles)
// 			}
// 		})
// 	}
// }
