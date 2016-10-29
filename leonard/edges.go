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

func thinEdgesIteration(m *boolMatrix, odd bool) bool {
	changed := false
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if !m.get(x, y) {
				continue
			}

			p2 := m.get(north.apply(x, y))
			p3 := m.get(northeast.apply(x, y))
			p4 := m.get(east.apply(x, y))
			p5 := m.get(southeast.apply(x, y))
			p6 := m.get(south.apply(x, y))
			p7 := m.get(southwest.apply(x, y))
			p8 := m.get(west.apply(x, y))
			p9 := m.get(northwest.apply(x, y))

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

			if count < 2 || count > 6 {
				continue
			}

			// last transition
			if !p9 && p2 {
				transitions++
			}

			switch transitions {
			case 1:
				if (odd && (!(p2 && p4 && p6) || !(p4 && p6 && p8))) ||
					(!odd && (!(p2 && p4 && p8) || !(p2 && p6 && p8))) {
					m.set(x, y, false)
					changed = true
					continue
				}
			case 2:
				if (odd && ((p4 && p6 && !p9) || (p4 && p2 && !p3 && !p7 && !p8))) ||
					(!odd && ((p2 && p8 && !p5) || (p6 && p8 && !p3 && !p4 && !p7))) {
					m.set(x, y, false)
					changed = true
					continue
				}
			}
		}
	}
	return changed
}

// Thin the edges of an image that went through Gradients() and return it.
// The image is modified in-place.
func ThinEdges(img *image.Gray16) *image.Gray16 {
	m := toBoolMatrix(img)

	// There are a bunch of algorithms to do edge-thinning; the simplest ones
	// being the slowest.
	//
	// I tried that one but it's sooo slow:
	// https://users.fmrib.ox.ac.uk/~steve/susan/thinning/node2.html
	//
	// See http://article.sciencepublishinggroup.com/pdf/10.11648.j.ajsea.20130201.11.pdf
	// for an overview of other existing algorithms.

	// Here we use the modified version of the Zhang-Suen's algorithm outlined
	// in Kocharyan's paper.

	changed := true

	for changed {
		changed = thinEdgesIteration(m, true)
		changed = thinEdgesIteration(m, false) || changed
	}

	bounds := img.Bounds()
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if !m.get(x, y) {
				img.SetGray16(x+bounds.Min.X, y+bounds.Min.Y, black)
			}
		}
	}

	return img
}
