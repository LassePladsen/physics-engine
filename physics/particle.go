package physics

import (
	"math"

	"github.com/LassePladsen/physics-engine/logger"
)

type Particle2D struct {
	Mass            float64 // kg
	Position        Vec2    // m
	Velocity        Vec2    // m/s
	Radius          float64 // m
}

// In-place advances the particle by deltaTime seconds.
func (p *Particle2D) Step(acceleration Vec2, deltaTime float64) {
	// Euler-Cromer
	p.Velocity = p.Velocity.Add(acceleration.Mul(deltaTime))
	p.Position = p.Position.Add(p.Velocity.Mul(deltaTime))
}

// Collide resolves a collision with the supplied coefficient of restitution and
// returns the two new repelled particle states (p, other). A restitution of 1 is a
// perfectly elastic collision; 0 is perfectly inelastic.
// Panics if the sum of the masses is zero or restitution is outside [0, 1].
//
// Momentum is conserved. Kinetic energy is also conserved for a perfectly elastic collision.
func (p Particle2D) Collide(other Particle2D, restitution float64) (Particle2D, Particle2D) {
	newP, newOther := applyCollision(p, other, restitution)
	newP, newOther = fixOverlap(newP, newOther)

	return newP, newOther
}

// Returns the distance from the CENTERS of the particles (not the circumference)
func (p Particle2D) DistanceTo(other Particle2D) float64 {
	return p.Position.DistanceTo(other.Position)
}

// Returns the distance from the particles' circumferences
func (p Particle2D) CircumferenceDistanceTo(other Particle2D) float64 {
	return p.DistanceTo(other) - p.Radius - other.Radius
}

func (p Particle2D) SeperationSpeed(other Particle2D) float64 {
	separation := p.Position.Sub(other.Position).Normalize()
	relativeVelocity := p.Velocity.Sub(other.Velocity)
	return separation.Dot(relativeVelocity)
}

func (p Particle2D) ApproachingSpeed(other Particle2D) float64 {
	return -p.SeperationSpeed(other)
}

// Returns whether the particles are approaching each other
func (p Particle2D) IsApproaching(other Particle2D) bool {
	return p.ApproachingSpeed(other) > 0
}

// Whether the particles are inside each other
func (p Particle2D) Overlaps(other Particle2D) bool {
	return p.CircumferenceDistanceTo(other) <= 0
}

// Moves the particles outside of each others circumference to avoid new collision next frame,
// and stop them getting stuck inside each other. Returns the new particle states
//
// https://photos.google.com/u/1/documents/ChdTa2plcm1iaWxkZXIgb2cgLW9wcHRhayIJEgcKBbIBAhgOKJ7y7auKNA%3D%3D/photo/AF1QipMhAxu-BfabOGc9wwcYw4Ta3g8hkk1vuAX3w8XK
func fixOverlap(p1, p2 Particle2D) (Particle2D, Particle2D) {
	sumMasses := p1.Mass + p2.Mass
	if sumMasses == 0 {
		panic("sum of particle masses are zero, the formulas will divide by zero.")
	}
	p1DistanceToOthersCircumerfence := p1.Radius + p2.Radius - p1.DistanceTo(p2)
	logger.Debugf("p1DistanceToOtherCircumerfence: %v", p1DistanceToOthersCircumerfence)

	// This is the vector p1 needs to move for it to leave the circumerfence of other
	// But, lets move them both instead of only moving p1, move the greater mass less by using its ratio of the sum of masses
	p1BounceVector := p1.Velocity.Normalize().Mul(p1DistanceToOthersCircumerfence * p2.Mass / sumMasses)
	logger.Debugf("p1BounceVector: %v", p1BounceVector)
	otherBounceVector := p2.Velocity.Normalize().Mul(p1DistanceToOthersCircumerfence * p1.Mass / sumMasses)
	logger.Debugf("otherBounceVector: %v", otherBounceVector)

	// Now, teleport them
	p1.Position = p1.Position.Add(p1BounceVector)
	p2.Position = p2.Position.Add(otherBounceVector)

	return p1, p2
}

// Returns the new particle states after a collision
func applyCollision(p1, p2 Particle2D, restitution float64) (Particle2D, Particle2D) {
	if math.IsNaN(restitution) || restitution < 0 || restitution > 1 {
		panic("Coefficient of restitution needs to be between 0 and 1, inclusive.")
	}

	// Only collide if the centers are approaching
	separation := p1.Position.Sub(p2.Position).Normalize()
	relativeVelocity := p1.Velocity.Sub(p2.Velocity)
	separationSpeed := separation.Dot(relativeVelocity)
	logger.Debugf("Seperation speed: %v", separationSpeed)
	if separationSpeed < 0 {
		impulse := -(1 + restitution) * separationSpeed // scalar
		p1.Velocity = p1.Velocity.Add(separation.Mul(impulse / p1.Mass))
		p2.Velocity = p2.Velocity.Sub(separation.Mul(impulse / p2.Mass))
	}
	return p1, p2
}
