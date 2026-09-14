package physics

import (
	"testing"
)

func TestClamp(t *testing.T) {
	tests := []struct {
		name            string
		value, min, max float64
		want            float64
	}{
		{"below minimum", -3, -2, 5, -2},
		{"at minimum", -2, -2, 5, -2},
		{"within range", 1.5, -2, 5, 1.5},
		{"at maximum", 5, -2, 5, 5},
		{"above maximum", 8, -2, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clamp(tt.value, tt.min, tt.max); got != tt.want {
				t.Fatalf("Clamp(%v, %v, %v) = %v, want %v", tt.value, tt.min, tt.max, got, tt.want)
			}
		})
	}
}
