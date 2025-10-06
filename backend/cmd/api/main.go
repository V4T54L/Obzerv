package main

import (
	"backend/internal/app"
	"backend/internal/infra/client"
	"backend/internal/infra/repo"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/V4T54L/goship/pkg/goship/server"
	"github.com/V4T54L/goship/pkg/goship/sse"
	"github.com/V4T54L/goship/pkg/goship/ws"
	"github.com/go-chi/chi/v5"
)

func main() {

	// Config
	logsPageSize := 50

	server := server.NewChiServer()
	server.AddDefaultMiddleware()
	server.AddPermissiveCORS()
	server.AddDefaultRoutes()

	r, ok := server.GetRouter().(*chi.Mux)
	if !ok {
		log.Fatal("Error obtaining the router")
	}

	logRepo := repo.NewLogRepository(logsPageSize)
	broadcaster := app.NewBroadcaster()

	uc := app.NewLogUsecase(logRepo, broadcaster)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start log generation/simulation
	_ = uc.Start(ctx)

	r.Get("/sse", sse.Handler(func(w sse.Writer, r *http.Request) error {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		w.Retry(2 * time.Second)

		logs, _ := uc.GetLastPage(ctx)
		data := ""
		for _, log := range logs {
			data += string(log) + "\n"
		}
		if err := w.Event("last_page", data); err != nil {
			return err
		}

		c := client.NewSSEClient(w)
		id := broadcaster.AddClient(c)
		defer broadcaster.RemoveClient(id)

		<-r.Context().Done()
		return nil
	}))

	r.Get("/ws", ws.Handler(func(conn ws.Conn) {
		log.Println("Client connected")
		c := client.NewWSClient(conn)

		// logs []string
		logs, _ := uc.GetLastPage(ctx)

		data, _ := json.Marshal(logs)
		_ = c.Send("last_page", data)

		id := broadcaster.AddClient(c)
		defer broadcaster.RemoveClient(id)

		// Keep connection open
		for {
			var msg any
			if err := conn.ReadJSON(&msg); err != nil {
				log.Println("WebSocket disconnected or error:", err)
				break
			}
		}

		log.Println("Client disconnected")
	}))

	// ---

	if err := server.Run("8000"); err != nil {
		log.Println("Error when running the server: ", err)
	}
}
