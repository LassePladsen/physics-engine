package graphics

import "github.com/LassePladsen/physics-engine/physics"

// Represents the PoV
type Camera struct {
	X, Y int // pixels. positive y is downwards
	Zoom float64
}

func NewCamera() Camera {
	return Camera{
		Zoom: 1,
	}
}

func (c Camera) ScaleInt(val int) int {
	return physics.Round(float64(val) * c.Zoom)
}

func (c Camera) ScaleFloat(val float64) float64 {
	return val * c.Zoom
}
