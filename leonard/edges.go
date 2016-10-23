package leonard

import (
	"image"
	"image/color"
	"math"
)

func gradients(img image.Image, fn func(image.Image, int, int) float64) image.Image {
	bounds := img.Bounds()
	grads := image.NewGray16(bounds)

	// Skip the borders
	minX, maxX := bounds.Min.X+1, bounds.Max.X-1
	minY, maxY := bounds.Min.Y+1, bounds.Max.Y-1

	maxGrad := 0.0

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			g := math.Abs(fn(img, x, y))
			if g > maxGrad {
				maxGrad = g
			}
			grads.SetGray16(x, y, color.Gray16{uint16(g)})
		}
	}

	// adjust based on the max value
	maxGrad /= 0xFFFF

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			g := grads.Gray16At(x, y)
			g.Y = uint16(float64(g.Y) / maxGrad)
			grads.SetGray16(x, y, g)
		}
	}

	return grads
}

// Return an image that represents the intensity of the horizontal gradients
func HorizontalGradients(img image.Image) image.Image {
	return gradients(img, func(img image.Image, x, y int) float64 {
		return float64(luminance(img.At(x+1, y)) - luminance(img.At(x-1, y)))
	})
}

// Return an image that represents the intensity of the vertical gradients
func VerticalGradients(img image.Image) image.Image {
	return gradients(img, func(img image.Image, x, y int) float64 {
		return float64(luminance(img.At(x, y+1)) - luminance(img.At(x, y-1)))
	})
}