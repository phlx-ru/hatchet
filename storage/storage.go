package storage

import (
	"os"
	"path"
	"path/filepath"
)

type Storage struct {
	basePath string
}

func New(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

func (s *Storage) Save(content []byte, pathChunks ...string) error {
	pathChunks = append([]string{s.basePath}, pathChunks...)
	filePath := path.Join(pathChunks...)

	dir := filepath.Dir(filePath)
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, content, os.ModePerm)
}

func (s *Storage) Read(pathChunks ...string) ([]byte, error) {
	pathChunks = append([]string{s.basePath}, pathChunks...)
	filePath := path.Join(pathChunks...)
	return os.ReadFile(filePath) //nolint:gosec // Mute `Potential file inclusion via variable` cause of basePath
}
