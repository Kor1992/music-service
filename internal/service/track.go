package service

import (
	"context"
	"errors"
	"music/internal/domain"
	"music/internal/repository"
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

func (s *TrackService) UpdateStatus(ctx context.Context, id int, status, audioURL string) error {
	return s.repo.UpdateStatus(ctx, id, status, audioURL)
}

func (s *TrackService) GenerateTrack(ctx context.Context, trackID int) {
	// track, err := s.repo.GetByID(ctx, trackID)
	// if err != nil {
	// 	log.Printf("GenerateTrack: failed to get track %d: %v", trackID, err)
	// 	return
	// }

	// apiURL := "https://api.riffusion.com/v1/generate"
	// body := fmt.Sprintf(`{"prompt":"%s","duration":30}`, track.Prompt)
	// req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(body))
	// if err != nil {
	// 	log.Printf("Riffusion request failed: %v", err)
	// 	s.repo.UpdateStatus(ctx, trackID, "failed", "")
	// 	return
	// }
	// req.Header.Set("Content-Type", "application/json")

	// resp, err := http.DefaultClient.Do(req)

	// if err != nil {
	// 	log.Printf("Riffusion API error: %v", err)
	// 	s.repo.UpdateStatus(ctx, trackID, "failed", "")
	// 	return
	// }
	// defer resp.Body.Close()

	// var riffResp struct {
	// 	AudioURL string `json:"audio_url"`
	// }

	// json.NewDecoder(resp.Body).Decode(&riffResp)

	// s.repo.UpdateStatus(ctx, trackID, "ready", riffResp.AudioURL)
	testAudioURL := "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3"
	s.repo.UpdateStatus(ctx, trackID, "ready", testAudioURL)
}
