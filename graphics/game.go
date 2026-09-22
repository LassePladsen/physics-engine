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

// The main state for the rendering game engine
type Game struct {
	World     physics.World
	Camera    Camera
	DeltaTime float64
}

func NewGame() Game {
	return Game{
		World:  physics.NewWorld(),
		Camera: NewCamera(),
	}
}

func (g *Game) DrawAll(screen *ebiten.Image) {
	screenWidth, screenHeight := ebiten.ScreenSize()

	// TODO: make bounds a variable in World and remove DisableBounds on Game maybe
	boundsWidth := g.Camera.ScaleInt(screenWidth)
	boundsHeight := g.Camera.ScaleInt(screenHeight)
	for _, particle := range g.World.Particles {
		// Physics uses positive Y upwards; screen coordinates use positive Y downwards.
		x := MetersToPixels(g.Camera.ScaleFloat(particle.Position.X)) - g.Camera.X
		y := boundsHeight - MetersToPixels(g.Camera.ScaleFloat(particle.Position.Y)) + g.Camera.Y
		DrawCircle(screen, x, y, MetersToPixels(g.Camera.ScaleFloat(particle.Radius)), color.White)
	}
	if !g.World.DisableBounds {
		DrawBounds(screen, g.Camera.X, g.Camera.Y, boundsWidth - g.Camera.X, boundsHeight + g.Camera.Y)
	}
}

func (g *Game) Update() error {
	RunKeybinds()
	if Paused {
		return nil
	}

	g.Step()
	return nil
}

func (g *Game) Step() {
	g.World.Step(g.DeltaTime)
	if !g.World.DisableBounds {
		screenWidth, screenHeight := ebiten.ScreenSize()
		g.World.EnsureParticlesInBoundary(PixelsToMeters(screenWidth), PixelsToMeters(screenHeight))
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawAll(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// Unfilled circle at (x, y)
func DrawCircle(screen *ebiten.Image, x, y int, radius int, clr color.Color) {
	// Draw pixels completely around the x,y center point in a 360 degree = 2pi radians
	for degrees := range 361 {
		radians := float64(degrees) * math.Pi / 180
		ix := float64(radius) * math.Cos(radians)
		iy := float64(radius) * math.Sin(radians)
		screen.Set(x+physics.Round(ix), y+physics.Round(iy), clr)
	}
}

func GetMonitorCenter() (int, int) {
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

// Draws rectangular bounds
func DrawBounds(screen *ebiten.Image, x int, y int, width int, height int) {
	clr := color.RGBA{64, 64, 64, 255}
	// TODO: fix bounds when moving camera coords
	// Top and bottom
	for xi := x; xi < width; xi++ {
		screen.Set(xi, y, clr)
		screen.Set(xi, y+height-1, clr)
	}
	// left and right
	for yi := y; yi < height; yi++ {
		screen.Set(x, yi, clr)
		screen.Set(x+width-1, yi, clr)
	}
}
