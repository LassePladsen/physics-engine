package physics

import (
	"math"
	"strings"
	"testing"
)

func TestParticle2DStep(t *testing.T) {
	tests := []struct {
		name               string
		particle           Particle2D
		dt                 float64
		position, velocity Vec2
	}{
		{"stationary", Particle2D{Position: Vec2{1, -0.5}}, 1, Vec2{1, -0.5}, Vec2Zero()},
		{"constant velocity", Particle2D{Velocity: Vec2{10, -1}}, 1, Vec2{10, -1}, Vec2{10, -1}},
		{"euler cromer acceleration", Particle2D{Position: Vec2{1, 2}, Velocity: Vec2{3, -4}, Acceleration: Vec2{2, 6}}, 0.5, Vec2{3, 1.5}, Vec2{4, -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.particle
			p.Step(tt.dt)
			if !p.Position.ApproxEquals(tt.position) || !p.Velocity.ApproxEquals(tt.velocity) {
				t.Fatalf("after Step: position=%v velocity=%v, want position=%v velocity=%v", p.Position, p.Velocity, tt.position, tt.velocity)
			}
		})
	}
}

func TestParticle2DCollide(t *testing.T) {
	tests := []struct {
		name                    string
		particle, other         Particle2D
		restitution             float64
		wantParticle, wantOther Vec2
	}{
		{"elastic equal masses exchange velocities", Particle2D{Mass: 2, Velocity: Vec2{3, -1}}, Particle2D{Mass: 2, Velocity: Vec2{-4, 5}}, 1, Vec2{-4, 5}, Vec2{3, -1}},
		{"elastic unequal masses with stationary particle", Particle2D{Mass: 1, Velocity: Vec2{6, -3}}, Particle2D{Mass: 2}, 1, Vec2{-2, 1}, Vec2{4, -2}},
		{"perfectly inelastic uses center of mass velocity", Particle2D{Mass: 1, Velocity: Vec2{6, -3}}, Particle2D{Mass: 2}, 0, Vec2{2, -1}, Vec2{2, -1}},
		{"partial restitution", Particle2D{Mass: 1, Velocity: Vec2{6, -3}}, Particle2D{Mass: 2}, 0.5, Vec2{0, 0}, Vec2{3, -1.5}},
		{"same velocity remains unchanged", Particle2D{Mass: 1, Velocity: Vec2{2, -7}}, Particle2D{Mass: 3, Velocity: Vec2{2, -7}}, 1, Vec2{2, -7}, Vec2{2, -7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParticle, gotOther := tt.particle.Collide(tt.other, tt.restitution)
			if !gotParticle.Velocity.ApproxEquals(tt.wantParticle) {
				t.Errorf("first particle velocity = %v, want %v", gotParticle.Velocity, tt.wantParticle)
			}
			if !gotOther.Velocity.ApproxEquals(tt.wantOther) {
				t.Errorf("second particle velocity = %v, want %v", gotOther.Velocity, tt.wantOther)
			}
		})
	}
}

func TestParticle2DCollidePreservesOtherStateAndInputs(t *testing.T) {
	p := Particle2D{Mass: 1, Position: Vec2{1, 2}, Velocity: Vec2{6, -3}, Acceleration: Vec2{4, 5}, Radius: 7}
	other := Particle2D{Mass: 2, Position: Vec2{8, 9}, Velocity: Vec2{-1, 2}, Acceleration: Vec2{-6, 3}, Radius: 4}

	gotP, gotOther := p.Collide(other, 1)
	if p.Velocity != (Vec2{6, -3}) || other.Velocity != (Vec2{-1, 2}) {
		t.Fatal("Collide mutated one of its inputs")
	}
	if gotP.Mass != p.Mass || gotP.Position != p.Position || gotP.Acceleration != p.Acceleration || gotP.Radius != p.Radius {
		t.Errorf("first particle state besides velocity changed: %+v", gotP)
	}
	if gotOther.Mass != other.Mass || gotOther.Position != other.Position || gotOther.Acceleration != other.Acceleration || gotOther.Radius != other.Radius {
		t.Errorf("second particle state besides velocity changed: %+v", gotOther)
	}
}

func TestParticle2DCollidePanicsForInvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		p, other    Particle2D
		restitution float64
		message     string
	}{
		{"zero total mass", Particle2D{Mass: 1}, Particle2D{Mass: -1}, 1, "sum of particle masses are zero"},
		{"negative restitution", Particle2D{Mass: 1}, Particle2D{Mass: 1}, -0.01, "Coefficient of restitution"},
		{"restitution above one", Particle2D{Mass: 1}, Particle2D{Mass: 1}, 1.01, "Coefficient of restitution"},
		{"not a number restitution", Particle2D{Mass: 1}, Particle2D{Mass: 1}, math.NaN(), "Coefficient of restitution"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				recovered := recover()
				message, ok := recovered.(string)
				if !ok || !strings.Contains(message, tt.message) {
					t.Fatalf("panic = %v, want message containing %q", recovered, tt.message)
				}
			}()
			tt.p.Collide(tt.other, tt.restitution)
		})
	}
}

func TestParticle2DShouldCollide(t *testing.T) {
	tests := []struct {
		name            string
		particle, other Particle2D
		want            bool
	}{
		{
			name:     "overlapping particles",
			particle: Particle2D{Position: Vec2{0, 0}, Radius: 2},
			other:    Particle2D{Position: Vec2{3, 0}, Radius: 2},
			want:     true,
		},
		{
			name:     "particles touching at their circumferences",
			particle: Particle2D{Position: Vec2{0, 0}, Radius: 2},
			other:    Particle2D{Position: Vec2{5, 0}, Radius: 3},
			want:     true,
		},
		{
			name:     "particles separated beyond their radii",
			particle: Particle2D{Position: Vec2{0, 0}, Radius: 2},
			other:    Particle2D{Position: Vec2{5.1, 0}, Radius: 3},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.particle.IsTouching(tt.other); got != tt.want {
				t.Errorf("ShouldCollide() = %t, want %t", got, tt.want)
			}
		})
	}
}
