package client

import (
	"backend/internal/domain"

	"github.com/V4T54L/goship/pkg/goship/sse"
)

type sseClient struct {
	writer sse.Writer
}

func NewSSEClient(writer sse.Writer) domain.BroadcastClient {
	return &sseClient{writer: writer}
}

func (c *sseClient) Send(event string, data []byte) error {
	return c.writer.Event(event, string(data))
}

func (c *sseClient) Close() error {
	// no-op for SSE, we rely on context cancellation
	return nil
}
