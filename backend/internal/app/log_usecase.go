package app

import (
	"backend/internal/domain"
	"backend/pkg/utils"
	"context"
	"log"
	"time"
)

type LogUsecase interface {
	Start(ctx context.Context) error
	GetLastPage(ctx context.Context) ([]domain.Log, error)
}

type logUsecaseImpl struct {
	repo        domain.LogRepository
	broadcaster *Broadcaster
}

func NewLogUsecase(repo domain.LogRepository, b *Broadcaster) LogUsecase {
	return &logUsecaseImpl{repo: repo, broadcaster: b}
}

func (uc *logUsecaseImpl) Start(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Context canceled, stopping log addition.")
				ticker.Stop() // Close the ticker
			case <-ticker.C:
				log := utils.GenerateLogEntry()

				uc.repo.AddLog(ctx, domain.Log(log))
				uc.broadcaster.Broadcast("log", []byte(log))
			}
		}
	}()

	return nil
}

func (uc *logUsecaseImpl) GetLastPage(ctx context.Context) ([]domain.Log, error) {
	return uc.repo.GetLastPage(ctx)
}
