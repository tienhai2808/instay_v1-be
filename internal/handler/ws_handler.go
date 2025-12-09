package handler

import (
	"net/http"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/hub"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	hub *hub.WSHub
}

func NewWSHandler(hub *hub.WSHub) *WSHandler {
	return &WSHandler{hub}
}

func (h *WSHandler) ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}

	clientID := c.GetInt64("client_id")
	clientType := c.GetString("client_type")
	departmentID := c.GetInt64("department_id")
	if clientID == 0 && clientType == "" {
		c.Error(common.ErrForbidden)
		return
	}

	var departmentIDP *int64
	if departmentID == 0 {
		departmentIDP = nil
	} else {
		departmentIDP = &departmentID
	}

	client := hub.NewWSClient(h.hub, conn, clientID, clientType, departmentIDP)
	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
