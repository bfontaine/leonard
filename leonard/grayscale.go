package leonard

import (
	"image"
	"image/color"
	"math"
)

func grayscale(r, g, b, a uint32) uint16 {
	alpha := float32(a) / 0xffff
	linear := luminanceRGB(r, g, b) * alpha

	return uint16(linear)
}

// Convert a colored image to a grayscaled one
func Grayscale(img image.Image) image.Image {
	grayscaled := image.NewGray16(img.Bounds())

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
	// TODO Implement Othsu's method to get the correct theshold:
	// https://en.wikipedia.org/wiki/Otsu%27s_method
	// http://ijetch.org/papers/260-T754.pdf
	return NewBinaryImage(img, int(math.Ceil(0xffff*0.6)))
}
