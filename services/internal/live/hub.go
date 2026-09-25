package live

import "sync"

// sendBuffer is how many messages a client may lag behind before it is dropped.
const sendBuffer = 32

// Client is one WebSocket connection subscribed to a match.
type Client struct {
	send chan []byte
}

// Messages is closed when the client is dropped or the hub shuts down.
func (c *Client) Messages() <-chan []byte {
	return c.send
}

// Hub fans out messages of each match to its subscribers.
// Many HTTP goroutines use it at once, so its map is guarded by a mutex.
type Hub struct {
	mu      sync.Mutex
	matches map[string]map[*Client]struct{}

	firstOnce sync.Once
	first     chan struct{}
}

func NewHub() *Hub {
	return &Hub{
		matches: map[string]map[*Client]struct{}{},
		first:   make(chan struct{}),
	}
}

// FirstSubscriber is closed when the first client subscribes to any match.
func (h *Hub) FirstSubscriber() <-chan struct{} {
	return h.first
}

func (h *Hub) Subscribe(matchID string) *Client {
	c := &Client{send: make(chan []byte, sendBuffer)}

	h.mu.Lock()
	if h.matches[matchID] == nil {
		h.matches[matchID] = map[*Client]struct{}{}
	}
	h.matches[matchID][c] = struct{}{}
	h.mu.Unlock()

	h.firstOnce.Do(func() { close(h.first) })
	return c
}

func (h *Hub) Unsubscribe(matchID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.remove(matchID, c)
}

// Publish sends msg to every subscriber of the match without waiting.
// A client whose buffer is full is too slow and gets dropped,
// so one slow phone never holds up the others.
func (h *Hub) Publish(matchID string, msg []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for c := range h.matches[matchID] {
		select {
		case c.send <- msg:
		default:
			h.remove(matchID, c)
		}
	}
}

// Close drops every client of every match.
func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for matchID, clients := range h.matches {
		for c := range clients {
			h.remove(matchID, c)
		}
	}
}

// remove must be called with h.mu held.
func (h *Hub) remove(matchID string, c *Client) {
	if _, ok := h.matches[matchID][c]; !ok {
		return
	}
	delete(h.matches[matchID], c)
	close(c.send)
}
