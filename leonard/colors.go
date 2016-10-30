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
	return n * 0xFF / 0xFFFF
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

// BinaryImage is a black & white image represented as a boolean matrix.
//
// White represents true pixels and black represents false ones. The underlying
// representation is optimized for sparse matrices; i.e. with a lot more false
// values than true ones.
type BinaryImage struct {
	height, width int
	pixels        map[image.Point]bool
}

var _ image.Image = &BinaryImage{}

// ColorModel implements the image.Image interface
func (b *BinaryImage) ColorModel() color.Model {
	return color.GrayModel
}

// Bounds implements the image.Image interface
func (b *BinaryImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, b.width, b.height)
}

// At implements the image.Image interface
func (b *BinaryImage) At(x, y int) color.Color {
	if b.Get(x, y) {
		return white
	}
	return black
}

// Set sets the value at a given pixel
func (b *BinaryImage) Set(x, y int, value bool) {
	p := image.Point{x, y}

	if !value {
		delete(b.pixels, p)
	}

	b.pixels[p] = value
}

// Get returns the boolean value at a given pixel
func (b *BinaryImage) Get(x, y int) bool {
	return b.pixels[image.Point{x, y}]
}

// Invert inverts the image.
//
// Black pixels become white and white ones become black.
func (b *BinaryImage) Invert() {
	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			b.Set(x, y, !b.Get(x, y))
		}
	}
}

// NewEmptyBinaryImage returns a new empty (= all black) binary image
func NewEmptyBinaryImage(height, width int) *BinaryImage {
	return &BinaryImage{
		height: height,
		width:  width,
		pixels: make(map[image.Point]bool),
	}
}

// DefaultBinaryThreshold is the default threshold used for binary images.
const DefaultBinaryThreshold = 0x28F6 // 0.16 * 0xFFFF

// NewBinaryImage creates a new binary image from a given one. The default
// threshold is used is the passed value is -1.
func NewBinaryImage(img image.Image, threshold int) *BinaryImage {
	bounds := img.Bounds()

	if threshold == -1 {
		threshold = DefaultBinaryThreshold
	}

	b := NewEmptyBinaryImage(bounds.Max.Y, bounds.Max.X)

	t := float32(threshold)

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			p := luminance(img.At(x, y))
			if p >= t {
				b.Set(x, y, true)
			}
		}
	}

	return b
}
