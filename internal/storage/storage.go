package storage

import (
	"errors"
	"fmt"
)

var (
	ArtifactNotFound = errors.New("artifact not found")
)

type Storager interface {
	Available() bool
	ArtifactExists(slug string, hash string) (bool, error)
	UploadArtifact(slug string, hash string, data []byte) error
	DownloadArtifact(slug string, hash string) ([]byte, error)
}

type Type string

const (
	Memory Type = "memory"
	Local  Type = "local"
	GCS    Type = "gcs"
)

func CreateStorage(typ Type) (Storager, error) {
	switch typ {
	case Memory:
		return newMemoryStorager(), nil
	case Local:
		// TODO: load basePath from config
		return newLocalStorager("./artifacts"), nil
	case GCS:
		//TODO load bucketName from config
		return newGcsStorager("sandbox-turbocache")
	default:
		return nil, fmt.Errorf("unknown storage type: %s", typ)
	}
}
