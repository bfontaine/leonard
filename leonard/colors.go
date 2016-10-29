package leonard

import (
	"image"
	"image/color"
)

var (
	black = color.Gray{0}
	white = color.Gray{0xff}
)

// Convert a value in a 0-0xFFFF range into one in a 0-0xFF one
func to255(n float64) float64 {
	return clip(n*0xFF/0xFFFF, 0, 0xFFFF)
}

func luminanceRGB(r, g, b uint32) float32 {
	// Ref:
	// https://en.wikipedia.org/wiki/Grayscale#Converting_color_to_grayscale
	return 0.2126*float32(r) + 0.7152*float32(g) + 0.0722*float32(b)
}

func luminance(c color.Color) float32 {
	r, g, b, _ := c.RGBA()
	return luminanceRGB(r, g, b)
}

type BinaryImage struct {
	height, width int
	pixels        []bool
}

func (b *BinaryImage) ColorModel() color.Model {
	return color.GrayModel
}

func (b *BinaryImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, b.width, b.height)
}

func (b *BinaryImage) At(x, y int) color.Color {
	if b.Get(x, y) {
		return white
	}
	return black
}

func (b *BinaryImage) Set(x, y int, value bool) {
	idx := y*b.width + x
	if idx < 0 || idx >= len(b.pixels) {
		return
	}

	b.pixels[idx] = value
}

func (b *BinaryImage) Get(x, y int) bool {
	idx := y*b.width + x
	if idx < 0 || idx >= len(b.pixels) {
		return false
	}
	return b.pixels[idx]
}

func NewEmptyBinaryImage(height, width int) *BinaryImage {
	return &BinaryImage{
		height: height,
		width:  width,
		pixels: make([]bool, height*width),
	}
}

func NewBinaryImage(img image.Image, threshold int) *BinaryImage {
	bounds := img.Bounds()

	b := NewEmptyBinaryImage(bounds.Max.Y, bounds.Max.X)

	t := float32(threshold)

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			p := luminance(img.At(x, y))
			if p >= t {
				b.pixels[y*b.width+x] = true
			}
		}
	}

	return b
}

var _ image.Image = &BinaryImage{}
