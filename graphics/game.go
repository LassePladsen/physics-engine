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

type Game struct {
	Particle physics.Particle2D
}

// Keeps particles in-bounds by bouncing it back, also truncate it back to current window sise
func (g *Game) ensureInBounds() {
	screenWidth, screenHeight := ebiten.ScreenSize()
	width := PixelsToMeters(screenWidth)
	height := PixelsToMeters(screenHeight)
	logger.Debugf("ensureInBounds: maxX, maxY: (%v, %v)", width, height)
	logger.Debugf("ensureInBounds: Particle: %+v", g.Particle)
	if g.Particle.Position.X+g.Particle.Radius > width ||
		g.Particle.Position.X-g.Particle.Radius < 0 {
		g.Particle.Velocity.X *= -1
		g.Particle.Position.X = min(max(g.Particle.Radius, g.Particle.Position.X+g.Particle.Radius), width-g.Particle.Radius)
	}
	if g.Particle.Position.Y+g.Particle.Radius > height ||
		g.Particle.Position.Y-g.Particle.Radius < 0 {
		g.Particle.Velocity.Y *= -1
		g.Particle.Position.Y = min(max(g.Particle.Radius, g.Particle.Position.Y+g.Particle.Radius), height-g.Particle.Radius)
	}
}

func (g *Game) Update() error {
	RunKeybinds()

	if Paused {
		return nil
	}

	g.Particle.Step(DeltaTime)
	g.ensureInBounds()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	_, screenHeight := ebiten.ScreenSize()
	// NB: positive y is down, so reverse the y by subtracting from height
	y := MetersToPixels(PixelsToMeters(screenHeight) - g.Particle.Position.Y)
	x := MetersToPixels(g.Particle.Position.X)
	drawCircle(screen, x, y, MetersToPixels(g.Particle.Radius), color.White)
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
