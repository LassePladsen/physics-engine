package main

import (
	"os"

	"github.com/LassePladsen/physics-engine/graphics"
	"github.com/LassePladsen/physics-engine/logger"
	"github.com/LassePladsen/physics-engine/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

const fps = 60.0
const gravityAcceleration = 9.81 // m/s^2

func main() {
	logger.Init()

	monitorWidth, monitorHeight := ebiten.Monitor().Size()
	windowWidth := monitorWidth * 2 / 3
	windowHeight := monitorHeight * 2 / 3
	windowX := (monitorWidth - windowWidth) / 2
	windowY := (monitorHeight - windowHeight) / 2

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowPosition(windowX, windowY)
	ebiten.SetWindowTitle("Physics Engine")
	ebiten.SetTPS(fps)

	radius := 0.2 // m
	graphics.DeltaTime = 1 / fps

	// Keybinds
	graphics.AddKeybind(graphics.TogglePause, true, ebiten.KeyP)
	exit := func() {
		os.Exit(0)
	}
	graphics.AddKeybind(exit, false, ebiten.KeyEscape)
	graphics.AddKeybind(exit, false, ebiten.KeyQ)
	graphics.AddKeybind(exit, false, ebiten.KeyC, ebiten.KeyControl)

	p1 := physics.Particle2D{
				Position:     physics.Vec2{X: radius, Y: graphics.PixelsToMeters(monitorHeight / 2)},
				Velocity:     physics.Vec2{X: 5, Y: 0},
				Acceleration: physics.Vec2Down().Mul(gravityAcceleration), // Constant acc as of now. NB: downwards is positive y
				Radius:       radius,
			}
	p2 := p1
	p2.Position.X = float64(windowWidth)-p1.Position.X
	p2.Velocity.X = -p1.Velocity.X
	game := graphics.Game{Particles: []physics.Particle2D{p1, p2}}
	if err := ebiten.RunGame(&game); err != nil {
		if err.Error() != "" {
			logger.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
}
