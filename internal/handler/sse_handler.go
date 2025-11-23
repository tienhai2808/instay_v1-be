package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/hub"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SSEHandler struct {
	hub *hub.SSEHub
}

func NewSSEHandler(hub *hub.SSEHub) *SSEHandler {
	return &SSEHandler{hub}
}

func (h *SSEHandler) HandleSSE(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Credentials", "true")
	origin := c.GetHeader("Origin")
	if origin != "" {
		c.Header("Access-Control-Allow-Origin", origin)
	}

	clientID := c.GetInt64("client_id")
	clientType := c.GetString("client_type")
	department := c.GetString("department")
	if clientID == 0 && clientType == "" {
		common.ToAPIResponse(c, http.StatusForbidden, common.ErrForbidden.Error(), nil)
		return
	}

	var departmentP *string
	if department == "" {
		departmentP = nil
	} else {
		departmentP = &department
	}

	client := &hub.SSEClient{
		ID:         uuid.NewString(),
		ClientID:   clientID,
		Type:       clientType,
		Department: departmentP,
		Chan:       make(chan []byte),
		Done:       make(chan bool),
	}

	h.hub.Register <- client
	defer func() {
		h.hub.Unregister <- client
	}()

	clientGone := c.Request.Context().Done()

	ticker := time.NewTicker(54 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message := <-client.Chan:
			var msg types.SSEEventData
			if err := json.Unmarshal(message, &msg); err != nil {
				sse.Encode(c.Writer, sse.Event{
					Event: "error",
					Data:  gin.H{"message": fmt.Sprintf("%v", err)},
				})
			} else {
				sse.Encode(c.Writer, sse.Event{
					Event: msg.Event,
					Data:  msg.Data,
				})
			}

			c.Writer.Flush()

		case <-client.Done:
			return

		case <-clientGone:
			return

		case <-ticker.C:
			c.Writer.Write([]byte(": keep-alive\n\n"))
			c.Writer.Flush()
		}
	}
}
