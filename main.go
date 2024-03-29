package main

import (
	"fmt"
	"image"
	"os"

	"github.com/bfontaine/leonard/leonard"
	"gopkg.in/urfave/cli.v1"
)

var transformFuncs = map[string]func(image.Image) image.Image{
	"gray":       leonard.Grayscale,
	"binary":     leonard.Binary,
	"vgradients": leonard.VerticalGradients,
	"hgradients": leonard.HorizontalGradients,
	"gradients":  leonard.Gradients,
	"downscale":  leonard.Downscale,
	"blur": func(i image.Image) image.Image {
		return leonard.GaussianFilter(i, 1.4)
	},
	"edges": func(i image.Image) image.Image {
		b := leonard.NewBinaryImage(
			leonard.Gradients(leonard.GaussianFilter(i, 5.0)), -1)

		b.ThinEdges()

		// acc := b.HoughTransform()
		// b.DrawLines(acc)

		return b
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
			for name := range transformFuncs {
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

		img, err := leonard.LoadImage(c.Args().First())

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

		if err := leonard.SaveImage(img, c.Args().Get(1)); err != nil {
			return cli.NewExitError(fmt.Sprintf("Write error: %s", err), 1)
		}
		return nil
	}

	app.Run(os.Args)
}
