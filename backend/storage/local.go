package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage implements StorageProvider using the local filesystem.
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage instance rooted at basePath.
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

func (s *LocalStorage) fullPath(filename string) string {
	return filepath.Join(s.basePath, filename)
}

func (s *LocalStorage) Save(filename string, data []byte) error {
	path := s.fullPath(filename)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *LocalStorage) LoadOnce(filename string) ([]byte, error) {
	return os.ReadFile(s.fullPath(filename))
}

func (s *LocalStorage) LoadStream(filename string) (io.ReadCloser, error) {
	return os.Open(s.fullPath(filename))
}

func (s *LocalStorage) Download(filename string, targetPath string) error {
	data, err := s.LoadOnce(filename)
	if err != nil {
		return err
	}
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(targetPath, data, 0644)
}

func (s *LocalStorage) Exists(filename string) (bool, error) {
	_, err := os.Stat(s.fullPath(filename))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *LocalStorage) Delete(filename string) error {
	err := os.Remove(s.fullPath(filename))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
