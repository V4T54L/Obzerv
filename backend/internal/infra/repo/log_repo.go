package repo

import (
	"backend/internal/domain"
	"backend/pkg/utils"
	"context"
	"sync"
)

type logRepoImpl struct {
	mu   sync.RWMutex
	logs []domain.Log
	idx  int // next write index
}

// NewLogRepository creates a circular buffer of logs
func NewLogRepository(pageSize int) domain.LogRepository {
	logs := make([]domain.Log, pageSize)
	for i := range logs {
		log := utils.GenerateLogEntry()
		logs[i] = domain.Log(log)
	}
	return &logRepoImpl{
		logs: logs,
		idx:  0,
	}
}

// AddLog adds a log to the circular buffer, replacing the oldest log
func (r *logRepoImpl) AddLog(ctx context.Context, log domain.Log) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logs[r.idx] = log
	r.idx = (r.idx + 1) % len(r.logs)
	return nil
}

// GetLastPage returns all logs in chronological order (oldest → newest)
func (r *logRepoImpl) GetLastPage(ctx context.Context) ([]domain.Log, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Log, 0, len(r.logs))

	// Start from the current idx (oldest log), and go through all logs in order
	for i := 0; i < len(r.logs); i++ {
		pos := (r.idx + i) % len(r.logs)
		result = append(result, r.logs[pos])
	}

	return result, nil
}
