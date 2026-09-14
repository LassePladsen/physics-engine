package main

import (
	"os"

	"github.com/LassePladsen/physics-engine/graphics"
	"github.com/LassePladsen/physics-engine/logger"
	"github.com/LassePladsen/physics-engine/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

const ticksPerSecond = 60.0
const deltaTime = 1 / ticksPerSecond
const gravityAcceleration = 9.81 // m/s^2
const particleRestitution = 0.2
const windowBoundsRestitution = 0.5
const friction = 0.1

func main() {
	logger.Init()

	// Graphics config
	monitorWidth, monitorHeight := ebiten.Monitor().Size()
	windowWidth := monitorWidth * 2 / 3
	windowHeight := monitorHeight * 2 / 3
	windowX := (monitorWidth - windowWidth) / 2
	windowY := (monitorHeight - windowHeight) / 2
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowPosition(windowX, windowY)
	ebiten.SetWindowTitle("Physics Engine")
	ebiten.SetTPS(ticksPerSecond)

	// Init the simulation
	world := physics.GenerateRandom(physics.WorldGenConfig{
		MinMass:       .1,
		MaxMass:       10,
		NumParticles:  100,
		MinRadius:     .06,
		MaxRadius:     .35,
		MinPositions:  physics.Vec2{X: 0.1, Y: 2},
		MaxPositions:  physics.Vec2{X: graphics.PixelsToMeters(windowWidth), Y: graphics.PixelsToMeters(windowHeight)},
		MinVelocities: physics.Vec2{X: .2, Y: .2},
		MaxVelocities: physics.Vec2{X: 50, Y: 50},
	})
	world.Gravity = physics.UnitDown().Mul(gravityAcceleration)
	world.Friction = friction
	world.WindowBoundsRestitution = windowBoundsRestitution
	// world.ParticleRestitution = particleRestitution
	world.DisableCollisions = true

	game := graphics.Game{
		DeltaTime: deltaTime,
		World:     world,
	}

	/// KEYBINDS
	graphics.AddKeybind(graphics.TogglePause, true, ebiten.KeyP) // pause

	// exit
	exit := func() {
		os.Exit(0)
	}
	graphics.AddKeybind(exit, false, ebiten.KeyEscape)
	graphics.AddKeybind(exit, false, ebiten.KeyQ)
	graphics.AddKeybind(exit, false, ebiten.KeyC, ebiten.KeyControl)

	// Step once
	graphics.AddKeybind(func() { game.Step() }, true, ebiten.KeyS)
	/// END KEYBINDS

	if err := ebiten.RunGame(&game); err != nil {
		if err.Error() != "" {
			logger.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
}
