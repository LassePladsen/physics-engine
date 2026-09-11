package physics

type Particle2D struct {
	Mass float64
	Position Vec2
	Velocity Vec2
	Acceleration Vec2
	Radius float64
}

// Moves the particle one step in the future, dt is the time delta change. low value is recommended for accurate results, e.g 0.001 or even lower.
func (p *Particle2D) Step(dt float64) {
	// Euler-Cromer
	p.Velocity = p.Velocity.Add(p.Acceleration.Mul(dt))
	p.Position = p.Position.Add(p.Velocity.Mul(dt))
}

// Adds to particles acceleration by dividing force by mass
func (p * Particle2D) ApplyForce(force Vec2) {
	p.Acceleration.Add(force)
}
