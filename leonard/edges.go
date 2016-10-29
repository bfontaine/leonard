package leonard

import (
	"image"
	"image/color"
	"math"
)

func gradients(img image.Image, fn func(image.Image, int, int) float64) *image.Gray16 {
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
			// Storing stuff in a [][]uint16 takes more memory and doesn't save
			// us time here
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

func horizontalGradient(img image.Image, x, y int) float64 {
	return float64(luminance(img.At(x+1, y)) - luminance(img.At(x-1, y)))
}

func verticalGradient(img image.Image, x, y int) float64 {
	return float64(luminance(img.At(x, y+1)) - luminance(img.At(x, y-1)))
}

// Return an image that represents the "intensity" of the horizontal gradients
func HorizontalGradients(img image.Image) *image.Gray16 {
	return gradients(img, horizontalGradient)
}

// Return an image that represents the "intensity" of the vertical gradients
func VerticalGradients(img image.Image) *image.Gray16 {
	return gradients(img, verticalGradient)
}

// Return an image that represents the magnitude of gradients
func Gradients(img image.Image) *image.Gray16 {
	return gradients(img, func(img image.Image, x, y int) float64 {
		// Read e.g. http://www.cse.psu.edu/~rtc12/CSE486/lecture02.pdf
		// Also: https://www.cs.umd.edu/~djacobs/CMSC426/ImageGradients.pdf
		//       https://en.wikipedia.org/wiki/Image_gradient
		h := horizontalGradient(img, x, y)
		v := verticalGradient(img, x, y)
		return math.Sqrt(h*h + v*v)
	})
}

func (b *BinaryImage) thinEdgesIteration(odd bool) (*BinaryImage, bool) {
	b2 := NewEmptyBinaryImage(b.height, b.width)

	changed := false
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			if !b.Get(x, y) {
				continue
			}

			b2.Set(x, y, true)

			// p9 p2 p3
			// p8 P1 p4
			// p7 p6 p5
			p2 := b.Get(north.apply(x, y))
			p3 := b.Get(northeast.apply(x, y))
			p4 := b.Get(east.apply(x, y))
			p5 := b.Get(southeast.apply(x, y))
			p6 := b.Get(south.apply(x, y))
			p7 := b.Get(southwest.apply(x, y))
			p8 := b.Get(west.apply(x, y))
			p9 := b.Get(northwest.apply(x, y))

			// B(P1)
			count := 0
			// A(P1)
			transitions := 0
			lastWas0 := false
			for _, p := range []bool{p2, p3, p4, p5, p6, p7, p8, p9} {
				if !p {
					// 0
					lastWas0 = true
				} else {
					count++

					if lastWas0 {
						// 0->1
						lastWas0 = false
						transitions++
					}
				}
			}

			// (a)
			if count < 2 || count > 6 {
				// preserve the node
				continue
			}

			// last transition
			if !p9 && p2 {
				transitions++
			}

			// (a), (b), (c), (d) below come from the original paper

			switch transitions {
			// (a)
			case 1:
				// (b)
				if odd {
					// first subiteration

					// (c)
					if p2 && p4 && p6 {
						continue
					}

					// (d)
					if p4 && p6 && p8 {
						continue
					}
				} else {
					// second subiteration

					// (c')
					if p2 && p4 && p8 {
						continue
					}

					// (d')
					if p2 && p6 && p8 {
						continue
					}
				}
			case 2:
				// Kocharyan (2013)'s modifications
				if odd {
					if !p4 || !p2 || p3 || p7 || p8 {
						continue
					}
				} else {
					if !p6 || !p8 || p3 || p4 || p7 {
						continue
					}
				}
			}

			// delete the pixel
			b2.Set(x, y, false)
			changed = true

		}
	}
	return b2, changed
}

// Thin the edges of an image [that went through Gradients()] and return it.
// The image is modified in-place.
func (b *BinaryImage) ThinEdges() *BinaryImage {
	// We use Zhang-Suen's algorithm (1984) + modifications from Kocharyan
	// (2013)

	// See:
	// https://dl.acm.org/citation.cfm?id=358023
	// http://article.sciencepublishinggroup.com/pdf/10.11648.j.ajsea.20130201.11.pdf

	// TODO check http://www.uel.br/pessoal/josealexandre/stuff/thinning/ftp/zhang-wang.pdf

	// There's also this algorithm but it's really slow:
	// https://users.fmrib.ox.ac.uk/~steve/susan/thinning/node2.html

	changed1 := true
	changed2 := true

	for changed1 || changed2 {
		b, changed1 = b.thinEdgesIteration(true)
		b, changed2 = b.thinEdgesIteration(false)
	}

	return b
}
