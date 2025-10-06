package client

import (
	"backend/internal/domain"

	"github.com/V4T54L/goship/pkg/goship/ws"
)

type wsClient struct {
	conn ws.Conn
}

func NewWSClient(conn ws.Conn) domain.BroadcastClient {
	return &wsClient{conn: conn}
}

func (c *wsClient) Send(event string, data []byte) error {
	// You could define a standard envelope format for events
	msg := map[string]interface{}{
		"type": event,
		"data": string(data),
	}
	return c.conn.WriteJSON(msg)
}

func (c *wsClient) Close() error {
	return c.conn.Close()
}
