package physics

import "testing"

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
		initial      Vec2
		forces       []Vec2
		acceleration Vec2
	}{
		{"single force", 2, Vec2Zero(), []Vec2{{10, -4}}, Vec2{5, -2}},
		{"accumulates forces", 4, Vec2{1, -1}, []Vec2{{8, 12}, {-4, 4}}, Vec2{2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Particle2D{Mass: tt.mass, Acceleration: tt.initial}
			for _, force := range tt.forces {
				p.ApplyForce(force)
			}
			if !p.Acceleration.ApproxEquals(tt.acceleration) {
				t.Fatalf("acceleration = %v, want %v", p.Acceleration, tt.acceleration)
			}
		})
	}
}
