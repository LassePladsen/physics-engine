package physics

import "github.com/LassePladsen/physics-engine/logger"

type Particle2D struct {
	Mass         float64 // kg
	Position     Vec2    // m
	Velocity     Vec2    // m/s
	Acceleration Vec2    // m/s^2
	Radius       float64 // m
}

// In-place advances the particle by deltaTime seconds.
func (p *Particle2D) Step(deltaTime float64) {
	// Euler-Cromer
	p.Velocity = p.Velocity.Add(p.Acceleration.Mul(deltaTime))
	p.Position = p.Position.Add(p.Velocity.Mul(deltaTime))
}

// Adds to particles acceleration using Newton's first law: acceleration = force / mass. Returns the new particle state
func (p Particle2D) ApplyForce(force Vec2) Particle2D {
	p.Acceleration = p.Acceleration.Add(force.Mul(1 / p.Mass))
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
	logger.Debugf("ElasticCollision: p new velocity: %v", newP.Velocity)

	newOther := other
	newOther.Velocity = p.Velocity.Mul(2 * p.Mass / sumMasses).
		Add(other.Velocity.Mul((other.Mass - p.Mass) / sumMasses))
	logger.Debugf("ElasticCollision: other new velocity: %v", newOther.Velocity)

	return newP, newOther
}

// Returns whether the particles' circumferences touch or overlap.
func (p Particle2D) IsTouching(other Particle2D) bool {
	// TODO: inelastic collision
	// TODO: make elastic vs inelastic one parameter 0 <= e <= 1.
	return p.CircumferenceDistanceTo(other) < tolerance
}

// Returns the distance from the CENTERS of the particles (not the circumference)
func (p Particle2D) DistanceTo(other Particle2D) float64 {
	return p.Position.DistanceTo(other.Position)
}

// Returns the distance from the particles' circumferences
func (p Particle2D) CircumferenceDistanceTo(other Particle2D) float64 {
	return p.DistanceTo(other)-p.Radius-other.Radius
}

// Returns whether the particles are approaching each other
func (p Particle2D) IsApproaching(other Particle2D) bool {
	separation := p.Position.Sub(other.Position)
	relativeVelocity := p.Velocity.Sub(other.Velocity)
	return separation.Dot(relativeVelocity) < 0
}
