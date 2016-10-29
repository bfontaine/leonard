package leonard

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
)

// SaveImage saves an image to a file. The correct format is guessed from the
// file extension. The default format for files without extension is PNG.
func SaveImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	switch ext := strings.ToLower(filepath.Ext(filename)); ext {
	case "jpg", "jpeg":
		return jpeg.Encode(f, img, nil)
	case "", "png":
		return png.Encode(f, img)
	case "gif":
		return gif.Encode(f, img, nil)
	default:
		return errors.New(fmt.Sprintf("Unknown format: %s", ext))
	}
}

// LoadImage loads an image from a file. If the filename is "-" the image is
// read from os.Stdin.
func LoadImage(filename string) (image.Image, error) {
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
