package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"music/internal/middleware"
	"music/internal/repository"
	"music/internal/service"
	"net/http"
	"strconv"
)

type TrackHandler struct {
	svc       *service.TrackService
	queueRepo repository.QueueRepository
}

func NewTrackHandler(svc *service.TrackService, queueRepo repository.QueueRepository) *TrackHandler {
	return &TrackHandler{svc: svc, queueRepo: queueRepo}
}

func (h *TrackHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title  string `json:"title"`
		Prompt string `json:"prompt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	userId, _ := r.Context().Value("user_id").(int)

	track, err := h.svc.Create(r.Context(), req.Title, req.Prompt, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(track)
}

func (h *TrackHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	track, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(track)
}

func (h *TrackHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok || userID == 0 {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}
	log.Printf("List tracks for userID: %d", userID)
	tracks, err := h.svc.ListByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}
func (h *TrackHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title  string `json:"title"`
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(int)

	track, err := h.svc.Create(r.Context(), req.Title, req.Prompt, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.queueRepo.Enqueue(r.Context(), track.ID); err != nil {
		http.Error(w, "failed to enqueue generation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(track)
}

func (h *TrackHandler) CancelGeneration(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userId, _ := r.Context().Value(middleware.UserIDKey).(int)
	track, err := h.svc.GetByID(r.Context(), id)
	if err != nil || track.UserID != userId {
		http.Error(w, "track not found", http.StatusNotFound)
		return
	}

	if err := h.queueRepo.CancelByTrackID(r.Context(), id); err != nil {
		http.Error(w, "failed to cancel generation", http.StatusInternalServerError)
		return
	}

	if err := h.svc.UpdateStatus(r.Context(), id, "cancelled", ""); err != nil {
		http.Error(w, "failed to update track status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "cancelled"})
}

func (h *TrackHandler) Stream(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	track, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "track not found", http.StatusNotFound)
		return
	}

	if track.Status != "ready" || track.AudioURL == nil || *track.AudioURL == "" {
		http.Error(w, "track is not ready", http.StatusServiceUnavailable)
		return
	}

	resp, err := http.Get(*track.AudioURL)
	if err != nil {
		http.Error(w, "failed to fetch audio", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Accept-Ranges", "bytes")

	if r.Header.Get("Range") != "" {
		rangeHeader := r.Header.Get("Range")
		var start, end int64 = 0, -1
		fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)

		contentLength := resp.ContentLength
		if end == -1 || end >= contentLength {
			end = contentLength - 1
		}
		chunkSize := end - start + 1

		resp.Body.Close()
		resp, err = http.Get(*track.AudioURL)
		if err != nil {
			http.Error(w, "failed to fetch audio", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if start > 0 {
			io.CopyN(io.Discard, resp.Body, start)
		}

		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, contentLength))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", chunkSize))
		w.WriteHeader(http.StatusPartialContent)

		io.CopyN(w, resp.Body, chunkSize)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}
