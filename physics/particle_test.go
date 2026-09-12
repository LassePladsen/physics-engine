package physics

import (
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

func TestParticle2DApplyForce(t *testing.T) {
	tests := []struct {
		name         string
		mass         float64
		initialAcceleration      Vec2
		forces       []Vec2
		wantAcceleration Vec2
	}{
		{"single force", 2, Vec2Zero(), []Vec2{{10, -4}}, Vec2{5, -2}},
		{"accumulates forces", 4, Vec2{1, -1}, []Vec2{{8, 12}, {-4, 4}}, Vec2{2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Particle2D{Mass: tt.mass, Acceleration: tt.initialAcceleration}
			for _, force := range tt.forces {
				p = p.ApplyForce(force)
			}
			if !p.Acceleration.ApproxEquals(tt.wantAcceleration) {
				t.Fatalf("acceleration = %v, want %v", p.Acceleration, tt.wantAcceleration)
			}
		})
	}
}

func TestParticle2DElasticCollision(t *testing.T) {
	tests := []struct {
		name                    string
		particle, other         Particle2D
		wantParticle, wantOther Vec2
	}{
		{
			name:         "equal masses exchange velocities",
			particle:     Particle2D{Mass: 2, Velocity: Vec2{3, -1}},
			other:        Particle2D{Mass: 2, Velocity: Vec2{-4, 5}},
			wantParticle: Vec2{-4, 5},
			wantOther:    Vec2{3, -1},
		},
		{
			name:         "unequal masses with stationary other particle",
			particle:     Particle2D{Mass: 1, Velocity: Vec2{6, -3}},
			other:        Particle2D{Mass: 2},
			wantParticle: Vec2{-2, 1},
			wantOther:    Vec2{4, -2},
		},
		{
			name:         "same velocity remains unchanged",
			particle:     Particle2D{Mass: 1, Velocity: Vec2{2, -7}},
			other:        Particle2D{Mass: 3, Velocity: Vec2{2, -7}},
			wantParticle: Vec2{2, -7},
			wantOther:    Vec2{2, -7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParticle, gotOther := tt.particle.ElasticCollision(tt.other)

			if !gotParticle.Velocity.ApproxEquals(tt.wantParticle) {
				t.Errorf("first particle velocity = %v, want %v", gotParticle.Velocity, tt.wantParticle)
			}
			if !gotOther.Velocity.ApproxEquals(tt.wantOther) {
				t.Errorf("second particle velocity = %v, want %v", gotOther.Velocity, tt.wantOther)
			}
		})
	}
}

func TestParticle2DElasticCollisionPanicsWhenMassesSumToZero(t *testing.T) {
	p := Particle2D{Mass: 1}
	other := Particle2D{Mass: -1}

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("ElasticCollision did not panic")
		}
		message, ok := recovered.(string)
		if !ok || !strings.Contains(message, "sum of particle masses are zero") {
			t.Fatalf("panic = %v, want zero total mass error", recovered)
		}
	}()

	p.ElasticCollision(other)
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
			if got := tt.particle.ShouldCollide(tt.other); got != tt.want {
				t.Errorf("ShouldCollide() = %t, want %t", got, tt.want)
			}
		})
	}
}
