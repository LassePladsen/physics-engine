package graphics

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func InitGraphics(ticksPerSecond int) {
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
}

func InitKeybinds(game Game) {
	AddKeybind(TogglePause, true, ebiten.KeyP) // pause

	// exit
	exit := func() {
		os.Exit(0)
	}
	AddKeybind(exit, false, ebiten.KeyEscape)
	AddKeybind(exit, false, ebiten.KeyQ)
	AddKeybind(exit, false, ebiten.KeyC, ebiten.KeyControl)

	// Step once
	AddKeybind(func() { game.Step() }, true, ebiten.KeyS)
}
