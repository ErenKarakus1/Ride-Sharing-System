package notification

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Notification struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

type client struct {
	userID string
	conn   *websocket.Conn
}

type outboundMessage struct {
	userID  string
	payload []byte
}

type Hub struct {
	register       chan client
	unregister     chan client
	outbound       chan outboundMessage
	clients        map[string]map[*websocket.Conn]struct{}
	allowedOrigins []string
	upgrader       websocket.Upgrader
}

func NewHub(allowedOrigins []string) *Hub {
	hub := &Hub{
		register:       make(chan client),
		unregister:     make(chan client),
		outbound:       make(chan outboundMessage, 64),
		clients:        make(map[string]map[*websocket.Conn]struct{}),
		allowedOrigins: allowedOrigins,
	}
	hub.upgrader = websocket.Upgrader{
		CheckOrigin: hub.originAllowed,
	}

	return hub
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for _, conns := range h.clients {
				for conn := range conns {
					_ = conn.Close()
				}
			}
			return
		case client := <-h.register:
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*websocket.Conn]struct{})
			}
			h.clients[client.userID][client.conn] = struct{}{}
		case client := <-h.unregister:
			h.removeClient(client)
		case message := <-h.outbound:
			for conn := range h.clients[message.userID] {
				if err := conn.WriteMessage(websocket.TextMessage, message.payload); err != nil {
					h.removeClient(client{userID: message.userID, conn: conn})
				}
			}
		}
	}
}

func (h *Hub) Send(userID string, notification Notification) error {
	payload, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	h.outbound <- outboundMessage{
		userID:  userID,
		payload: payload,
	}
	return nil
}

func (h *Hub) HandleWebSocket(ctx *gin.Context, userID string) {
	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	current := client{userID: userID, conn: conn}
	h.register <- current
	go h.readLoop(current)
}

func (h *Hub) readLoop(current client) {
	defer func() {
		h.unregister <- current
	}()

	for {
		if _, _, err := current.conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (h *Hub) removeClient(client client) {
	conns := h.clients[client.userID]
	if conns == nil {
		return
	}

	if _, ok := conns[client.conn]; ok {
		delete(conns, client.conn)
		_ = client.conn.Close()
	}
	if len(conns) == 0 {
		delete(h.clients, client.userID)
	}
}

func (h *Hub) originAllowed(request *http.Request) bool {
	origin := request.Header.Get("Origin")
	return origin != "" && slices.Contains(h.allowedOrigins, origin)
}
