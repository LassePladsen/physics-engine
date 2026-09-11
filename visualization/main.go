package main

import (
	"errors"
	"image/color"
	"math"
	"os"

	"github.com/LassePladsen/physics-engine/logger"
	"github.com/LassePladsen/physics-engine/physics"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const fps = 60.0
const dt = 1 / fps
const gravityAcceleration = 9.81 // m/s^2
const pixelsPerMeter = 100

var Paused = false

type Game struct {
	Particle physics.Particle2D
}

// Keeps particles in-bounds by bouncing it back, also truncate it back to current window sise
func (g *Game) ensureInBounds() {
	screenWidth, screenHeight := ebiten.ScreenSize()
	width := pixelsToMeters(screenWidth)
	height := pixelsToMeters(screenHeight)
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
	// Exit on esc, q, and ctrl-c
	if ebiten.IsKeyPressed(ebiten.KeyEscape) ||
		ebiten.IsKeyPressed(ebiten.KeyQ) ||
		(ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyC)) {
		return errors.New("")
	}

	// Toggle pause on p
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		Paused = !Paused
		msg := "Simulation unpaused"
		if Paused {
			msg = "Simulation paused"
		}
		logger.Info(msg)
	}
	if Paused {
		return nil
	}

	g.Particle.Step(dt)
	g.ensureInBounds()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	_, screenHeight := ebiten.ScreenSize()
	// NB: positive y is down, so reverse the y by subtracting from height
	y := metersToPixels(pixelsToMeters(screenHeight) - g.Particle.Position.Y)
	x := metersToPixels(g.Particle.Position.X)
	drawCircle(screen, x, y, metersToPixels(g.Particle.Radius), color.White)
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
		screen.Set(x+round(ix), y+round(iy), clr)
	}
}

func getMonitorCenter() (int, int) {
	w, h := ebiten.Monitor().Size()
	return w / 2, h / 2
}

func round(x float64) int {
	return int(math.Round(x))
}

func pixelsToMeters(pixels int) float64 {
	return float64(pixels) / pixelsPerMeter
}

// Rounds to nearest int
func metersToPixels(meters float64) int {
	return round(meters * pixelsPerMeter)
}

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

	game := Game{physics.Particle2D{
		Position:     physics.Vec2{X: radius, Y: pixelsToMeters(mh / 2)},
		Velocity:     physics.Vec2{X: 1, Y: 0},
		Acceleration: physics.Vec2Down().Mul(gravityAcceleration), // Constant acc as of now. NB: downwards is positive y
		Radius:       radius,
	}}
	if err := ebiten.RunGame(&game); err != nil {
		if err.Error() != "" {
			logger.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
}
