package newsse

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Writer interface {
	Retry(dur time.Duration) error
	Event(event, data string) error
	Close() error
	Comment(comment string) error
	Flush() error
}

type sseWriter struct {
	w       http.ResponseWriter
	flusher http.Flusher
	mu      sync.Mutex
	closed  bool
}

func newSSEWriter(w http.ResponseWriter) (*sseWriter, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming unsupported: ResponseWriter is not a Flusher")
	}

	// set SSE headers (do this before any Write)
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // for nginx

	// explicitly send headers
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	return &sseWriter{w: w, flusher: flusher}, nil
}

func (s *sseWriter) writeRaw(b []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return errors.New("writer closed")
	}
	_, err := s.w.Write(b)
	if err != nil {
		return err
	}
	s.flusher.Flush()
	return nil
}

func (s *sseWriter) Event(event, data string) error {
	// SSE data lines must be prefixed with "data: "
	var buf bytes.Buffer
	if event != "" {
		buf.WriteString("event: ")
		buf.WriteString(event)
		buf.WriteString("\n")
	}
	// preserve newlines in data (split, then "data: " each line)
	for _, line := range strings.Split(strings.ReplaceAll(data, "\r\n", "\n"), "\n") {
		buf.WriteString("data: ")
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	buf.WriteString("\n") // end of event
	return s.writeRaw(buf.Bytes())
}

func (s *sseWriter) Retry(d time.Duration) error {
	err := s.writeRaw([]byte("retry: " + strconv.FormatInt(int64(d/time.Millisecond), 10) + "\n\n"))
	return err
}

func (s *sseWriter) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

func (s *sseWriter) Comment(comment string) error {
	// Per SSE spec: comments start with ":"
	var buf bytes.Buffer
	buf.WriteString(": ")
	buf.WriteString(comment)
	buf.WriteString("\n\n") // End of event
	return s.writeRaw(buf.Bytes())
}

func (s *sseWriter) Flush() error {
	return nil
}
