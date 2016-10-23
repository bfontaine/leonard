package leonard

import (
	"image"
	"math"
)

func absint(n int) int {
	return int(math.Abs(float64(n)))
}

func abs32(n float32) float64 { return math.Abs(float64(n)) }

func clip(n, min, max float64) float64 {
	return math.Min(math.Max(n, min), max)
}

func toBoolMatrix(img *image.Gray16) *boolMatrix {
	bounds := img.Bounds()

	height := bounds.Max.Y - bounds.Min.Y
	width := bounds.Max.X - bounds.Min.X

	m := newBoolMatrix(height, width)

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			m.set(i, j, img.Gray16At(bounds.Min.X+i, bounds.Min.Y+j).Y > 0)
		}
	}

	return m
}
