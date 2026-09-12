package physics

import "github.com/LassePladsen/physics-engine/logger"

// World contains the particles that make up a simulation.
type World struct {
	Particles []Particle2D
}

// Advances every particle in the world by deltaTime seconds.
func (w *World) Step(deltaTime float64) {
	for i := range w.Particles {
		w.Particles[i].Step(deltaTime)
		w.DoCollisions(true) // change to inelastic here
	}
}

func (w *World) DoCollisions(elastic bool) {
	for i := range w.Particles {
		u := w.Particles[i]
		for j := range w.Particles {
			v := w.Particles[j]
			if u == v {
				continue
			}
			logger.Debugf("Checking collision for %+v and %+v", u, v)
			if elastic && u.ShouldCollide(v) {
				logger.Debug("They should elastically collide")
				w.Particles[i], w.Particles[j] = u.ElasticCollision(v)
			} else if u.ShouldCollide(v) {
				logger.Debug("They should INelastically collide")
				// TODO: w.Particles[i], w.Particles[j] = u.InelasticCollision(v)
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
