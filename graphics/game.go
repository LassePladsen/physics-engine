// The main game state logic for the rendering engine

// The graphics rendering logic
package graphics

import (
	"image/color"
	"math"

	"github.com/LassePladsen/physics-engine/logger"
	"github.com/LassePladsen/physics-engine/physics"
	"github.com/hajimehoshi/ebiten/v2"
)

const pixelsPerMeter = 100

var Paused = false
var DeltaTime = 0.01 // seconds

// The main state for the rendering game engine
type Game struct {
	World physics.World
}

func (g *Game) DrawAll(screen *ebiten.Image) {
	_, screenHeight := ebiten.ScreenSize()
	for _, particle := range g.World.Particles {
		// NB: positive y is down, so reverse the y by subtracting from height
		y := MetersToPixels(PixelsToMeters(screenHeight) - particle.Position.Y)
		x := MetersToPixels(particle.Position.X)
		drawCircle(screen, x, y, MetersToPixels(particle.Radius), color.White)
	}
}

func (g *Game) Update() error {
	RunKeybinds()

	if Paused {
		return nil
	}

	g.World.Step(DeltaTime)
	screenWidth, screenHeight := ebiten.ScreenSize()
	g.World.EnsureParticlesInBounds(PixelsToMeters(screenWidth), PixelsToMeters(screenHeight))
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawAll(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// Unfilled circle at (x, y)
func drawCircle(screen *ebiten.Image, x, y int, radius int, clr color.Color) {
	// Draw pixels completely around the x,y center point in a 360 degree = 2pi radians
	for degrees := range 361 {
		radians := float64(degrees) * math.Pi / 180
		ix := float64(radius) * math.Cos(radians)
		iy := float64(radius) * math.Sin(radians)
		screen.Set(x+physics.Round(ix), y+physics.Round(iy), clr)
	}
}

func getMonitorCenter() (int, int) {
	w, h := ebiten.Monitor().Size()
	return w / 2, h / 2
}

func PixelsToMeters(pixels int) float64 {
	return float64(pixels) / pixelsPerMeter
}

// Rounds to nearest int
func MetersToPixels(meters float64) int {
	return physics.Round(meters * pixelsPerMeter)
}

func TogglePause() {
	Paused = !Paused
	msg := "Simulation unpaused"
	if Paused {
		msg = "Simulation paused"
	}
	logger.Info(msg)
}
