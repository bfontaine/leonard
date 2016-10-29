package leonard

import (
	"image"
	"image/color"
	"math"
)

func gaussianKernel(x float64, sigma float64) float64 {
	// The Gaussian filter is the convolution between a kernel and the image
	// matrix [1,2,3,4,5,6].
	//
	// This convolution is separable, i.e. instead of using a matrix MxM we use
	// 1xM and Mx1, effectively applying the filter first horizontally then
	// vertically [7].
	//
	// See also Grigory Dryapak's concurrent implementation:
	//     https://github.com/disintegration/imaging/blob/master/effects.go
	//
	// Formula from:
	//     http://www.stat.wisc.edu/~mchung/teaching/MIA/reading/diffusion.gaussian.kernel.pdf.pdf
	//
	// [1] http://homepages.inf.ed.ac.uk/rbf/HIPR2/gsmooth.htm
	// [2] https://en.wikipedia.org/wiki/Canny_edge_detector#Gaussian_filter
	// [3] http://homepages.inf.ed.ac.uk/rbf/HIPR2/convolve.htm
	// [4] http://aishack.in/tutorials/image-convolution-examples/
	// [5] http://graphics.cs.cmu.edu/courses/15-463/2005_fall/www/Lectures/convolution.pdf
	// [6] https://www.cs.cornell.edu/courses/cs6670/2011sp/lectures/lec02_filter.pdf
	// [7] http://www.songho.ca/dsp/convolution/convolution.html#separable_convolution
	return math.Exp(-(x * x / (2 * sigma * sigma))) / (sigma * math.Sqrt(2*math.Pi))
}

// GaussianFilter applies a gaussian filter with the given sigma parameter on
// the image.
func GaussianFilter(img image.Image, sigma float64) image.Image {
	// The radius should grow with sigma. Mathematica uses a factor of 2 [1]
	// while G. Dryapak uses 3 [2].
	//
	// [1] http://dsp.stackexchange.com/a/10067/24352
	// [2] https://github.com/disintegration/imaging/blob/5b7e226/effects.go#L26
	radius := int(math.Ceil(sigma * 2.0))

	// <center pixel> + radius, hence radius+1
	kernel := make([]float64, radius+1)

	for i := 0; i <= radius; i++ {
		kernel[i] = gaussianKernel(float64(i), sigma)
	}

	bounds := img.Bounds()
	blured := image.NewRGBA(bounds)

	// columns
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		start := x - radius
		if start < 0 {
			start = 0
		}

		end := x + radius
		if end >= bounds.Max.X {
			end = bounds.Max.X - 1
		}

		weightsSum := 0.0

		for ix := start; ix <= end; ix++ {
			weightsSum += kernel[absint(x-ix)]
		}

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			var r, g, b, a float64

			for ix := start; ix <= end; ix++ {
				ir, ig, ib, ia := img.At(ix, y).RGBA()

				weight := kernel[absint(x-ix)]
				r += weight * to255(float64(ir))
				g += weight * to255(float64(ig))
				b += weight * to255(float64(ib))
				a += weight * to255(float64(ia))
			}

			blured.Set(x, y, color.RGBA{
				uint8(clip(r/weightsSum, 0.0, 255.0)),
				uint8(clip(g/weightsSum, 0.0, 255.0)),
				uint8(clip(b/weightsSum, 0.0, 255.0)),
				uint8(clip(a/weightsSum, 0.0, 255.0)),
			})
		}
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		start := y - radius
		if start < 0 {
			start = 0
		}

		end := y + radius
		if end >= bounds.Max.Y {
			end = bounds.Max.Y - 1
		}

		weightsSum := 0.0

		for iy := start; iy <= end; iy++ {
			weightsSum += kernel[absint(y-iy)]
		}

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a float64

			for iy := start; iy <= end; iy++ {
				ir, ig, ib, ia := blured.At(x, iy).RGBA()

				weight := kernel[absint(y-iy)]
				r += weight * to255(float64(ir))
				g += weight * to255(float64(ig))
				b += weight * to255(float64(ib))
				a += weight * to255(float64(ia))
			}

			blured.Set(x, y, color.RGBA{
				uint8(clip(r/weightsSum, 0.0, 255.0)),
				uint8(clip(g/weightsSum, 0.0, 255.0)),
				uint8(clip(b/weightsSum, 0.0, 255.0)),
				uint8(clip(a/weightsSum, 0.0, 255.0)),
			})
		}
	}

	return blured
}
