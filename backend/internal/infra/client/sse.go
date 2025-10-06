package client

import (
	"backend/extra/newsse"
	"backend/internal/domain"
	"errors"
)

type sseClient struct {
	writer newsse.Writer
	done   chan struct{}
}

func NewSSEClient(w newsse.Writer) domain.BroadcastClient {
	return &sseClient{writer: w, done: make(chan struct{})}
}

func (c *sseClient) Send(event string, data []byte) error {
	select {
	case <-c.done:
		return errors.New("client closed")
	default:
		return c.writer.Event(event, string(data))
	}
}

func (c *sseClient) Close() error {
	select {
	case <-c.done:
		return nil
	default:
		close(c.done)
		return c.writer.Close()
	}
}
