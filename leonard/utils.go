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

// Convert a value in a 0-0xFFFF range into one in a 0-0xFF one
func to255(n float64) float64 {
	return clip(n*0xFF/0xFFFF, 0, 0xFFFF)
}
