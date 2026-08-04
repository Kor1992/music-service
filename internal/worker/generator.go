package worker

import (
	"context"
	"log"
	"music/internal/repository"
	"music/internal/service"
	"time"
)

type GeneratorWorker struct {
	queueRepo    repository.QueueRepository
	trackService *service.TrackService
}

func NewGeneratorWorker(queueRepo repository.QueueRepository, trackService *service.TrackService) *GeneratorWorker {
	return &GeneratorWorker{
		queueRepo:    queueRepo,
		trackService: trackService,
	}
}

func (w *GeneratorWorker) Start(ctx context.Context) {
	log.Println("Generator worker started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Generator worker stopped")
			return
		default:
			item, err := w.queueRepo.Dequeue(ctx)
			if err != nil {
				log.Printf("Worker: dequeue error: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			if item == nil {
				time.Sleep(2 * time.Second)
				continue
			}

			track, err := w.trackService.GetByID(ctx, item.TrackID)
			if err != nil || track.Status == "cancelled" {
				w.queueRepo.UpdateStatus(ctx, item.ID, "cancelled")
				continue
			}

			log.Printf("Worker: processing task %d for track %d", item.ID, item.TrackID)
			w.trackService.GenerateTrack(ctx, item.TrackID)

			w.queueRepo.UpdateStatus(ctx, item.ID, "completed")
		}
	}
}
