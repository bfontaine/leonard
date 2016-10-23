package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/bfontaine/leonard/leonard"
)

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

	transformed := leonard.VerticalGradients(img)

	out, err := os.Create("out.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if err = writePNG(transformed, out); err != nil {
		log.Fatal(err)
	}
}
