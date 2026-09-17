package physics

import (
	"testing"
	"github.com/LassePladsen/physics-engine/logger"
)

var _ = func () int {
	logger.Init()
	return 0

}()

func TestParticle2DStep(t *testing.T) {
	tests := []struct {
		name               string
		particle           Particle2D
		acceleration       Vec2
		dt                 float64
		position, velocity Vec2
	}{
		{"stationary", Particle2D{Position: Vec2{1, -0.5}}, Zero(), 1, Vec2{1, -0.5}, Zero()},
		{"constant velocity", Particle2D{Velocity: Vec2{10, -1}}, Zero(), 1, Vec2{10, -1}, Vec2{10, -1}},
		{"euler cromer acceleration", Particle2D{Position: Vec2{1, 2}, Velocity: Vec2{3, -4}}, Vec2{2, 6}, 0.5, Vec2{3, 1.5}, Vec2{4, -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.particle
			p.Step(tt.acceleration, tt.dt)
			if !p.Position.ApproxEquals(tt.position) || !p.Velocity.ApproxEquals(tt.velocity) {
				t.Fatalf("after Step: position=%v velocity=%v, want position=%v velocity=%v", p.Position, p.Velocity, tt.position, tt.velocity)
			}
		})
	}
}

func TestParticle2DOverlaps(t *testing.T) {
	tests := []struct {
		name            string
		particle, other Particle2D
		want            bool
	}{
		{
			name:     "partially overlapping circles",
			particle: Particle2D{Position: Vec2{0, 0}, Radius: 2},
			other:    Particle2D{Position: Vec2{3, 0}, Radius: 2},
			want:     true,
		},
		{
			name:     "one circle completely inside the other",
			particle: Particle2D{Position: Vec2{0, 0}, Radius: 5},
			other:    Particle2D{Position: Vec2{1, 0}, Radius: 1},
			want:     true,
		},
		{
			name:     "concentric circles",
			particle: Particle2D{Position: Vec2{2, -3}, Radius: 4},
			other:    Particle2D{Position: Vec2{2, -3}, Radius: 1},
			want:     true,
		},
		{
			name:     "circles are separated",
			particle: Particle2D{Position: Vec2{0, 0}, Radius: 2},
			other:    Particle2D{Position: Vec2{5.1, 0}, Radius: 3},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.particle.Overlaps(tt.other); got != tt.want {
				t.Errorf("Overlaps() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestParticle2DIsApproaching(t *testing.T) {
	tests := []struct {
		name            string
		particle, other Particle2D
		want            bool
	}{
		{
			name:     "both approaching",
			particle: Particle2D{Position: Vec2{0, 0}, Velocity: Vec2{3, 0}},
			other:    Particle2D{Position: Vec2{3, 0}, Velocity: Vec2{-1, 0}},
			want:     true,
		},
		{
			name:     "separating one chasing the other",
			particle: Particle2D{Position: Vec2{0, 0}, Velocity: Vec2{3, 0}},
			other:    Particle2D{Position: Vec2{3, 0}, Velocity: Vec2{5, 0}},
			want:     false,
		},
		{
			name:     "approaching one chasing the other",
			particle: Particle2D{Position: Vec2{0, 0}, Velocity: Vec2{3, 0}},
			other:    Particle2D{Position: Vec2{3, 0}, Velocity: Vec2{1, 0}},
			want:     true,
		},
		{
			name:     "one stationary separating",
			particle: Particle2D{Position: Vec2{0, 0}, Velocity: Vec2{0, 0}},
			other:    Particle2D{Position: Vec2{3, 0}, Velocity: Vec2{1, 0}},
			want:     false,
		},
		{
			name:     "one stationary approaching",
			particle: Particle2D{Position: Vec2{0, 0}, Velocity: Vec2{0, 0}},
			other:    Particle2D{Position: Vec2{3, 0}, Velocity: Vec2{-1, 0}},
			want:     true,
		},
		{
			name:     "both stationary",
			particle: Particle2D{Position: Vec2{0, 0}, Velocity: Vec2{0, 0}},
			other:    Particle2D{Position: Vec2{3, 0}, Velocity: Vec2{0, 0}},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.particle.IsApproaching(tt.other); got != tt.want {
				t.Errorf("IsApproaching() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestParticle2DCollide(t *testing.T) {
	tests := []struct {
		name          string
		p, other      Particle2D
		restitution   float64
		wantPVelocity Vec2
		wantOVelocity Vec2
	}{
		{
			name: "equal masses head-on elastic",
			p: Particle2D{
				Mass:     1,
				Position: Vec2{X: 0, Y: 0},
				Velocity: Vec2{X: 1, Y: 0},
				Radius:   1,
			},
			other: Particle2D{
				Mass:     1,
				Position: Vec2{X: 1.5, Y: 0},
				Velocity: Vec2{X: -1, Y: 0},
				Radius:   1,
			},
			restitution:   1,
			wantPVelocity: Vec2{X: -1, Y: 0},
			wantOVelocity: Vec2{X: 1, Y: 0},
		},
		{
			name: "different masses elastic",
			p: Particle2D{
				Mass:     2,
				Position: Vec2{X: 0, Y: 0},
				Velocity: Vec2{X: 3, Y: 0},
				Radius:   1,
			},
			other: Particle2D{
				Mass:     1,
				Position: Vec2{X: 1.5, Y: 0},
				Velocity: Vec2{X: 0, Y: 0},
				Radius:   1,
			},
			restitution:   1,
			wantPVelocity: Vec2{X: 1, Y: 0},
			wantOVelocity: Vec2{X: 4, Y: 0},
		},
		{
			name: "equal masses perfectly inelastic",
			p: Particle2D{
				Mass:     1,
				Position: Vec2{X: 0, Y: 0},
				Velocity: Vec2{X: 2, Y: 0},
				Radius:   1,
			},
			other: Particle2D{
				Mass:     1,
				Position: Vec2{X: 1.5, Y: 0},
				Velocity: Vec2{X: 0, Y: 0},
				Radius:   1,
			},
			restitution:   0,
			wantPVelocity: Vec2{X: 1, Y: 0},
			wantOVelocity: Vec2{X: 1, Y: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP, gotOther := tt.p.Collide(tt.other, tt.restitution)

			if !gotP.Velocity.ApproxEquals(tt.wantPVelocity) {
				t.Errorf("p velocity = %v, want %v",
					gotP.Velocity, tt.wantPVelocity)
			}

			if !gotOther.Velocity.ApproxEquals(tt.wantOVelocity) {
				t.Errorf("other velocity = %v, want %v",
					gotOther.Velocity, tt.wantOVelocity)
			}

			// Momentum must be conserved.
			beforeMomentum := tt.p.Momentum().Add(tt.other.Momentum())
			afterMomentum := gotP.Momentum().Add(gotOther.Momentum())

			if !beforeMomentum.ApproxEquals(afterMomentum) {
				t.Errorf(
					"momentum not conserved: before=%v, after=%v",
					beforeMomentum, afterMomentum,
				)
			}

			// Kinetic energy must be conserved for a perfectly
			// elastic collision.
			if tt.restitution == 1 {
				beforeKE := tt.p.KineticEnergy() + tt.other.KineticEnergy()
				afterKE := gotP.KineticEnergy() + gotOther.KineticEnergy()

				if !FloatEquals(beforeKE, afterKE) {
					t.Errorf(
						"kinetic energy not conserved: before=%v, after=%v",
						beforeKE, afterKE,
					)
				}
			}
		})
	}
}

func TestParticle2DCollidePanics(t *testing.T) {
	base := Particle2D{
		Mass:     1,
		Position: Vec2{X: 0, Y: 0},
		Velocity: Vec2{X: 1, Y: 0},
		Radius:   1,
	}

	tests := []struct {
		name        string
		p, other    Particle2D
		restitution float64
	}{
		{
			name:        "restitution below zero",
			p:           base,
			other:       base,
			restitution: -0.01,
		},
		{
			name:        "restitution above one",
			p:           base,
			other:       base,
			restitution: 1.01,
		},
		{
			name: "zero total mass",
			p: Particle2D{
				Mass: 0,
			},
			other: Particle2D{
				Mass: 0,
			},
			restitution: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("Collide did not panic")
				}
			}()

			tt.p.Collide(tt.other, tt.restitution)
		})
	}
}
