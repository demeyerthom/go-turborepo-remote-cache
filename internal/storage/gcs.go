package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"io"
)

type gcsStorager struct {
	bucketName string
	client     *storage.Client
}

func newGcsStorager(bucketName string) (*gcsStorager, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &gcsStorager{
		bucketName: bucketName,
		client:     client,
	}, nil
}

func (g *gcsStorager) Available() bool {
	ctx := context.Background()
	_, err := g.client.Bucket(g.bucketName).Attrs(ctx)
	return err == nil
}

func (g *gcsStorager) ArtifactExists(slug string, hash string) (bool, error) {
	ctx := context.Background()
	_, err := g.client.Bucket(g.bucketName).Object(g.createFilePath(slug, hash)).Attrs(ctx)
	if errors.Is(err, storage.ErrObjectNotExist) {
		return false, nil
	}
	return err == nil, err
}

func (g *gcsStorager) UploadArtifact(slug string, hash string, data []byte) error {
	ctx := context.Background()
	wc := g.client.Bucket(g.bucketName).Object(g.createFilePath(slug, hash)).NewWriter(ctx)
	if _, err := wc.Write(data); err != nil {
		return fmt.Errorf("failed to write data to GCS: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close GCS writer: %w", err)
	}
	return nil
}

func (g *gcsStorager) DownloadArtifact(slug string, hash string) ([]byte, error) {
	ctx := context.Background()
	rc, err := g.client.Bucket(g.bucketName).Object(g.createFilePath(slug, hash)).NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, ArtifactNotFound
		}
		return nil, fmt.Errorf("failed to create GCS reader: %w", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from GCS: %w", err)
	}
	return data, nil
}

func (g *gcsStorager) createFilePath(slug string, hash string) string {
	return fmt.Sprintf("%s/%s", sanitizeSlug(slug), hash)
}
