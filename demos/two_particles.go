package main

import (
	"os"

	"github.com/LassePladsen/physics-engine/graphics"
	"github.com/LassePladsen/physics-engine/logger"
	"github.com/LassePladsen/physics-engine/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

func runTwoParticlesDemo() {
	monitorWidth, monitorHeight := ebiten.Monitor().Size()
	windowWidth := monitorWidth * 2 / 3
	windowHeight := monitorHeight * 2 / 3 

	// Init the simulation
	game := graphics.NewGame()
	game.DeltaTime = deltaTime
	game.World = physics.NewWorld()
	game.World.Gravity = physics.UnitDown().Mul(gravityAcceleration)
	game.World.Friction = friction
	game.World.WindowBoundsRestitution = windowBoundsRestitution
	game.World.ParticleRestitution = particleRestitution

	// Particles
	radius := 0.2 // m
	p1 := physics.Particle2D{
		Mass:     1,
		Radius:   radius,
		Position: physics.Vec2{X: radius * 1.5, Y: graphics.PixelsToMeters(windowHeight * 2 / 3 )},
		Velocity: physics.Vec2{X: 5, Y: 0},
	}
	p2 := p1
	p2.Position.X = graphics.PixelsToMeters(windowWidth) - p1.Position.X
	p2.Velocity.X = -p1.Velocity.X
	p2.Mass *= 1.8
	p2.Radius *= 1.8
	game.World.Particles = []physics.Particle2D{p1, p2}
	game.Camera.Zoom /= 2

	graphics.InitKeybinds(&game)

	if err := ebiten.RunGame(&game); err != nil {
		if err.Error() != "" {
			logger.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
}
