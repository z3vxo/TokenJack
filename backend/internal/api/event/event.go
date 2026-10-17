package event

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/z3vxo/TokenJack/internal/api"
	"github.com/z3vxo/TokenJack/internal/auth"
	"github.com/z3vxo/TokenJack/internal/events"
)

const pingInterval = 25 * time.Second

type Handler struct {
	Auth   *auth.Auth
	Events *events.Hub
}

func New(a *auth.Auth, hub *events.Hub) *Handler {
	return &Handler{Auth: a, Events: hub}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", h.Stream)
	return r
}

func (h *Handler) Stream(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	token := bearerToken(r)

	if id == "" || token == "" {
		api.Error(w, http.StatusUnauthorized, "missing id or auth token")
		return
	}

	if _, err := h.Auth.Verify(token); err != nil {
		api.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		api.Error(w, http.StatusInternalServerError, "streaming unsupported")
		return
	}

	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	header.Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	client := h.Events.Subscribe(id)
	defer h.Events.Unsubscribe(client)

	ctx := r.Context()
	ping := time.NewTicker(pingInterval)
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
