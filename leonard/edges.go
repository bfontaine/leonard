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

// HorizontalGradients returns an image that represents the magnitude of the
// horizontal gradients
func HorizontalGradients(img image.Image) image.Image {
	return gradients(img, horizontalGradient)
}

// VerticalGradients returns an image that represents the magnitude of the
// vertical gradients
func VerticalGradients(img image.Image) image.Image {
	return gradients(img, verticalGradient)
}

// Gradients returns an image that represents the magnitude of gradients
func Gradients(img image.Image) image.Image {
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
	b.EachPixel(func(x, y int) {
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
			goto preservePixel
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
					goto preservePixel
				}

				// (d)
				if p4 && p6 && p8 {
					goto preservePixel
				}
			} else {
				// second subiteration

				// (c')
				if p2 && p4 && p8 {
					goto preservePixel
				}

				// (d')
				if p2 && p6 && p8 {
					goto preservePixel
				}
			}
		case 2:
			// Kocharyan (2013)'s modifications
			if odd {
				if !p4 || !p2 || p3 || p7 || p8 {
					goto preservePixel
				}
			} else {
				if !p6 || !p8 || p3 || p4 || p7 {
					goto preservePixel
				}
			}
		}

		// delete the pixel
		changed = true
		return

	preservePixel:
		b2.Set(x, y, true)
	})
	return b2, changed
}

// ThinEdges thins the edges of an image [that went through Gradients()] and
// return it. The image is modified in-place.
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

	b2 := b.Clone()

	for changed1 || changed2 {
		b2, changed1 = b2.thinEdgesIteration(true)
		b2, changed2 = b2.thinEdgesIteration(false)
	}

	// Modify in-place
	b.pixels = b2.pixels

	return b
}

const (
	tetaV = 0
	tetaH = 1
)

// 2 = number of teta we support. For now only vertical & horizontal
type houghAccumulator struct {
	bins     [][2]uint8
	binWidth int
}

func newHoughAccumulator(bins, binWidth int) *houghAccumulator {
	return &houghAccumulator{
		bins:     make([][2]uint8, bins),
		binWidth: binWidth,
	}
}

func (acc *houghAccumulator) Inc(bin, teta int) {
	acc.bins[bin][teta]++
}

func (acc houghAccumulator) Max() (m uint8) {
	for _, bin := range acc.bins {
		for _, v := range bin {
			if v > m {
				m = v
			}
		}
	}
	return m
}

// HoughTransform performs a Hough Transform on the image and return an
// (r, teta) accumulator.
func (b *BinaryImage) HoughTransform() *houghAccumulator {
	// Algorithm: https://en.wikipedia.org/wiki/Hough_transform#Implementation
	// For now we only check vertical & horizontal lines

	minR := 0
	// max length of the r parameter; that is if it's a north-west ->
	// south-east diagonal at the top-right of the image.
	maxR := b.width * b.height

	// Arbitrary value. Higher values mean faster but more approximate results.
	binWidth := 6
	bins := (maxR-minR)/binWidth + 1

	accumulator := newHoughAccumulator(bins, binWidth)

	b.EachPixel(func(x, y int) {
		// FIXME we should maybe store min/max coordinates for each bin so we
		// know where to start & end the lines
		// -> problem: there might be multiple lines in the same bin
		// -> solution: store all coordinates then run some heuristics
		// (*magic*)

		if b.Get(north.apply(x, y)) && b.Get(south.apply(x, y)) {
			// vertical
			accumulator.Inc(x/binWidth, tetaV)
		}

		if b.Get(west.apply(x, y)) && b.Get(east.apply(x, y)) {
			// horizontal
			accumulator.Inc(y/binWidth, tetaH)
		}
	})

	return accumulator
}

// DrawLines takes an (r, teta) accumulator as returned by HoughTransform and
// draw the corresponding lines on the image.
func (b *BinaryImage) DrawLines(acc *houghAccumulator) {
	threshold := uint8(acc.binWidth * 5) // arbitrary

	for r, ts := range acc.bins {
		for teta, n := range ts {
			if n < threshold {
				continue
			}

			// draw a line in the middle of the bin
			rm := r*acc.binWidth + acc.binWidth/2
			if teta == tetaV {
				for y := 0; y < b.height; y++ {
					b.Set(rm, y, true)
				}
			} else if teta == tetaH {
				for x := 0; x < b.width; x++ {
					b.Set(x, rm, true)
				}
			}
		}
	}
}
