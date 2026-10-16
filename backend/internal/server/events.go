package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

const eventPingInterval = 25 * time.Second

func (s *Server) HandleEvents(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	token := bearerToken(r)

	if id == "" || token == "" {
		s.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "missing id or auth token"})
		return
	}

	if _, err := s.Auth.Verify(token); err != nil {
		s.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		s.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "streaming unsupported"})
		return
	}

	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	header.Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	client := s.Events.Subscribe(id)
	defer s.Events.Unsubscribe(client)

	ctx := r.Context()
	ping := time.NewTicker(eventPingInterval)
	defer ping.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ping.C:
			if _, err := fmt.Fprint(w, ": ping\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case message := <-client.Events():
			if _, err := fmt.Fprintf(w, "data: %s\n\n", message); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func bearerToken(r *http.Request) string {
	scheme, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	if ok && strings.EqualFold(scheme, "Bearer") {
		return token
	}
	return r.URL.Query().Get("auth")
}
