package graphics

import (
	"testing"

	"github.com/LassePladsen/physics-engine/physics"
)

func TestUnitConversions(t *testing.T) {
	pixelTests := []struct {
		pixels int
		want   float64
	}{{0, 0}, {1, 0.01}, {100, 1}, {250, 2.5}, {-50, -0.5}}
	for _, tt := range pixelTests {
		if got := PixelsToMeters(tt.pixels); got != tt.want {
			t.Errorf("PixelsToMeters(%d) = %v, want %v", tt.pixels, got, tt.want)
		}
	}

	meterTests := []struct {
		meters float64
		want   int
	}{{0, 0}, {0.0149, 1}, {0.015, 2}, {1, 100}, {-0.015, -2}}
	for _, tt := range meterTests {
		if got := MetersToPixels(tt.meters); got != tt.want {
			t.Errorf("MetersToPixels(%v) = %d, want %d", tt.meters, got, tt.want)
		}
	}
}

func TestGameLayoutPreservesOutsideDimensions(t *testing.T) {
	tests := []struct{ width, height int }{{0, 0}, {640, 480}, {1920, 1080}}
	game := Game{}
	for _, tt := range tests {
		width, height := game.Layout(tt.width, tt.height)
		if width != tt.width || height != tt.height {
			t.Errorf("Layout(%d, %d) = (%d, %d), want unchanged dimensions", tt.width, tt.height, width, height)
		}
	}
}

func TestGameStepAllUpdatesParticlesInPlace(t *testing.T) {
	previousDeltaTime := DeltaTime
	DeltaTime = 0.5
	defer func() { DeltaTime = previousDeltaTime }()

	game := Game{Particles: []physics.Particle2D{{
		Position: physics.Vec2{X: 1, Y: 2},
		Velocity: physics.Vec2{X: 4, Y: -2},
	}}}

	game.StepAll()

	want := physics.Vec2{X: 3, Y: 1}
	if got := game.Particles[0].Position; got != want {
		t.Errorf("particle position after StepAll = %v, want %v", got, want)
	}
}
