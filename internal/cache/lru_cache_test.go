package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLRUCache_SetAndGet тестирует Set и Get методы.
func TestLRUCache_SetAndGet(t *testing.T) {
	tempDir := t.TempDir()
	PreviewDir = tempDir

	c := NewLRUCache(2)

	err := c.Set("key1", []byte("image1"))
	require.NoError(t, err)

	value, ok := c.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, []byte("image1"), value)

	path := filepath.Join(tempDir, "key1")
	_, err = os.Stat(path)
	assert.NoError(t, err)
}

// TestLRUCache_RemoveLRU тестирует удаление наименее используемого элемента.
func TestLRUCache_RemoveLRU(t *testing.T) {
	tempDir := t.TempDir()
	PreviewDir = tempDir

	c := NewLRUCache(2)

	err := c.Set("key1", []byte("image1"))
	require.NoError(t, err)
	err = c.Set("key2", []byte("image2"))
	require.NoError(t, err)

	err = c.Set("key3", []byte("image3"))
	require.NoError(t, err)

	_, ok := c.Get("key1")
	assert.False(t, ok)

	path := filepath.Join(tempDir, "key1")
	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err))

	value, ok := c.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, []byte("image2"), value)

	value, ok = c.Get("key3")
	assert.True(t, ok)
	assert.Equal(t, []byte("image3"), value)
}

// TestLRUCache_FileOperations тестирует операции с файлами (сохранение и удаление).
func TestLRUCache_FileOperations(t *testing.T) {
	tempDir := t.TempDir()
	PreviewDir = tempDir

	c := NewLRUCache(1)

	err := c.Set("key1", []byte("image1"))
	require.NoError(t, err)

	path := filepath.Join(tempDir, "key1")
	_, err = os.Stat(path)
	assert.NoError(t, err)

	err = c.Set("key2", []byte("image2"))
	require.NoError(t, err)

	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err))

	newPath := filepath.Join(tempDir, "key2")
	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}
