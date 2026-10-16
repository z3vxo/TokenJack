package events

import (
	"errors"
	"sync"
)

var ErrNoClient = errors.New("events: no client for id")

const clientBuffer = 16

type Client struct {
	ID     string
	events chan []byte
}

func (c *Client) Events() <-chan []byte {
	return c.events
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]struct{}
}

func New() *Hub {
	return &Hub{
		clients: make(map[string]map[*Client]struct{}),
	}
}

func (h *Hub) Subscribe(id string) *Client {
	client := &Client{ID: id, events: make(chan []byte, clientBuffer)}

	h.mu.Lock()
	defer h.mu.Unlock()
	set, ok := h.clients[id]
	if !ok {
		set = make(map[*Client]struct{})
		h.clients[id] = set
	}
	set[client] = struct{}{}

	return client
}

func (h *Hub) Unsubscribe(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	set, ok := h.clients[client.ID]
	if !ok {
		return
	}
	delete(set, client)
	if len(set) == 0 {
		delete(h.clients, client.ID)
	}
}

func (h *Hub) Send(id string, message []byte) error {
	clients := h.clientsFor(id)
	if len(clients) == 0 {
		return ErrNoClient
	}
	for _, client := range clients {
		client.push(message)
	}
	return nil
}

func (h *Hub) Broadcast(message []byte) {
	for _, client := range h.allClients() {
		client.push(message)
	}
}

func (h *Hub) clientsFor(id string) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	set := h.clients[id]
	clients := make([]*Client, 0, len(set))
	for client := range set {
		clients = append(clients, client)
	}
	return clients
}

func (h *Hub) allClients() []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var clients []*Client
	for _, set := range h.clients {
		for client := range set {
			clients = append(clients, client)
		}
	}
	return clients
}

func (c *Client) push(message []byte) {
	select {
	case c.events <- message:
	default:
	}
}
