package physics

import "math"

type Particle2D struct {
	Mass     float64 // kg
	Position Vec2    // m
	Velocity Vec2    // m/s
	Radius   float64 // m
}

// In-place advances the particle by deltaTime seconds.
func (p *Particle2D) Step(acceleration Vec2, deltaTime float64) {
	// Euler-Cromer
	p.Velocity = p.Velocity.Add(acceleration.Mul(deltaTime))
	p.Position = p.Position.Add(p.Velocity.Mul(deltaTime))
}

// Collide resolves a collision with the supplied coefficient of restitution and
// returns the two new particle states (p, other). A restitution of 1 is a
// perfectly elastic collision; 0 is perfectly inelastic.
// Panics if the sum of the masses is zero or restitution is outside [0, 1].
//
// Conservation of momentum: m_1*v_1i + m_2*v_2i = m_1*v_1f + m_2*v_2f
//
// Kinetic energy is also conserved when restitution is 1.
func (p Particle2D) Collide(other Particle2D, restitution float64) (Particle2D, Particle2D) {
	if math.IsNaN(restitution) || restitution < 0 || restitution > 1 {
		panic("Coefficient of restitution needs to be between 0 and 1, inclusive.")
	}
	sumMasses := p.Mass + other.Mass
	if sumMasses == 0 {
		panic("sum of particle masses are zero, the formulas will divide by zero.")
	}
	momentum := AddVectors(p.Velocity.Mul(p.Mass), other.Velocity.Mul(other.Mass))
	newP := p
	newP.Velocity = AddVectors(momentum, other.Velocity.Sub(p.Velocity).Mul(restitution*other.Mass))
	newP.Velocity = newP.Velocity.Mul(1 / sumMasses)

	newOther := other
	newOther.Velocity = AddVectors(momentum, p.Velocity.Sub(other.Velocity).Mul(restitution*p.Mass))
	newOther.Velocity = newOther.Velocity.Mul(1 / sumMasses)

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
	return p.DistanceTo(other) - p.Radius - other.Radius
}

// Returns whether the particles are approaching each other
func (p Particle2D) IsApproaching(other Particle2D) bool {
	separation := p.Position.Sub(other.Position)
	relativeVelocity := p.Velocity.Sub(other.Velocity)
	return separation.Dot(relativeVelocity) < 0
}
