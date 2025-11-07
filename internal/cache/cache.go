package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// DefaultRefreshInterval is the default cache refresh interval
	DefaultRefreshInterval = 10 * time.Minute

	// CacheDir is the directory where cache files are stored
	CacheDir = ".local/k4a/cache"
)

// Manager handles caching of kafkactl responses
type Manager struct {
	cacheDir string
}

// NewManager creates a new cache manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, CacheDir)

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Manager{
		cacheDir: cacheDir,
	}, nil
}

// getCacheKey generates a cache key from resource type and context
func (m *Manager) getCacheKey(resourceType string, args ...string) string {
	// Create a unique key based on resource type and arguments
	key := resourceType
	for _, arg := range args {
		key += "_" + arg
	}

	// Hash the key to create a safe filename
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// GetCachePath returns the full path to a cache file
func (m *Manager) GetCachePath(resourceType string, args ...string) string {
	filename := m.getCacheKey(resourceType, args...) + ".yaml"
	return filepath.Join(m.cacheDir, filename)
}

// Get retrieves data from cache if it exists and is fresh
func (m *Manager) Get(resourceType string, maxAge time.Duration, args ...string) ([]byte, bool, error) {
	cachePath := m.GetCachePath(resourceType, args...)

	// Check if cache file exists
	info, err := os.Stat(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil // Cache miss
		}
		return nil, false, fmt.Errorf("failed to stat cache file: %w", err)
	}

	// Check if cache is still fresh
	if time.Since(info.ModTime()) > maxAge {
		return nil, false, nil // Cache expired
	}

	// Read cache file
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read cache file: %w", err)
	}

	return data, true, nil
}

// Set stores data in cache
func (m *Manager) Set(resourceType string, data []byte, args ...string) error {
	cachePath := m.GetCachePath(resourceType, args...)

	// Write data to cache file
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// Invalidate removes a specific cache entry
func (m *Manager) Invalidate(resourceType string, args ...string) error {
	cachePath := m.GetCachePath(resourceType, args...)

	err := os.Remove(cachePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cache file: %w", err)
	}

	return nil
}

// InvalidateAll removes all cache entries
func (m *Manager) InvalidateAll() error {
	entries, err := os.ReadDir(m.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			path := filepath.Join(m.cacheDir, entry.Name())
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove cache file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// GetCacheAge returns the age of a cache entry
func (m *Manager) GetCacheAge(resourceType string, args ...string) (time.Duration, error) {
	cachePath := m.GetCachePath(resourceType, args...)

	info, err := os.Stat(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("cache entry does not exist")
		}
		return 0, fmt.Errorf("failed to stat cache file: %w", err)
	}

	return time.Since(info.ModTime()), nil
}
