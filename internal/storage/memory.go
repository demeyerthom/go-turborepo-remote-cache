package storage

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type memoryStorager struct {
	data map[string][]byte
}

func newMemoryStorager() *memoryStorager {
	data := make(map[string][]byte)

	// Print all the keys in the data map every second
	go func() {
		for {
			time.Sleep(time.Second * 5)
			if len(data) == 0 {
				log.Debugf("No keys")
			} else {
				k := make([]string, 0, len(data))
				for key := range data {
					k = append(k, key)
				}
				log.WithField("keys", k).Debugf("Keys in memory")
			}
		}
	}()

	return &memoryStorager{
		data: data,
	}
}

func (m *memoryStorager) Available() bool {
	return true
}

func (m *memoryStorager) ArtifactExists(slug string, hash string) (bool, error) {
	if _, ok := m.data[createMemoryKey(slug, hash)]; ok {
		return true, nil
	}
	return false, nil
}

func (m *memoryStorager) UploadArtifact(slug string, hash string, data []byte) error {
	m.data[createMemoryKey(slug, hash)] = data
	return nil
}

func (m *memoryStorager) DownloadArtifact(slug string, hash string) ([]byte, error) {
	if data, ok := m.data[createMemoryKey(slug, hash)]; ok {
		return data, nil
	}
	return nil, ArtifactNotFound
}

func createMemoryKey(slug string, hash string) string {
	return strings.Join([]string{slug, hash}, "-")
}
