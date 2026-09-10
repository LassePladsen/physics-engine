package main

import (
	"errors"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func getMonitorCenter() (int, int) {
	w, h := ebiten.Monitor().Size()
	return w / 2, h / 2
}

type Game struct{}

func (g *Game) Update() error {
	// Exit on esc, q, and ctrl-c
	if ebiten.IsKeyPressed(ebiten.KeyEscape) ||
		ebiten.IsKeyPressed(ebiten.KeyQ) ||
		(ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyC)) {
		return errors.New("")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	bounds := screen.Bounds().Size()
	drawCircle(screen, bounds.X / 2, bounds.Y / 2, 15, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// Unfilled circle at (x, y)
func drawCircle(screen *ebiten.Image, x, y int, radius int, clr color.Color) {
	// Draw pixels completely around the x,y center point in a 360 degree = 2pi radians
	for degrees := range(361) {
		radians := float64(degrees) * math.Pi / 180
		ix := float64(radius) * math.Cos(radians)
		iy := float64(radius) * math.Sin(radians)
		screen.Set(x + int(math.Round(ix)), y + int(math.Round(iy)), clr)
	}
}

func main() {
	mw, mh := ebiten.Monitor().Size()
	w := mw * 2/3
	h := mh * 2/3
	x := (mw - w) / 2
	y := (mh - h) / 2
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowPosition(x, y)
	ebiten.SetWindowTitle("Physics Engine")
	if err := ebiten.RunGame(&Game{}); err != nil {
		if err.Error() != "" {
			log.Fatal(err)
		}
	}
}
