package physics

import (
	"math"

	"github.com/LassePladsen/physics-engine/logger"
)

// World contains the particles that make up a simulation.
type World struct {
	DisableCollisions bool
	DisableBounds bool
	Particles               []Particle2D
	ParticleRestitution     float64 // coefficient of restitution 'e'. Should be 0 <= e <= 1. https://en.wikipedia.org/wiki/Coefficient_of_restitution
	WindowBoundsRestitution float64 // coefficient of restitution 'e'. Should be 0 <= e <= 1. https://en.wikipedia.org/wiki/Coefficient_of_restitution
	Friction                float64
	lastDeltaTime           float64
	Gravity                 Vec2 // Gravitatonal acceleration (m/s^2)
	ParticleToParticleGravitationalStrength float64 // m3 kg-1 s-2. newtonian gravitation Big G
	ParticleMergeEnabled bool // whether particles should absorb each other and combine to one bigger particle. Takes precedence over collisions

}

func NewWorld() World {
	return World{
		ParticleToParticleGravitationalStrength: 1,
		ParticleRestitution: 1,
		WindowBoundsRestitution: 1,
	}
}

// Advances every particle in the world by deltaTime seconds.
func (w *World) Step(deltaTime float64) {
	w.lastDeltaTime = deltaTime
	indexPairsToMerge := make(map[[2]int]bool)
	for i := range w.Particles {
		// Particle-to-particle gravity
		sumParticleGravity := Zero()
		for j := range w.Particles {
			if i == j {
				continue
			}

			if (w.ShouldMerge(i, j)) {
				// Dont add duplicate (0, 1) (1, 0) 
				a, b := canonicalPair(i, j)
				indexPairsToMerge[[2]int{a, b}] = true
			}


			sumParticleGravity = sumParticleGravity.Add(w.Particles[i].GravityFrom(w.Particles[j], w.ParticleToParticleGravitationalStrength).Mul(1/w.Particles[i].Mass))
		}
		logger.Debugf("sumParticleGravity: %v", sumParticleGravity)
		acceleration := w.Gravity.Add(sumParticleGravity)
		logger.Debugf("acceleration: %v", acceleration)
		w.Particles[i].Step(acceleration, deltaTime)

		// Collisions
		w.DoCollisions(w.ParticleRestitution) // change to inelastic here
	}
	
	// Do merges, first filter unique
	for pair := range indexPairsToMerge {
		w.MergeParticlesAt(pair[0], pair[1])
	}
}

// Merges particles at given indexes of w.Particles array
func (w *World) MergeParticlesAt(index1, index2 int) {
	logger.Debugf("LP index1: %v", index1)
	logger.Debugf("LP index2: %v", index2)
	// p1 particle becomes the newly combined particle, then remove p2 from the array
	w.Particles[index1] = w.Particles[index1].Merge(w.Particles[index2])
	w.Particles = append(w.Particles[:index2], w.Particles[index2+1:]...)
}

// Whether particles at indexes should merge
func (w *World) ShouldMerge(index1, index2 int) bool {
	if !w.ParticleMergeEnabled {
		return false
	}
	p1 := w.Particles[index1]
	p2 := w.Particles[index2]
	return p1.Overlaps(p2) || p1.IsTouching(p2)
}

// Checks and runs collisions for each relevant particle. Panics if not: 0 <= resitution <= 1 (https://en.wikipedia.org/wiki/Coefficient_of_restitution)
func (w *World) DoCollisions(restitution float64) {
	if w.DisableCollisions {
		return
	}
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

			if u.ShouldCollide(v) {
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

func canonicalPair(a, b int) (int, int) {
	if a > b {
		b, a = a, b
	}
	return a, b
}
