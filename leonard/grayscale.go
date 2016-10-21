package leonard

import (
	"image"
	"image/color"
	"math"
)

var (
	black = color.Gray16{0}
	white = color.Gray16{0xffff}
)

func luminance(r, g, b uint32) float32 {
	// Ref:
	// https://en.wikipedia.org/wiki/Grayscale#Converting_color_to_grayscale
	return 0.2126*float32(r) + 0.7152*float32(g) + 0.0722*float32(b)
}

func grayscale(r, g, b, a uint32) uint16 {
	alpha := float32(a) / 0xffff
	linear := luminance(r, g, b) * alpha

	return uint16(linear)
}

// Convert a colored image to a grayscaled one
func Grayscale(img image.Image) image.Image {
	grayscaled := image.NewNRGBA(img.Bounds())

	bd := img.Bounds()
	for y := bd.Min.Y; y < bd.Max.Y; y++ {
		for x := bd.Min.X; x < bd.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			grayscaled.Set(x, y, color.Gray16{grayscale(r, g, b, a)})
		}
	}

	return grayscaled
}

func Binary(img image.Image) image.Image {
	return BinaryWithThreshold(img, uint16(math.Ceil(0xffff*0.9)))
}

func BinaryWithThreshold(img image.Image, threshold uint16) image.Image {
	binary := image.NewNRGBA(img.Bounds())

	bd := img.Bounds()
	for y := bd.Min.Y; y < bd.Max.Y; y++ {
		for x := bd.Min.X; x < bd.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			gray := grayscale(r, g, b, a)

			if gray >= threshold {
				binary.Set(x, y, black)
			} else {
				binary.Set(x, y, white)
			}
		}
	}

	return binary
}
