package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 2048

	eventMarkRead = "mark_read"

	eventNewMessage = "new_message"

	eventError = "error"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type WSClient struct {
	Hub         *WSHub
	Conn        *websocket.Conn
	Send        chan []byte
	ID          string
	ClientID    int64
	StaffData   *types.StaffData
	Type        string
	ActiveChats map[int64]bool
}

func NewWSClient(hub *WSHub, conn *websocket.Conn, clientID int64, clientType string, staffData *types.StaffData) *WSClient {
	return &WSClient{
		hub,
		conn,
		make(chan []byte, 256),
		uuid.NewString(),
		clientID,
		staffData,
		clientType,
		make(map[int64]bool),
	}
}

type WSHub struct {
	Clients     map[string]map[string]*WSClient
	Register    chan *WSClient
	Unregister  chan *WSClient
	SendMessage chan *MessagePayload
	ChatSvc     service.ChatService
}

type MessagePayload struct {
	TargetKey string
	Data      []byte
}

func NewWSHub(chatSvc service.ChatService) *WSHub {
	return &WSHub{
		make(map[string]map[string]*WSClient),
		make(chan *WSClient),
		make(chan *WSClient),
		make(chan *MessagePayload),
		chatSvc,
	}
}

func (c *WSClient) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var req types.WSRequest
		if err := json.Unmarshal(data, &req); err != nil {
			c.sendError("invalid message")
			continue
		}

		switch req.Event {
		case "send_message":
			c.handleSendMessage(req.Data)
		case "mark_read":
			c.handleMarkRead(req.Data)
		default:
			c.sendError("unknown action")
		}
	}
}

func (c *WSClient) handleSendMessage(content []byte) {
	var req types.CreateMessageRequest
	if err := json.Unmarshal(content, &req); err != nil {
		c.sendError("invalid message")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message, err := c.Hub.ChatSvc.CreateMessage(ctx, req.ChatID, c.ClientID, c.Type, req)
	if err != nil {
		c.sendError("send message failed")
		return
	}

	res := types.WSResponse{
		Event: eventNewMessage,
		Data:  common.ToMessageResponse(message),
	}

	resBytes, _ := json.Marshal(res)

	targets := []string{
		"staff",
		fmt.Sprintf("guest_%d", message.Chat.OrderRoomID),
	}

	for _, t := range targets {
		msgPayload := &MessagePayload{
			TargetKey: t,
			Data:      resBytes,
		}

		c.Hub.SendMessage <- msgPayload
	}
}

func (c *WSClient) handleMarkRead(content []byte) {
	var req types.UpdateReadMessagesRequest
	if err := json.Unmarshal(content, &req); err != nil {
		c.sendError("invalid message")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chat, err := c.Hub.ChatSvc.UpdateReadMessages(ctx, req.ChatID, c.ClientID, c.Type)
	if err != nil {
		c.sendError("read message failed")
		return
	}

	res := types.WSResponse{
		Event: eventMarkRead,
		Data: types.UpdateReadMessagesResponse{
			ChatID:     req.ChatID,
			ReaderType: c.Type,
			ReadAt:     *chat.LastMessageAt,
			Reader:     (*types.BasicUserResponse)(c.StaffData),
		},
	}

	resBytes, _ := json.Marshal(res)
	targets := []string{
		"staff",
		fmt.Sprintf("guest_%d", chat.OrderRoomID),
	}

	for _, t := range targets {
		msgPayload := &MessagePayload{
			TargetKey: t,
			Data:      resBytes,
		}

		c.Hub.SendMessage <- msgPayload
	}
}

func (c *WSClient) sendError(msg string) {
	resp := types.WSResponse{
		Event: eventError,
		Data:  map[string]string{"message": msg},
	}
	b, _ := json.Marshal(resp)
	c.Send <- b
}

func (c *WSClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for range n {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.Register:
			key := client.getKey()
			if _, ok := h.Clients[key]; !ok {
				h.Clients[key] = make(map[string]*WSClient)
			}
			h.Clients[key][client.ID] = client

		case client := <-h.Unregister:
			key := client.getKey()
			if conns, ok := h.Clients[key]; ok {
				if _, exists := conns[client.ID]; exists {
					delete(conns, client.ID)
					close(client.Send)
					if len(conns) == 0 {
						delete(h.Clients, key)
					}
				}
			}

		case msg := <-h.SendMessage:
			if conns, ok := h.Clients[msg.TargetKey]; ok {
				for _, client := range conns {
					select {
					case client.Send <- msg.Data:
					default:
						close(client.Send)
						delete(conns, client.ID)
					}
				}

				if len(conns) == 0 {
					delete(h.Clients, msg.TargetKey)
				}
			}
		}
	}
}

func (c *WSClient) getKey() string {
	if c.Type == "guest" {
		return fmt.Sprintf("guest_%d", c.ClientID)
	}

	return "staff"
}
