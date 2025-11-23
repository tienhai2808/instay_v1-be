package hub

import (
	"encoding/json"
	"sync"

	"github.com/InstaySystem/is-be/internal/types"
)

type SSEClient struct {
	ID         string
	ClientID   int64
	Type       string
	Department *string
	Chan       chan []byte
	Done       chan bool
}

type SSEHub struct {
	Clients    map[string]*SSEClient
	Register   chan *SSEClient
	Unregister chan *SSEClient
	Broadcast  chan []byte
	Mutex      sync.RWMutex
}

func NewSSEHub() *SSEHub {
	return &SSEHub{
		Clients:    make(map[string]*SSEClient),
		Register:   make(chan *SSEClient),
		Unregister: make(chan *SSEClient),
		Broadcast:  make(chan []byte),
	}
}

func (h *SSEHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client.ID] = client
			h.Mutex.Unlock()

		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Chan)
				close(client.Done)
			}
			h.Mutex.Unlock()

		case message := <-h.Broadcast:
			h.Mutex.RLock()
			for _, client := range h.Clients {
				select {
				case client.Chan <- message:
				default:
					delete(h.Clients, client.ID)
					close(client.Chan)
					close(client.Done)
				}
			}
			h.Mutex.RUnlock()
		}
	}
}

func (h *SSEHub) SendToClient(clientID int64, event types.SSEEventData) {
	data, _ := json.Marshal(event)

	if event.Type == "staff" && event.Department != nil {
		for _, client := range h.Clients {
			if client.ClientID == clientID && client.Department != nil && client.Department == event.Department {
				select {
				case client.Chan <- data:
				default:
					h.Mutex.RUnlock()
					h.Mutex.Lock()
					if _, ok := h.Clients[client.ID]; ok {
						delete(h.Clients, client.ID)
						close(client.Chan)
						close(client.Done)
					}
					h.Mutex.Unlock()
					h.Mutex.RLock()
				}
			}
		}
	} else if event.Type == "guest" && event.Department == nil {
		for _, client := range h.Clients {
			if client.ClientID == clientID && client.Department == nil && client.Type == event.Type {
				select {
				case client.Chan <- data:
				default:
					h.Mutex.RUnlock()
					h.Mutex.Lock()
					if _, ok := h.Clients[client.ID]; ok {
						delete(h.Clients, client.ID)
						close(client.Chan)
						close(client.Done)
					}
					h.Mutex.Unlock()
					h.Mutex.RLock()
				}
			}
		}
	}
	h.Mutex.RUnlock()
}
