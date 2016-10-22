package leonard

import (
	"image/color"
	"math"
)

var (
	black = color.Gray16{0}
	white = color.Gray16{0xffff}
)

func absint(n int) int {
	return int(math.Abs(float64(n)))
}

func clip(n, min, max float64) float64 {
	return math.Min(math.Max(n, min), max)
}
