package images

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/webp"
)

// Takes an image reader and returns a normalized version of the image without metadata.
// The returned image data will be a PNG format image, regardless of what is passed in.
// Also returns the decoded image.Image for use in generating thumbnails, etc.
func NormalizeImageToPNG(r io.Reader, contentType string) ([]byte, image.Image, error) {
	var (
		img image.Image
		err error
	)
	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(r)
	case "image/png":
		img, err = png.Decode(r)
	case "image/webp":
		img, err = webp.Decode(r)
	default:
		return nil, nil, fmt.Errorf("unsupported image content type: %s", contentType)
	}

	if err != nil {
		return nil, nil, err
	}

	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		return nil, nil, err
	}

	return out.Bytes(), img, nil
}
