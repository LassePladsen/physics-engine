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

	mw, mh := ebiten.Monitor().Size()
	w := mw * 2 / 3
	h := mh * 2 / 3
	x := (mw - w) / 2
	y := (mh - h) / 2

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowPosition(x, y)
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

	game := graphics.Game{Particles: []physics.Particle2D{{
		Position:     physics.Vec2{X: radius, Y: graphics.PixelsToMeters(mh / 2)},
		Velocity:     physics.Vec2{X: 5, Y: 0},
		Acceleration: physics.Vec2Down().Mul(gravityAcceleration), // Constant acc as of now. NB: downwards is positive y
		Radius:       radius,
	}}}
	if err := ebiten.RunGame(&game); err != nil {
		if err.Error() != "" {
			logger.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
}
