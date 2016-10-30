package leonard

import (
	"image"
	"math"
)

const (
	// we only support 3 directions for now
	thetaV  int = iota // |
	thetaH             // -
	thetaNW            // \

	thetaCount
)

// TODO check out the Randomized Hough Transform:
// https://en.wikipedia.org/wiki/Randomized_Hough_transform
// http://homepages.inf.ed.ac.uk/rbf/CVonline/LOCAL_COPIES/AV1011/macdonald.pdf

// We store all points for each bin to know where to start/end each line
type houghBin []image.Point

type houghAccumulator struct {
	bins     [][thetaCount]houghBin
	binWidth int
}

func newHoughAccumulator(b *BinaryImage, binWidth int) *houghAccumulator {
	minR := 0
	// max length of the r parameter; that is if it's a north-west ->
	// south-east diagonal at the top-right of the image.
	maxR := b.width * b.height

	// Arbitrary value. Higher values mean faster but more approximate results.
	// An odd value is more suitable because we'll draw a line in the middle of
	// it.
	bins := (maxR-minR)/binWidth + 1

	return &houghAccumulator{
		bins:     make([][thetaCount]houghBin, bins),
		binWidth: binWidth,
	}
}

func (acc *houghAccumulator) Inc(x, y, theta int) {
	var bin int

	switch theta {
	case thetaV:
		bin = y / acc.binWidth
	case thetaH:
		bin = x / acc.binWidth
	case thetaNW:
		bin = int(math.Sqrt(float64(x*x+y*y))) / acc.binWidth
	default:
		panic("Invalid theta")
	}

	acc.bins[bin][theta] = append(acc.bins[bin][theta], image.Point{x, y})
}

func (acc houghAccumulator) Max() (m int) {
	for _, bin := range acc.bins {
		for _, v := range bin {
			if l := len(v); l > m {
				m = l
			}
		}
	}
	return m
}

// HoughTransform performs a Hough Transform on the image and return an
// (r, theta) accumulator.
func (b *BinaryImage) HoughTransform() *houghAccumulator {
	// Algorithm: https://en.wikipedia.org/wiki/Hough_transform#Implementation
	// For now we only check vertical & horizontal lines

	// Arbitrary bin width. Higher values mean faster but more approximate
	// results. An odd value is more suitable because we'll draw a line in the
	// middle of it.
	accumulator := newHoughAccumulator(b, 5)

	b.EachPixel(func(x, y int) {
		if b.Get(north.apply(x, y)) && b.Get(south.apply(x, y)) {
			// vertical
			accumulator.Inc(x, y, thetaV)
		}

		if b.Get(west.apply(x, y)) && b.Get(east.apply(x, y)) {
			// horizontal
			accumulator.Inc(x, y, thetaH)
		}

		if b.Get(northwest.apply(x, y)) && b.Get(southeast.apply(x, y)) {
			// north-west -> south-east
			accumulator.Inc(x, y, thetaNW)
		}
	})

	return accumulator
}

// DrawLines takes an (r, theta) accumulator as returned by HoughTransform and
// draw the corresponding lines on the image.
func (b *BinaryImage) DrawLines(acc *houghAccumulator) {
	threshold := acc.binWidth * 5 // arbitrary

	for r, ts := range acc.bins {
		for theta, n := range ts {
			if len(n) < threshold {
				continue
			}

			// draw a line in the middle of the bin
			off := r*acc.binWidth + acc.binWidth/2

			if theta == thetaV {
				for y := 0; y < b.height; y++ {
					b.Set(off, y, true)
				}
			} else if theta == thetaH {
				for x := 0; x < b.width; x++ {
					b.Set(x, off, true)
				}
			}

			// TODO diagonals NW-SE
		}
	}
}
