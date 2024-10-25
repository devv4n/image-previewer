package service

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) ([]byte, bool) {
	args := m.Called(key)
	if args.Get(0) != nil {
		return args.Get(0).([]byte), args.Bool(1)
	}

	return nil, args.Bool(1)
}

func (m *MockCache) Set(key string, value []byte) error {
	args := m.Called(key, value)
	return args.Error(0)
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if resp := args.Get(0); resp != nil {
		return resp.(*http.Response), args.Error(1)
	}

	return nil, args.Error(1)
}

// TestGeneratePreview_Success тестирует успешную генерацию превью.
func TestGeneratePreview_Success(t *testing.T) {
	cache := new(MockCache)
	client := new(MockHTTPClient)

	img := imaging.New(100, 100, image.White)

	var imgBuf bytes.Buffer
	err := jpeg.Encode(&imgBuf, img, nil)
	assert.NoError(t, err)

	cache.On("Get", mock.Anything).Return(nil, false)
	cache.On("Set", mock.Anything, mock.Anything).Return(nil)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(imgBuf.Bytes())),
		Header:     http.Header{"Content-Type": []string{"image/jpeg"}},
	}
	client.On("Do", mock.Anything).Return(resp, nil)

	svc := NewService(cache)
	svc.Client = client

	req := httptest.NewRequest(http.MethodGet, "https://example.com", http.NoBody)

	result, err := svc.GeneratePreview(req, 50, 50, "https://example.com/image.jpg")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	cache.AssertExpectations(t)
	client.AssertExpectations(t)
}

// TestGeneratePreview_CacheHit тестирует случай кэш-хита.
func TestGeneratePreview_CacheHit(t *testing.T) {
	cache := new(MockCache)
	client := new(MockHTTPClient)

	cache.On("Get", mock.Anything).Return([]byte("cached image data"), true)

	svc := NewService(cache)
	svc.Client = client

	req := httptest.NewRequest(http.MethodGet, "https://example.com", http.NoBody)

	result, err := svc.GeneratePreview(req, 50, 50, "https://example.com/image.jpg")
	assert.NoError(t, err)
	assert.Equal(t, []byte("cached image data"), result)

	client.AssertNotCalled(t, "Do", mock.Anything)
	cache.AssertExpectations(t)
}

// TestGeneratePreview_UnsupportedContentType тестирует случай с неподдерживаемым типом контента.
func TestGeneratePreview_UnsupportedContentType(t *testing.T) {
	cache := new(MockCache)
	client := new(MockHTTPClient)

	cache.On("Get", mock.Anything).Return(nil, false)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("image data")),
		Header:     http.Header{"Content-Type": []string{"image/gif"}},
	}
	client.On("Do", mock.Anything).Return(resp, nil)

	svc := NewService(cache)
	svc.Client = client

	req := httptest.NewRequest(http.MethodGet, "https://example.com", http.NoBody)

	_, err := svc.GeneratePreview(req, 50, 50, "https://example.com/image.gif")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported content type")

	cache.AssertExpectations(t)
	client.AssertExpectations(t)
}

// TestGeneratePreview_DownloadError тестирует случай ошибки загрузки изображения.
func TestGeneratePreview_DownloadError(t *testing.T) {
	cache := new(MockCache)
	client := new(MockHTTPClient)

	cache.On("Get", mock.Anything).Return(nil, false)

	client.On("Do", mock.Anything).Return(nil, errors.New("download error"))

	svc := NewService(cache)
	svc.Client = client

	req := httptest.NewRequest(http.MethodGet, "https://example.com", http.NoBody)

	_, err := svc.GeneratePreview(req, 50, 50, "https://example.com/image.jpg")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "download error")

	cache.AssertExpectations(t)
	client.AssertExpectations(t)
}
