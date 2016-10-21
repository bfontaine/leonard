package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
)

func toGrayscale(img image.Image) image.Image {
	grayscaled := image.NewRGBA(img.Bounds())

	bd := img.Bounds()
	for y := bd.Min.Y; y < bd.Max.Y; y++ {
		for x := bd.Min.X; x < bd.Max.X; x++ {
			// https://en.wikipedia.org/wiki/Grayscale#Converting_color_to_grayscale
			r, g, b, _ := img.At(x, y).RGBA()
			linear := 0.2126*float32(r) + 0.7152*float32(g) + 0.0722*float32(b)

			y_linear := uint16(linear)

			grayscaled.Set(x, y, color.Gray16{y_linear})
		}
	}

	return grayscaled
}

func writePNG(img image.Image, w io.Writer) error {
	return png.Encode(w, img)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage:\n\t%s <image>\n\n", os.Args[0])
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	gray := toGrayscale(img)

	out, err := os.Create("out.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if err = writePNG(gray, out); err != nil {
		log.Fatal(err)
	}
}
