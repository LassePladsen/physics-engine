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
	Particles []physics.Particle2D
}

// In-place ensures particle is in-bounds by bouncing it back if outside, also keep within the current window size
func ensureInBounds(particle *physics.Particle2D) {
	screenWidth, screenHeight := ebiten.ScreenSize()
	width := PixelsToMeters(screenWidth)
	height := PixelsToMeters(screenHeight)
	logger.Debugf("ensureInBounds: maxX, maxY: (%v, %v)", width, height)
	logger.Debugf("ensureInBounds: Particle: %+v", particle)
	if particle.Position.X+particle.Radius > width ||
		particle.Position.X-particle.Radius < 0 {
		particle.Velocity.X *= -1
		particle.Position.X = min(max(particle.Radius, particle.Position.X+particle.Radius), width-particle.Radius)
	}
	if particle.Position.Y+particle.Radius > height ||
		particle.Position.Y-particle.Radius < 0 {
		particle.Velocity.Y *= -1
		particle.Position.Y = min(max(particle.Radius, particle.Position.Y+particle.Radius), height-particle.Radius)
	}
}

func (g *Game) EnsureAllInBounds() {
	for i := range g.Particles {
		ensureInBounds(&g.Particles[i])
	}
}

func (g *Game) DrawAll(screen *ebiten.Image) {
	_, screenHeight := ebiten.ScreenSize()
	for _, particle := range g.Particles {
		// NB: positive y is down, so reverse the y by subtracting from height
		y := MetersToPixels(PixelsToMeters(screenHeight) - particle.Position.Y)
		x := MetersToPixels(particle.Position.X)
		drawCircle(screen, x, y, MetersToPixels(particle.Radius), color.White)
	}
}

func (g *Game) StepAll() {
	for i := range g.Particles {
		g.Particles[i].Step(DeltaTime)
	}
}

func (g *Game) Update() error {
	RunKeybinds()

	if Paused {
		return nil
	}

	for i := range g.Particles {
		g.Particles[i].Step(DeltaTime)
		ensureInBounds(&g.Particles[i])
	}
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
