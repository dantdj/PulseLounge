package httpapi

import (
	"errors"
	"io"
	"log"
	"net/http"
	"pulselounge/internal/images"
	"pulselounge/internal/media"
	"sort"
	"strings"

	"github.com/google/uuid"
)

const maxUploadBytes = 10 << 20

var errUnsupportedUploadType = errors.New("unsupported upload type")

// Using a map for O(1) lookup. This isn't strictly necessary,
// but useful for the future.
var allowedUploadTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

type UploadHandler struct {
	store media.Store
}

func NewUploadHandler(store media.Store) UploadHandler {
	return UploadHandler{store: store}
}

func (h UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)

	err := r.ParseMultipartForm(maxUploadBytes)
	if err != nil {
		log.Printf("failed to parse upload form: %v", err)
		writeJSONError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to retrieve file")
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("error closing uploaded file: %s", closeErr.Error())
		}
	}()

	if handler.Size > maxUploadBytes {
		writeJSONError(w, http.StatusRequestEntityTooLarge, "file exceeds 10MB limit")
		return
	}

	contentType, err := detectAllowedUploadType(file)
	if err != nil {
		if errors.Is(err, errUnsupportedUploadType) {
			writeJSONError(w, http.StatusUnsupportedMediaType, "allowed file types: "+allowedUploadTypeList())
			return
		}

		writeJSONError(w, http.StatusBadRequest, "failed to read file")
		return
	}

	normalizedImageBytes, img, err := images.NormalizeImageToPNG(file, contentType)
	if err != nil {
		log.Printf("failed to normalize uploaded file %q: %v", handler.Filename, err)
		writeJSONError(w, http.StatusInternalServerError, "failed to normalize file")
		return
	}

	id := uuid.New().String() + ".png"

	url, err := h.store.Save(id, "image/png", normalizedImageBytes)
	if err != nil {
		log.Printf("failed to save uploaded file %q as %q: %v", handler.Filename, id, err)
		writeJSONError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	// Best effort attempt to make a thumbnail, if it fails we don't want to fail the whole upload
	// The UI should handle missing thumbnails gracefully
	thumbnailBytes, err := images.ResizeImage(img, 200, 200)
	if err != nil {
		log.Printf("failed to create thumbnail for uploaded file %q: %v", handler.Filename, err)
	}
	if thumbnailBytes != nil {
		thumbnailID := media.ThumbnailKey(id)
		// We don't need the URL here so we can ignore it
		if _, err := h.store.Save(thumbnailID, "image/png", thumbnailBytes); err != nil {
			log.Printf("failed to save thumbnail for uploaded file %q: %v", handler.Filename, err)
		}
	}

	response := struct {
		Url string `json:"url"`
		Key string `json:"key"`
	}{
		Key: id,
		Url: url,
	}

	writeJSON(w, http.StatusOK, response)
}

// detectAllowedUploadType reads the first 512 bytes of the file to determine its content type,
// returning an error if the content type cannot be determined or is not allowed.
func detectAllowedUploadType(file io.ReadSeeker) (contentType string, err error) {
	var buffer [512]byte

	n, readErr := file.Read(buffer[:])
	if readErr != nil && readErr != io.EOF {
		return "", readErr
	}

	if _, seekErr := file.Seek(0, io.SeekStart); seekErr != nil {
		return "", seekErr
	}

	contentType = http.DetectContentType(buffer[:n])
	if !allowedUploadTypes[contentType] {
		return "", errUnsupportedUploadType
	}

	return contentType, nil
}

func allowedUploadTypeList() string {
	types := make([]string, 0, len(allowedUploadTypes))
	for contentType := range allowedUploadTypes {
		types = append(types, contentType)
	}

	sort.Strings(types)
	return strings.Join(types, ", ")
}
