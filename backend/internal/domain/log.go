package domain

import "context"

type Log string

type LogRepository interface {
	AddLog(ctx context.Context, log Log) error
	GetLastPage(ctx context.Context) ([]Log, error)
}
