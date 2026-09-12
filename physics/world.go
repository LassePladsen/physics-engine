package physics

import (
	"math"

	"github.com/LassePladsen/physics-engine/logger"
)

// World contains the particles that make up a simulation.
type World struct {
	Particles []Particle2D
	Restitution float64  // coefficient of restitution 'e'. Should be 0 <= e <= 1. https://en.wikipedia.org/wiki/Coefficient_of_restitution
}

// Advances every particle in the world by deltaTime seconds.
func (w *World) Step(deltaTime float64) {
	for i := range w.Particles {
		w.Particles[i].Step(deltaTime)
	}
	w.DoCollisions(w.Restitution) // change to inelastic here
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

			if u.IsTouching(v) && u.IsApproaching(v) {
				logger.Debug("They should collide")
				w.Particles[i], w.Particles[j] = u.Collide(v, restitution)
			} else {
				logger.Debug("NO collision")
			}
		}
	}
}

// EnsureParticlesInBounds keeps every particle inside a world with the given
// dimensions, reflecting its velocity when it reaches an edge.
func (w *World) EnsureParticlesInBounds(width, height float64) {
	for i := range w.Particles {
		ensureParticleInBounds(&w.Particles[i], width, height)
	}
}

func ensureParticleInBounds(particle *Particle2D, width, height float64) {
	if particle.Position.X+particle.Radius > width ||
		particle.Position.X-particle.Radius < 0 {
		particle.Velocity.X *= -1
		particle.Position.X = min(max(particle.Radius, particle.Position.X+particle.Radius), width-particle.Radius)
	}
	if particle.Position.Y+particle.Radius > height ||
		particle.Position.Y-particle.Radius < 0 {
		particle.Velocity.Y *= -1
		particle.Position.Y = min(max(particle.Radius, particle.Position.Y+particle.Radius), height-particle.Radius)
	}
}
