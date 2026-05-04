package httpapi

import (
	"io"
	"log"
	"net/http"
	"pulselounge/internal/media"
	"sort"
	"strings"
)

const maxUploadBytes = 10 << 20

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
		writeJSONError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to retrieve file")
		return
	}
	defer file.Close()

	if handler.Size > maxUploadBytes {
		writeJSONError(w, http.StatusRequestEntityTooLarge, "file exceeds 10MB limit")
		return
	}

	if ok, err := hasAllowedImageContent(file); err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to read file")
		return
	} else if !ok {
		writeJSONError(w, http.StatusUnsupportedMediaType, "allowed file types: "+allowedUploadTypeList())
		return
	}

	// Get bytes of the file to save to blob storage
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to read file")
		return
	}
	url, err := h.store.Save(handler.Filename, fileBytes)
	if err != nil {
		log.Printf("error saving file: %s", err.Error())
		writeJSONError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	exists, err := h.store.Exists(handler.Filename)
	if err != nil {
		log.Printf("error checking file existence: %s", err.Error())
		writeJSONError(w, http.StatusInternalServerError, "failed to check file existence")
		return
	}

	response := struct {
		FileName string `json:"fileName"`
		FileSize int64  `json:"fileSize"`
		Url      string `json:"url"`
		Exists   bool   `json:"exists"`
	}{
		FileName: handler.Filename,
		FileSize: handler.Size,
		Url:      url,
		Exists:   exists,
	}

	writeJSON(w, http.StatusOK, response)
}

func hasAllowedImageContent(file io.ReadSeeker) (bool, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false, err
	}

	contentType := http.DetectContentType(buffer[:n])
	return allowedUploadTypes[contentType], nil
}

func allowedUploadTypeList() string {
	types := make([]string, 0, len(allowedUploadTypes))
	for contentType := range allowedUploadTypes {
		types = append(types, contentType)
	}

	sort.Strings(types)
	return strings.Join(types, ", ")
}
