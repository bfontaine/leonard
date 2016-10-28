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
	return ThinEdges(gradients(img, func(img image.Image, x, y int) float64 {
		// Read e.g. http://www.cse.psu.edu/~rtc12/CSE486/lecture02.pdf
		// Also: https://www.cs.umd.edu/~djacobs/CMSC426/ImageGradients.pdf
		//       https://en.wikipedia.org/wiki/Image_gradient
		h := horizontalGradient(img, x, y)
		v := verticalGradient(img, x, y)
		return math.Sqrt(h*h + v*v)
	}))
}

// Return true if the point at the given coordinates is the sole connection
// between some of its neighbours; i.e. if it splits its neighborhood in two
// when removed.
//
// Note: We use the Moore neighborhood: https://en.wikipedia.org/wiki/Moore_neighborhood
//
// Examples:
//
//     ooo            o..            o..            .o.             ooo
//     oxo -> false   ox. -> false   .x. -> true    oxo -> false    .x. -> true
//     ooo            o..            ooo            ...             ooo
//
func isSoleConnection(neighborhood *boolMatrix, x, y int) bool {

	// FIXME This is wrong; it checks if there's an isolated neighbour instead
	// of looking for two components.

	// for each neighbour...
	for _, no := range neighboursOffsets {
		hasNeighbour := false
		px, py := no.apply(1, 1)

		// check each of its neighboors...
		for _, nno := range neighboursOffsets {
			// if it has at least one it's not isolated: stop
			if neighborhood.get(nno.apply(px, py)) {
				hasNeighbour = true
				break
			}
		}

		if !hasNeighbour {
			return true
		}
	}

	return false
}

func edgePointShouldBeRemoved(m *boolMatrix, x, y int) bool {
	// Build the neighborhood; excluding self
	neighborhood := newBoolMatrix(3, 3)
	for _, no := range neighboursOffsets {
		mx, my := no.apply(x, y)
		nx, ny := no.apply(1, 1)

		neighborhood.set(nx, ny, m.get(mx, my))
	}

	// see http://homepages.inf.ed.ac.uk/rbf/HIPR2/thin.htm

	count := neighborhood.count(true)

	// Not at a region boundary
	if count == 8 {
		return false
	}

	// isolated
	if count == 0 {
		return true
	}

	return !isSoleConnection(neighborhood, x, y)
}

// Thin the edges of an image that went through Gradients() and return it.
// The image is modified in-place.
func ThinEdges(img *image.Gray16) *image.Gray16 {
	m := toBoolMatrix(img)

	// NOTE: We could probably optimize the neighbors lookup we perform for
	// *each* pixel at *each* step

	passes := 4 // this should be configurable

	for p := 0; p < passes; p++ {
		for y := 0; y < m.height; y++ {
			for x := 0; x < m.width; x++ {
				if edgePointShouldBeRemoved(m, x, y) {
					m.set(x, y, false)
				}
			}
		}
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
