package hub

import (
	"encoding/json"
	"fmt"
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
	evtDept := "nil"
	if event.Department != nil {
		evtDept = *event.Department
	}
	fmt.Printf("\n[SSE-DEBUG] ---> Start SendToClient. TargetID: %d | EventType: %s | EventDept: %s\n", clientID, event.Type, evtDept)

	h.Mutex.RLock()

	fmt.Printf("[SSE-DEBUG] Current connected clients: %d\n", len(h.Clients))

	data, _ := json.Marshal(event)

	if event.Type == "staff" && event.Department != nil {
		for _, client := range h.Clients {
			clientDept := "nil"
			if client.Department != nil {
				clientDept = *client.Department
			}

			if client.ClientID == clientID {
				fmt.Printf("[SSE-DEBUG] Found matching ClientID in Hub. HubClient: {Type: %s, Dept: %s} vs Event: {Type: %s, Dept: %s}\n",
					client.Type, clientDept, event.Type, evtDept)

				isDeptMatch := client.Department != nil && *client.Department == *event.Department
				fmt.Printf("[SSE-DEBUG] Condition Check: Dept != nil? %v | Values Match? %v\n", client.Department != nil, isDeptMatch)
			}

			if client.ClientID == clientID && client.Department != nil && *client.Department == *event.Department {
				fmt.Println("[SSE-DEBUG] >>> MATCHED! Attempting to send to channel...")
				select {
				case client.Chan <- data:
					fmt.Println("[SSE-DEBUG] >>> SUCCESS: Sent data to client channel.")
				default:
					fmt.Println("[SSE-DEBUG] >>> FAILED: Channel full/blocked. Removing client.")
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
		fmt.Println("[SSE-DEBUG] Processing logic for GUEST")
		for _, client := range h.Clients {
			if client.ClientID == clientID {
				fmt.Printf("[SSE-DEBUG] Found guest client. Type: %s vs EventType: %s\n", client.Type, event.Type)
			}

			if client.ClientID == clientID && client.Department == nil && client.Type == event.Type {
				fmt.Println("[SSE-DEBUG] >>> MATCHED GUEST! Attempting to send...")
				select {
				case client.Chan <- data:
					fmt.Println("[SSE-DEBUG] >>> SUCCESS: Sent data to guest.")
				default:
					fmt.Println("[SSE-DEBUG] >>> FAILED: Guest channel blocked.")
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
	} else {
		fmt.Println("[SSE-DEBUG] >>> SKIPPED: Event type or Department condition invalid.")
	}
	h.Mutex.RUnlock()
	fmt.Println("[SSE-DEBUG] <--- End SendToClient\n")
}
