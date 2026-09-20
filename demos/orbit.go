package main

import (
	"math"
	"os"

	"github.com/LassePladsen/physics-engine/graphics"
	"github.com/LassePladsen/physics-engine/logger"
	"github.com/LassePladsen/physics-engine/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

const G = 6.6743e-11

func runOrbitDemo() {
	monitorWidth, monitorHeight := ebiten.Monitor().Size()
	windowWidth := monitorWidth * 2 / 3
	windowHeight := monitorHeight * 2 / 3 

	// Particles
	radius := 0.2 // m
	p1 := physics.Particle2D{
		Mass:     1e11, // kg
		Radius:   radius,
		Position: physics.Vec2{X: graphics.PixelsToMeters(physics.Round(float64(windowWidth) * 1.2 / 3.0)), Y: graphics.PixelsToMeters(windowHeight / 2)},
		GravityEnabled: true,
	}

	p2 := p1
	p2.Position.X = graphics.PixelsToMeters(windowWidth) - p1.Position.X

	// Fix stable orbit
	separation := p1.DistanceTo(p2)
	p1.Velocity = physics.Vec2{X: 0, Y: math.Sqrt(G*p1.Mass/(2*separation))}
	p2.Velocity = p1.Velocity.Mul(-1)


	// Init the simulation
	game := graphics.Game{
		DeltaTime: deltaTime,
		World: physics.World{
			DisableBounds: true,
			Particles:               []physics.Particle2D{p1, p2},
			ParticleToParticleGravitationalStrength: G,
		},
	}

	graphics.InitKeybinds(game)

	if err := ebiten.RunGame(&game); err != nil {
		if err.Error() != "" {
			logger.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
}
