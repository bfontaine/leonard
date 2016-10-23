package leonard

import "math"

func absint(n int) int {
	return int(math.Abs(float64(n)))
}

func clip(n, min, max float64) float64 {
	return math.Min(math.Max(n, min), max)
}
