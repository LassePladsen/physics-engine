package physics

import (
	"math"
	"math/rand/v2"
)

// Round to nearest int
func Round(x float64) int {
	return int(math.Round(x))
}

// Clamp value to [minVal, maxVal]
func Clamp(val, minVal, maxVal float64) float64 {
	return max(min(val, maxVal), minVal)
}

func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
