package images

import (
	"bytes"
	"image"
	"image/png"

	"github.com/nfnt/resize"
)

// ResizeImage takes an image and resizes it to the specified width and height, returning the resized image as a PNG byte slice.
// It uses nfnt/resize for resizing, which is an archived library. In the future we may want to consider switching to a more actively maintained library for image resizing.
func ResizeImage(image image.Image, width, height int) ([]byte, error) {
	resizedImg := resize.Resize(uint(width), uint(height), image, resize.Lanczos3)

	var out bytes.Buffer
	if err := png.Encode(&out, resizedImg); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
