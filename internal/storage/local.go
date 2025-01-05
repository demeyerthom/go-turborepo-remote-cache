package storage

import (
	"errors"
	"os"
	"path/filepath"
)

type localStorager struct {
	basePath string
}

func newLocalStorager(basePath string) *localStorager {
	return &localStorager{
		basePath: basePath,
	}
}

func (l *localStorager) Available() bool {
	testFilePath := filepath.Join(l.basePath, "testfile.tmp")
	if err := os.WriteFile(testFilePath, []byte("test"), 0644); err != nil {
		return false
	}
	if err := os.Remove(testFilePath); err != nil {
		return false
	}
	return true
}

func (l *localStorager) ArtifactExists(slug string, hash string) (bool, error) {
	path := l.createFilePath(slug, hash)
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func (l *localStorager) UploadArtifact(slug string, hash string, data []byte) error {
	path := l.createFilePath(slug, hash)
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (l *localStorager) DownloadArtifact(slug string, hash string) ([]byte, error) {
	path := l.createFilePath(slug, hash)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ArtifactNotFound
		}
		return nil, err
	}
	return data, nil
}

func (l *localStorager) createFilePath(slug string, hash string) string {
	return filepath.Join(l.basePath, sanitizeSlug(slug), hash)
}
