package physics

import (
	"math"

	"github.com/LassePladsen/physics-engine/logger"
)

// World contains the particles that make up a simulation.
type World struct {
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
	w.DoCollisions(w.ParticleRestitution) // change to inelastic here
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
// dimensions. Side walls and the ceiling reflect velocity; the bottom boundary
// is a resting surface, so it removes downward velocity. Physics coordinates
// use positive Y upwards.
func (w *World) EnsureParticlesInBounds(width, height float64) {
	for i := range w.Particles {
		w.ensureParticleInBounds(&w.Particles[i], width, height)
	}
}

func (w *World) ensureParticleInBounds(particle *Particle2D, width, height float64) {
	// Left/right bounce
	if particle.Position.X+particle.Radius > width ||
		particle.Position.X-particle.Radius < 0 {
		particle.Velocity.X *= -w.WindowBoundsRestitution
		particle.Position.X = min(max(particle.Radius, particle.Position.X), width-particle.Radius)
	}

	// Top bounce
	if particle.Position.Y+particle.Radius > height {
		particle.Velocity.Y *= -w.WindowBoundsRestitution
		particle.Position.Y = height - particle.Radius //
	}

	// Groundbounce
	if particle.Position.Y-particle.Radius < 0 {
		particle.Velocity.Y *= -w.WindowBoundsRestitution
		particle.Position.Y = particle.Radius
	}

	// Apply friction by dampening x speed when on the ground. F = mu * g
	// CODEX HERE

	// A rebound that is smaller than gravity adds in one tick cannot lift the
	// particle off the wall/ground, so treat it as resting instead of endless micro-bouncing.
	v_g := particle.Velocity.LengthInDirection(w.Gravity)
	if v_g <= -w.Gravity.Length()*w.lastDeltaTime {
		particle.Velocity = particle.Velocity.SetLengthInDirection(0, w.Gravity)
	}
}
