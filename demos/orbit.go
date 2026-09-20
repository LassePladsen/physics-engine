package main

import (
	"os"

	"github.com/LassePladsen/physics-engine/graphics"
	"github.com/LassePladsen/physics-engine/logger"
	"github.com/LassePladsen/physics-engine/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

func runOrbitDemo() {
	monitorWidth, monitorHeight := ebiten.Monitor().Size()
	windowWidth := monitorWidth * 2 / 3
	windowHeight := monitorHeight * 2 / 3 

	// Particles
	radius := 0.2 // m
	p1 := physics.Particle2D{
		Mass:     3e5, // kg
		Radius:   radius,
		Position: physics.Vec2{X: graphics.PixelsToMeters(physics.Round(float64(windowWidth) * 1.2 / 3.0)), Y: graphics.PixelsToMeters(windowHeight / 2)},
		Velocity: physics.Vec2{X: 0.1, Y: 0.8},
		GravityEnabled: true,
	}
	p2 := p1
	p2.Position.X = graphics.PixelsToMeters(windowWidth) - p1.Position.X
	p2.Velocity = p1.Velocity.Mul(-1)
	// p2.Mass *= 1.8
	// p2.Radius *= 1.8

	// Init the simulation
	game := graphics.Game{
		DeltaTime: deltaTime,
		World: physics.World{
			DisableBounds: true,
			Particles:               []physics.Particle2D{p1, p2},
			ParticleToParticleGravitationalStrength: 6.6743e-11,
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
