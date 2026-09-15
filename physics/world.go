package physics

import (
	"math"

	"github.com/LassePladsen/physics-engine/logger"
)

// World contains the particles that make up a simulation.
type World struct {
	DisableCollisions bool
	Particles               []Particle2D
	ParticleRestitution     float64 // coefficient of restitution 'e'. Should be 0 <= e <= 1. https://en.wikipedia.org/wiki/Coefficient_of_restitution
	WindowBoundsRestitution float64 // coefficient of restitution 'e'. Should be 0 <= e <= 1. https://en.wikipedia.org/wiki/Coefficient_of_restitution
	Friction                float64
	lastDeltaTime           float64
	Gravity                 Vec2 // Gravitatonal acceleration (m/s^2)
}

// Advances every particle in the world by deltaTime seconds.
func (w *World) Step(deltaTime float64) {
	w.lastDeltaTime = deltaTime
	for i := range w.Particles {
		w.Particles[i].Step(w.Gravity, deltaTime)
	}
	if !w.DisableCollisions {
		w.DoCollisions(w.ParticleRestitution) // change to inelastic here
	}
}

// Checks and runs collisions for each relevant particle. Panics if not: 0 <= resitution <= 1 (https://en.wikipedia.org/wiki/Coefficient_of_restitution)
func (w *World) DoCollisions(restitution float64) {
	if math.IsNaN(restitution) || restitution > 1 || restitution < 0 {
		panic("Coefficient of restitution needs to be between 0 and 1, inclusive.")
	}
	for i := range w.Particles {
		u := w.Particles[i]
		// Each pair is unordered, so resolve it exactly once. Resolving both
		// (i, j) and (j, i) reverses a head-on elastic collision immediately.
		for j := i + 1; j < len(w.Particles); j++ {
			v := w.Particles[j]
			logger.Debugf("Checking collision for %+v and %+v", u, v)

			if u.Overlaps(v) {
				logger.Debug("They should collide")
				w.Particles[i], w.Particles[j] = u.Collide(v, restitution)
			} else {
				logger.Debug("NO collision")
			}
		}
	}
}

// keeps every particle inside a world with the given boundary dimensions
func (w *World) EnsureParticlesInBoundary(width, height float64) {
	for i := range w.Particles {
		w.ensureParticleInBoundary(&w.Particles[i], width, height)
	}
}

func (w *World) ensureParticleInBoundary(particle *Particle2D, width, height float64) {
	var hitBoundaryNormal Vec2
	hitBoundary := false

	// Right bounce
	if particle.Position.X+particle.Radius > width {
		particle.Velocity.X *= -w.WindowBoundsRestitution
		particle.Position.X = width - particle.Radius
		normal := UnitLeft()
		w.applyBoundaryFriction(particle, normal)
		if w.Gravity.Dot(normal) < 0 {
			hitBoundaryNormal, hitBoundary = normal, true
		}
	}

	// Left bounce
	if particle.Position.X-particle.Radius < 0 {
		particle.Velocity.X *= -w.WindowBoundsRestitution
		particle.Position.X = particle.Radius
		normal := UnitRight()
		w.applyBoundaryFriction(particle, normal)
		if w.Gravity.Dot(normal) < 0 {
			hitBoundaryNormal, hitBoundary = normal, true
		}
	}

	// Top bounce
	if particle.Position.Y+particle.Radius > height {
		particle.Velocity.Y *= -w.WindowBoundsRestitution
		particle.Position.Y = height - particle.Radius //
		normal := UnitDown()
		w.applyBoundaryFriction(particle, normal)
		if w.Gravity.Dot(normal) < 0 {
			hitBoundaryNormal, hitBoundary = normal, true
		}
	}

	// Ground bounce
	if particle.Position.Y-particle.Radius < 0 {
		particle.Velocity.Y *= -w.WindowBoundsRestitution
		particle.Position.Y = particle.Radius
		normal := UnitUp()
		w.applyBoundaryFriction(particle, normal)
		if w.Gravity.Dot(normal) < 0 {
			hitBoundaryNormal, hitBoundary = normal, true
		}
	}

	// Fix endless microbouncing: A rebound that is smaller than gravity adds in one tick cannot lift the
	// particle off the wall/ground, so treat it as resting instead of bounce
	if hitBoundary {
		normalVelocity := particle.Velocity.Dot(hitBoundaryNormal)
		gravityPerTick := -w.Gravity.Dot(hitBoundaryNormal) * w.lastDeltaTime
		if normalVelocity <= gravityPerTick {
			particle.Velocity = particle.Velocity.SetLengthInDirection(0, hitBoundaryNormal)
		}
	}
}

// applyBoundaryFriction applies kinetic friction for a boundary with the given
// inward normal. A boundary supports a particle only when gravity pushes it
// into that boundary; its normal force is the gravity component along normal.
func (w *World) applyBoundaryFriction(particle *Particle2D, normal Vec2) {
	normalGravity := w.Gravity.Dot(normal)
	if normalGravity >= 0 {
		return
	}

	// Friction removes only the velocity parallel to the boundary. It is capped
	// at the tangential speed, so it cannot reverse the particle's direction.
	tangentialVelocity := particle.Velocity.Sub(normal.Mul(particle.Velocity.Dot(normal)))
	frictionDeltaV := min(w.Friction*-normalGravity*w.lastDeltaTime, tangentialVelocity.Length())
	particle.Velocity = particle.Velocity.Sub(tangentialVelocity.Normalize().Mul(frictionDeltaV))
}

// Parameters for random world generation
type WorldGenConfig struct {
	NumParticles  int
	MinMass       float64 // kg
	MaxMass       float64 // kg
	MinPositions  Vec2    // m
	MaxPositions  Vec2    // m
	MinVelocities Vec2    // m / s
	MaxVelocities Vec2    // m / s
	MinRadius     float64 // m
	MaxRadius     float64 // m
}

// Randomizes a world with given number of particles
func GenerateRandom(config WorldGenConfig) World {
	var world World
	for range config.NumParticles {
		particle := Particle2D{
			Position: Vec2{RandomFloat(config.MinPositions.X, config.MaxPositions.X), RandomFloat(config.MinPositions.X, config.MaxPositions.Y)},
			Velocity: Vec2{RandomFloat(config.MinVelocities.X, config.MaxVelocities.X), RandomFloat(config.MinVelocities.Y, config.MaxVelocities.Y)},
			Radius:   RandomFloat(config.MinRadius, config.MaxRadius),
			Mass:     RandomFloat(config.MinMass, config.MaxMass),
		}
		// Check boundaries
		particle.Position.X = Clamp(particle.Position.X, config.MinPositions.X + particle.Radius*1.1, config.MaxPositions.X-particle.Radius*1.1)
		particle.Position.Y = Clamp(particle.Position.Y, config.MinPositions.Y + particle.Radius*1.1, config.MaxPositions.Y-particle.Radius*1.1)
		world.Particles = append(world.Particles, particle)

	}
	return world
}
