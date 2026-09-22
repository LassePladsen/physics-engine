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

func InitKeybinds(game *Game) {
	RegisterKeybind(TogglePause, true, ebiten.KeyP) // pause

	// exit
	exit := func() {
		os.Exit(0)
	}
	RegisterKeybind(exit, false, ebiten.KeyEscape)
	RegisterKeybind(exit, false, ebiten.KeyQ)
	RegisterKeybind(exit, false, ebiten.KeyC, ebiten.KeyControl)

	// Camera actions
	// TODO: add on mouse buttons, need to extend keybinds with mouse binds
	zoom := 1.01
	RegisterKeybind(func() {game.Camera.Zoom *= zoom}, false, ebiten.KeyMinus) // KeyMinus on NO keyboards is plus
	RegisterKeybind(func() {game.Camera.Zoom /= zoom}, false, ebiten.KeySlash) // KeySlash on NO keyboards is minus
	move := 20
	RegisterKeybind(func() {game.Camera.X -= move}, false, ebiten.KeyLeft)
	RegisterKeybind(func() {game.Camera.X += move}, false, ebiten.KeyRight)
	RegisterKeybind(func() {game.Camera.Y += move}, false, ebiten.KeyUp)
	RegisterKeybind(func() {game.Camera.Y -= move}, false, ebiten.KeyDown)

	// Step once
	RegisterKeybind(func() { game.Step() }, true, ebiten.KeyS)
}
