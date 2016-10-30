package leonard

import (
	"image"
	"image/color"
	"math"
)

func averageColor(colors []color.Color) color.Color {
	var r, g, b, a float64 = 0, 0, 0, 0

	for _, c := range colors {
		rc, gc, bc, ac := c.RGBA()
		r += float64(rc)
		g += float64(gc)
		b += float64(bc)
		a += float64(ac)
	}

	n := float64(len(colors))

	return color.NRGBA{
		uint8(to255(r / n)),
		uint8(to255(g / n)),
		uint8(to255(b / n)),
		uint8(to255(a / n)),
	}
}

// Downscale reduces the size of an image by 4x (width/2 and height/2) by
// averaging the values of 4-pixels squares.
func Downscale(img image.Image) image.Image {
	// https://en.wikipedia.org/wiki/Pyramid_(image_processing)

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	width2 := int(math.Ceil(float64(width) / 2.0))
	height2 := int(math.Ceil(float64(height) / 2.0))

	downscaled := image.NewRGBA(image.Rect(0, 0, width2, height2))

	for y := 0; y < height2; y++ {
		for x := 0; x < width2; x++ {
			x1 := x * 2
			y1 := y * 2

			downscaled.Set(x, y,
				averageColor([]color.Color{
					img.At(x1, y1),
					img.At(x1+1, y1),
					img.At(x1, y1+1),
					img.At(x1+1, y1+1)}))
		}
	}

	return downscaled
}
