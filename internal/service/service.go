package service

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/disintegration/imaging"
)

// Cache интерфейс кэша.
type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte) error
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service struct {
	Cache  Cache
	Client HTTPClient
}

// NewService создаёт новый сервис.
func NewService(cache Cache) *Service {
	return &Service{
		Cache:  cache,
		Client: &http.Client{},
	}
}

// GeneratePreview скачивает изображение, изменяет его размеры и сохраняет в кэш.
func (s *Service) GeneratePreview(req *http.Request, width, height int, imageURL string) ([]byte, string, error) {
	key := fmt.Sprintf("%d_%d_%s", width, height, url.QueryEscape(imageURL))

	// Check the cache
	if data, ok := s.Cache.Get(key); ok {
		slog.Debug("cache hit for key", "key", key)
		return data, "", nil
	}

	slog.Debug("cache miss", "key", key)

	// Prepare the request for the image with metadata from the client request.
	clientReq, err := http.NewRequestWithContext(req.Context(), http.MethodGet, imageURL, http.NoBody)
	if err != nil {
		return nil, "", fmt.Errorf("error creating request for image: %w", err)
	}

	// Copy relevant headers from the original request.
	for header, values := range req.Header {
		for _, value := range values {
			clientReq.Header.Add(header, value)
		}
	}

	resp, err := s.Client.Do(clientReq)
	if err != nil {
		return nil, "", fmt.Errorf("error downloading image: %w", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			slog.Error("error while response body closing", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("remote server returned status: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/jpg" {
		return nil, "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error decoding image: %w", err)
	}

	preview := imaging.Resize(img, width, height, imaging.Lanczos)

	var buf bytes.Buffer

	switch contentType {
	case "image/png":
		if err = png.Encode(&buf, preview); err != nil {
			return nil, "", fmt.Errorf("error while encoding PNG preview: %w", err)
		}
	case "image/jpeg", "image/jpg":
		if err = jpeg.Encode(&buf, preview, nil); err != nil {
			return nil, "", fmt.Errorf("error while encoding JPEG preview: %w", err)
		}
	}

	previewData := buf.Bytes()

	if err = s.Cache.Set(key, previewData); err != nil {
		return nil, "", fmt.Errorf("error while caching image: %w", err)
	}

	slog.Debug("preview successfully generated and cached", "key", key)

	return previewData, contentType, nil
}
