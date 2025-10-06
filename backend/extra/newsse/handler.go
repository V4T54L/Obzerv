package newsse

import (
	"log"
	"net/http"
)

// Handler wraps a function that uses the SSE Writer.
func Handler(fn func(Writer, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sseW, err := newSSEWriter(w)
		if err != nil {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}
		// Allow handler to use the writer
		if err := fn(sseW, r); err != nil {
			log.Printf("SSE handler error: %v", err)
			// ensure close
			_ = sseW.Close()
			return
		}
		_ = sseW.Close()
	}
}
