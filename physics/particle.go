package physics

type Particle2D struct {
	Mass float64
	Position Vec2
	Velocity Vec2
	Radius float64
	// TODO: acceleration
}

// Moves the particle one step in the future, dt is the time delta change. low value is recommended for accurate results, e.g 0.001 or even lower.
func (p *Particle2D) Step(dt float64) {
	// TODO: acceleration
	// Euler-Cromer
	p.Position = p.Position.Add(p.Velocity.Mul(dt))
}
