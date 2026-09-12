package physics

type Particle2D struct {
	Mass         float64 // kg
	Position     Vec2    // m
	Velocity     Vec2    // m/s
	Acceleration Vec2    // m/s^2
	Radius       float64 // m
}

// In-place moves the particle one step in the future, dt is the time delta change. low value is recommended for accurate results, e.g 0.001 or even lower.
// For graphics rendering at x fps it should be: 1/x
func (p *Particle2D) Step(deltaTime float64) {
	// Euler-Cromer
	p.Velocity = p.Velocity.Add(p.Acceleration.Mul(deltaTime))
	p.Position = p.Position.Add(p.Velocity.Mul(deltaTime))
}

// Adds to particles acceleration using Newton's first law: acceleration = force / mass. Returns the new particle state
func (p Particle2D) ApplyForce(force Vec2) Particle2D {
	p.Acceleration = p.Acceleration.Add(force.Mul(1/p.Mass))
	return p
}

// Performs an elastic collision i.e no kinetic energy or momentum loss, returns the two new particle states (p, other).
// Panics if the sum of their masses are zero, which should't never be the case.
//
// Conversation of momentum: m_1*v_1i + m_2*v_2i = m_1*v_1f + m_2*v_2f
//
// Conservation of kinetic energy: 1/2*m_1*v_1i^2 + 1/2*m_2*v_2i^2 = 1/2*m_1*v_1f^2 + 1/2*m_2*v_2f^2
func (p Particle2D) ElasticCollision(other Particle2D) (Particle2D, Particle2D) {
	sumMasses := p.Mass + other.Mass
	if sumMasses == 0 {
		panic("Particle2D.ElasticCollision: sum of particle masses are zero, the formulas will divide by zero.")
	}
	newP := p
	newP.Velocity = p.Velocity.Mul((p.Mass - other.Mass) / sumMasses).
		Add(other.Velocity.Mul(2 * other.Mass / sumMasses))

	newOther := other
	newOther.Velocity = p.Velocity.Mul(2 * p.Mass / sumMasses).
		Add(other.Velocity.Mul((other.Mass - p.Mass) / sumMasses))

	return newP, newOther
}

