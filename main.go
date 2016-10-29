package main

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/bfontaine/leonard/leonard"
	"gopkg.in/urfave/cli.v1"
)

func decodeImage(filename string) (image.Image, error) {
	var f *os.File

	if filename == "-" {
		f = os.Stdin
	} else {
		var err error
		if f, err = os.Open(filename); err != nil {
			return nil, err
		}
		defer f.Close()
	}

	img, _, err := image.Decode(f)

	return img, err
}

func encodeImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	switch ext := strings.ToLower(filepath.Ext(filename)); ext {
	case "jpg", "jpeg":
		return jpeg.Encode(f, img, nil)
	case "png":
		return png.Encode(f, img)
	case "gif":
		return gif.Encode(f, img, nil)
	default:
		return errors.New(fmt.Sprintf("Unknown format: %s", ext))
	}

}

var transformFuncs = map[string]func(image.Image) image.Image{
	"gray":       leonard.Grayscale,
	"binary":     leonard.Binary,
	"vgradients": leonard.VerticalGradients,
	"hgradients": leonard.HorizontalGradients,
	"gradients":  leonard.Gradients,
	"blur":       func(i image.Image) image.Image { return leonard.GaussianFilter(i, 1.4) },
	"edges": func(i image.Image) image.Image {
		// arbitrary
		threshold := 0.16 * 0xFFFF
		return leonard.NewBinaryImage(
			leonard.Gradients(
				leonard.GaussianFilter(i, 5.0)), int(threshold)).ThinEdges()
	},
}

func main() {

	app := cli.NewApp()
	app.Name = "Leonard"
	app.Usage = "Apply various transforms on images"
	app.UsageText = "leonard [options] <image> <output image>"
	// No "help" command, please. Unfortunately this also hides the flags.
	app.HideHelp = true
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "transform, t",
			Usage: "Transformation to apply. You can chain them.",
		},
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "List the available transformations",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("list") {
			for name, _ := range transformFuncs {
				fmt.Println(name)
			}
			return nil
		}

		switch c.NArg() {
		case 0:
			return cli.NewExitError("Please give me an image.", 1)
		case 1:
			return cli.NewExitError("Please give me an output file.", 1)
		}

		img, err := decodeImage(c.Args().First())

		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Decoding error: %s", err), 1)
		}

		for _, t := range c.StringSlice("transform") {
			fn, ok := transformFuncs[t]
			if !ok {
				return cli.NewExitError(
					fmt.Sprintf("Unknown transform '%s'", t), 1)
			}
			img = fn(img)
		}

		if err := encodeImage(img, c.Args().Get(1)); err != nil {
			return cli.NewExitError(fmt.Sprintf("Write error: %s", err), 1)
		}
		return nil
	}

	app.Run(os.Args)
}
