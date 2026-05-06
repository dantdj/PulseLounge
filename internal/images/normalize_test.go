package images

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"
)

func TestNormalizeImageToPNGReturnsPNG(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		input       []byte
	}{
		{
			name:        "jpeg",
			contentType: "image/jpeg",
			input:       createTestJPEG(t),
		},
		{
			name:        "png",
			contentType: "image/png",
			input:       createTestPNG(t),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeImageToPNG(bytes.NewReader(tt.input), tt.contentType)
			if err != nil {
				t.Fatalf("NormalizeImageToPNG returned error: %v", err)
			}

			if !bytes.HasPrefix(got, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}) {
				t.Fatal("expected normalized image to have a PNG signature")
			}

			if _, err := png.Decode(bytes.NewReader(got)); err != nil {
				t.Fatalf("failed to decode normalized png: %v", err)
			}
		})
	}
}

func TestNormalizeImageToPNGStripsJPEGExif(t *testing.T) {
	const exifMarker = "PulseLounge GPS metadata test marker"

	jpegWithExif := addJPEGExif(t, createTestJPEG(t), []byte(exifMarker))
	if !bytes.Contains(jpegWithExif, []byte(exifMarker)) {
		t.Fatal("expected jpeg fixture to contain test EXIF marker")
	}

	got, err := NormalizeImageToPNG(bytes.NewReader(jpegWithExif), "image/jpeg")
	if err != nil {
		t.Fatalf("NormalizeImageToPNG returned error: %v", err)
	}

	if bytes.Contains(got, []byte(exifMarker)) {
		t.Fatal("expected normalized image to strip EXIF marker")
	}
	if bytes.Contains(got, []byte("Exif\x00\x00")) {
		t.Fatal("expected normalized image to strip EXIF header")
	}
	if _, err := png.Decode(bytes.NewReader(got)); err != nil {
		t.Fatalf("failed to decode normalized png: %v", err)
	}
}

func TestNormalizeImageToPNGRejectsUnsupportedContentType(t *testing.T) {
	_, err := NormalizeImageToPNG(bytes.NewReader(createTestPNG(t)), "image/gif")
	if err == nil {
		t.Fatal("expected error for unsupported content type")
	}
	if !strings.Contains(err.Error(), "unsupported image content type: image/gif") {
		t.Fatalf("expected unsupported content type error, got %v", err)
	}
}

func TestNormalizeImageToPNGReturnsDecodeErrors(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		input       []byte
	}{
		{
			name:        "invalid jpeg",
			contentType: "image/jpeg",
			input:       []byte("not a jpeg"),
		},
		{
			name:        "invalid png",
			contentType: "image/png",
			input:       []byte("not a png"),
		},
		{
			name:        "invalid webp",
			contentType: "image/webp",
			input:       []byte("not a webp"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NormalizeImageToPNG(bytes.NewReader(tt.input), tt.contentType)
			if err == nil {
				t.Fatal("expected decode error")
			}
		})
	}
}

func createTestJPEG(t *testing.T) []byte {
	t.Helper()

	src := image.NewGray(image.Rect(0, 0, 1, 1))

	var input bytes.Buffer
	if err := jpeg.Encode(&input, src, nil); err != nil {
		t.Fatalf("encode jpeg fixture: %v", err)
	}

	return input.Bytes()
}

func createTestPNG(t *testing.T) []byte {
	t.Helper()

	src := image.NewGray(image.Rect(0, 0, 1, 1))

	var input bytes.Buffer
	if err := png.Encode(&input, src); err != nil {
		t.Fatalf("encode png fixture: %v", err)
	}

	return input.Bytes()
}

func addJPEGExif(t *testing.T, jpegData []byte, payload []byte) []byte {
	t.Helper()

	if len(jpegData) < 2 || jpegData[0] != 0xff || jpegData[1] != 0xd8 {
		t.Fatal("jpeg fixture is missing SOI marker")
	}

	exif := append([]byte("Exif\x00\x00"), payload...)
	segmentLength := len(exif) + 2
	if segmentLength > 0xffff {
		t.Fatal("test EXIF payload is too large")
	}

	out := make([]byte, 0, len(jpegData)+len(exif)+4)
	out = append(out, jpegData[:2]...)
	out = append(out, 0xff, 0xe1, byte(segmentLength>>8), byte(segmentLength))
	out = append(out, exif...)
	out = append(out, jpegData[2:]...)
	return out
}
