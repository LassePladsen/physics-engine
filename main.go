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
const windowBoundsRestitution = 0.55
const friction = 0.8

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

	// Particles
	radius := 0.2 // m
	p1 := physics.Particle2D{
		Mass:         1,
		Radius:       radius,
		Position:     physics.Vec2{X: radius * 1.5, Y: graphics.PixelsToMeters(monitorHeight / 2)},
		Velocity:     physics.Vec2{X: 5, Y: 0},
		Acceleration: physics.Vec2Down().Mul(gravityAcceleration), // Constant downward acceleration
	}
	p2 := p1
	p2.Position.X = graphics.PixelsToMeters(windowWidth) - p1.Position.X
	p2.Velocity.X = -p1.Velocity.X
	p2.Mass *= 1.8
	p2.Radius *= 1.8

	// Init the simulation
	game := graphics.Game{
		DeltaTime: 1 / ticksPerSecond,
		World: physics.World{
			Friction:                friction,
			WindowBoundsRestitution: windowBoundsRestitution,
			ParticleRestitution:     particleRestitution,
			Particles:               []physics.Particle2D{p1, p2},
		},
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
