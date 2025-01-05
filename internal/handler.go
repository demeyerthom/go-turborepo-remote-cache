package internal

import (
	"encoding/json"
	"errors"
	"github.com/demeyerthom/go-turborepo-remote-cache/internal/storage"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ServerHandler struct {
	storager storage.Storager
}

func NewServerHandler(storager storage.Storager) *ServerHandler {
	return &ServerHandler{
		storager: storager,
	}
}

type Status string

const (
	StatusEnabled  Status = "enabled"
	StatusDisabled Status = "disabled"
	//StatusOverLimit Status = "over_limit"
	//StatusPaused    Status = "paused"
)

type StatusResponse struct {
	Status Status `json:"status"`
}

func (h *ServerHandler) Status(w http.ResponseWriter, r *http.Request, params StatusParams) {
	resp := StatusResponse{
		Status: StatusEnabled,
	}

	if h.storager.Available() {
		resp.Status = StatusDisabled
	}

	d, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(d)
	return
}

func (h *ServerHandler) ArtifactQuery(w http.ResponseWriter, r *http.Request, params ArtifactQueryParams) {
	log.Warnf("ArtifactQuery not implemented")
	_, _ = w.Write([]byte("{}"))
	w.WriteHeader(http.StatusOK)
	return
}

func (h *ServerHandler) RecordEvents(w http.ResponseWriter, r *http.Request, params RecordEventsParams) {
	log.Warnf("RecordEvents not implemented")
	w.WriteHeader(http.StatusOK)
	return
}

func (h *ServerHandler) DownloadArtifact(w http.ResponseWriter, r *http.Request, hash string, params DownloadArtifactParams) {
	slug := GetSlug(params.TeamId, params.Slug)
	if slug == nil {
		http.Error(w, "invalid teamId or slug", http.StatusBadRequest)
		return
	}

	data, err := h.storager.DownloadArtifact(*slug, hash)
	if err != nil {
		if errors.Is(err, storage.ArtifactNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Errorf("failed to download artifact: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debugf("downloaded artifact: %s", hash)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
	return
}

func (h *ServerHandler) ArtifactExists(w http.ResponseWriter, r *http.Request, hash string, params ArtifactExistsParams) {
	slug := GetSlug(params.TeamId, params.Slug)
	if slug == nil {
		http.Error(w, "invalid teamId or slug", http.StatusBadRequest)
		return
	}

	exists, err := h.storager.ArtifactExists(*slug, hash)
	if err != nil {
		log.Errorf("failed to check if artifact exists: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Debugf("artifact does not exist: %s", hash)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Debugf("artifact exists: %s", hash)
	w.WriteHeader(http.StatusOK)
	return
}

func (h *ServerHandler) UploadArtifact(w http.ResponseWriter, r *http.Request, hash string, params UploadArtifactParams) {
	slug := GetSlug(params.TeamId, params.Slug)
	if slug == nil {
		http.Error(w, "invalid teamId or slug", http.StatusBadRequest)
		return
	}

	data := make([]byte, r.ContentLength)
	_, err := r.Body.Read(data)
	if err != nil {
		log.Errorf("failed to read request body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.storager.UploadArtifact(*slug, hash, data)
	if err != nil {
		log.Errorf("failed to upload artifact: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debugf("uploaded artifact: %s", hash)
	w.WriteHeader(http.StatusOK)
	return
}
