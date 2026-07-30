package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"music/internal/domain"
	"music/internal/repository"
	"net/http"
	"os"
)

type TrackService struct {
	repo repository.TrackRepository
}

func NewTrackService(repo repository.TrackRepository) *TrackService {
	return &TrackService{
		repo: repo,
	}
}

func (s *TrackService) Create(ctx context.Context, title, prompt string, userID int) (*domain.Track, error) {
	switch {
	case title == "":
		return nil, errors.New("title cannot be empty")
	case prompt == "":
		return nil, errors.New("prompt cannot be empty")
	}

	return s.repo.Create(ctx, title, prompt, userID)
}

func (s *TrackService) GetByID(ctx context.Context, id int) (*domain.Track, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TrackService) List(ctx context.Context) ([]*domain.Track, error) {
	return s.repo.List(ctx)
}

func (s *TrackService) ListByUser(ctx context.Context, userID int) ([]*domain.Track, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *TrackService) UpdateStatus(ctx context.Context, id int, status, audioURL string) error {
	return s.repo.UpdateStatus(ctx, id, status, audioURL)
}

func (s *TrackService) GenerateTrack(ctx context.Context, trackID int) {
	track, err := s.repo.GetByID(ctx, trackID)
	if err != nil {
		log.Println(err)
		return
	}

	body := map[string]any{
		"version": "meta/musicgen:671ac645ce5e552cc63a54a2bbff63fcf798043055d2dac5fc9e36a837eedcfb",
		"input": map[string]any{
			"prompt":                 track.Prompt,
			"model_version":          "stereo-large",
			"output_format":          "mp3",
			"normalization_strategy": "peak",
		},
	}

	data, _ := json.Marshal(body)

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.replicate.com/v1/predictions",
		bytes.NewReader(data),
	)
	if err != nil {
		log.Println(err)
		s.repo.UpdateStatus(ctx, trackID, "failed", "")
		return
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("REPLICATE_API_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Prefer", "wait")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		s.repo.UpdateStatus(ctx, trackID, "failed", "")
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated {

		log.Println(string(bodyBytes))
		s.repo.UpdateStatus(ctx, trackID, "failed", "")
		return
	}

	type Prediction struct {
		Status string `json:"status"`
		Output string `json:"output"`
	}

	var prediction Prediction

	if err := json.Unmarshal(bodyBytes, &prediction); err != nil {
		log.Println(err)
		s.repo.UpdateStatus(ctx, trackID, "failed", "")
		return
	}

	if prediction.Status != "succeeded" {
		log.Println(string(bodyBytes))
		s.repo.UpdateStatus(ctx, trackID, "failed", "")
		return
	}

	if prediction.Output == "" {
		log.Println("empty output")
		s.repo.UpdateStatus(ctx, trackID, "failed", "")
		return
	}

	if err := s.repo.UpdateStatus(
		ctx,
		trackID,
		"ready",
		prediction.Output,
	); err != nil {
		log.Println(err)
	}
}
