package app

import (
	"backend/internal/domain"
	"log"
	"sync"
)

type Broadcaster struct {
	mu      sync.RWMutex
	clients map[int]domain.BroadcastClient
	counter int
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[int]domain.BroadcastClient),
	}
}

func (b *Broadcaster) AddClient(client domain.BroadcastClient) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	id := b.counter
	b.counter++
	b.clients[id] = client
	return id
}

func (b *Broadcaster) RemoveClient(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if client, ok := b.clients[id]; ok {
		_ = client.Close()
		delete(b.clients, id)
	}
}

func (b *Broadcaster) Broadcast(event string, data []byte) {
	b.mu.RLock()
	clients := make(map[int]domain.BroadcastClient, len(b.clients))
	for id, c := range b.clients {
		clients[id] = c
	}
	b.mu.RUnlock()

	for id, c := range clients {
		if err := c.Send(event, data); err != nil {
			log.Printf("Broadcast error: %v (removing client %d)", err, id)
			go b.RemoveClient(id)
		}
	}
}
