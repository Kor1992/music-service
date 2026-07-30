package service

import (
	"context"
	"errors"
	"log"
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

func (s *TrackService) ListByUser(ctx context.Context, userID int) ([]*domain.Track, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *TrackService) UpdateStatus(ctx context.Context, id int, status, audioURL string) error {
	return s.repo.UpdateStatus(ctx, id, status, audioURL)
}

func (s *TrackService) GenerateTrack(ctx context.Context, trackID int) {
	log.Printf("GenerateTrack started for track %d", trackID)

	// Заглушка: сразу ставим статус ready и тестовую ссылку
	testAudioURL := "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3"
	log.Printf("Setting status ready with URL: %s", testAudioURL)

	err := s.repo.UpdateStatus(ctx, trackID, "ready", testAudioURL)
	if err != nil {
		log.Printf("GenerateTrack: failed to update status: %v", err)
	} else {
		log.Printf("GenerateTrack: successfully updated track %d to ready", trackID)
	}
	// log.Println("Generate handler called")
	// track, err := s.repo.GetByID(ctx, trackID)
	// if err != nil {
	// 	log.Printf("GenerateTrack: failed to get track %d: %v", trackID, err)
	// 	return
	// }

	// apiURL := "https://api.replicate.com/v1/predictions"
	// requestBody := fmt.Sprintf(`{
	//     "version": "8fa09c7d0e1b38c5e7adc0cf44c5e6e6e7c3d15f7e7b8f7c4d4e7a8a3b3e8c8",
	//     "input": {
	//         "prompt": %q,
	//         "duration": 30
	//     }
	// }`, track.Prompt)

	// req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(requestBody))

	// if err != nil {
	// 	log.Printf("Replicate request failed: %v", err)
	// 	s.repo.UpdateStatus(ctx, trackID, "failed", "")
	// 	return
	// }

	// req.Header.Set("Authorization", "Token "+os.Getenv("REPLICATE_API_TOKEN"))
	// req.Header.Set("Content-Type", "application/json")

	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	log.Printf("Replicate API error: %v", err)
	// 	s.repo.UpdateStatus(ctx, trackID, "failed", "")
	// 	return
	// }

	// defer resp.Body.Close()

	// var prediction struct {
	// 	ID     string `json:"id"`
	// 	Status string `json:"status"`
	// 	Output struct {
	// 		Audio string `json:"audio"`
	// 	} `json:"output"`
	// }

	// json.NewDecoder(resp.Body).Decode(&prediction)

	// for prediction.Status != "succeeded" && prediction.Status != "failed" {
	// 	time.Sleep(5 * time.Second)
	// 	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL+"/"+prediction.ID, nil)
	// 	req.Header.Set("Authorization", "Token "+os.Getenv("REPLICATE_API_TOKEN"))
	// 	resp, _ := http.DefaultClient.Do(req)
	// 	if resp != nil {
	// 		json.NewDecoder(resp.Body).Decode(&prediction)
	// 		resp.Body.Close()
	// 	}
	// }

	// if prediction.Status == "succeeded" && prediction.Output.Audio != "" {
	// 	s.repo.UpdateStatus(ctx, trackID, "ready", prediction.Output.Audio)
	// } else {
	// 	s.repo.UpdateStatus(ctx, trackID, "failed", "")
	// }
}
