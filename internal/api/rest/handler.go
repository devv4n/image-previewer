package rest

import (
	"log/slog"
	"net/http"
	"strconv"
)

// PreviewHandler обрабатывает запросы на генерацию и получение превью.
// Пример URL: http://localhost:8080/fill/300/200/<image_url_encoded>
func (s *Server) PreviewHandler(w http.ResponseWriter, r *http.Request) {
	widthStr, heightStr, imageURL := r.PathValue("width"), r.PathValue("height"), r.PathValue("img_url")

	width, err := strconv.Atoi(widthStr)
	if err != nil || width <= 0 {
		http.Error(w, "Invalid 'width' parameter", http.StatusBadRequest)
		return
	}

	height, err := strconv.Atoi(heightStr)
	if err != nil || height <= 0 {
		http.Error(w, "Invalid 'height' parameter", http.StatusBadRequest)
		return
	}

	previewData, contentType, err := s.Service.GeneratePreview(r, width, height, imageURL)
	if err != nil {
		slog.Error("Error generating preview", "error", err, "imageURL", imageURL)
		http.Error(w, "Error while generating preview", http.StatusBadGateway)

		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(previewData)))

	if _, err = w.Write(previewData); err != nil {
		slog.Error("Error writing response", "error", err)
		return
	}
}
