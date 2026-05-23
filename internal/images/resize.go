package images

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"math"
	"time"

	"golang.org/x/image/draw"
)

func ResizeImage(img image.Image, width int) ([]byte, error) {
	start := time.Now()
	defer func() {
		log.Printf("resizing image took %v", time.Since(start))
	}()
	ratio := (float64)(img.Bounds().Max.Y) / (float64)(img.Bounds().Max.X)
	newHeight := int(math.Round(float64(width) * ratio))

	dst := image.NewRGBA(image.Rect(0, 0, width, newHeight))

	draw.CatmullRom.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

	var out bytes.Buffer
	if err := png.Encode(&out, dst); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
