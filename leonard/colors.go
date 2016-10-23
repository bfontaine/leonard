package leonard

import "image/color"

var (
	black = color.Gray16{0}
	white = color.Gray16{0xffff}
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
