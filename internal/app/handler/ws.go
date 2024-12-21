package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/tranminhquanq/gomess/internal/app/usecase"
)

type WsAction string

const (
	ActionSubscribe     WsAction = "subscribe"
	ActionSendMessage   WsAction = "send_message"
	ActionUpdateProfile WsAction = "update_profile"
	ActionDisconnect    WsAction = "disconnect"
)

type WsMessage struct {
	Version    string      `json:"version"`    // Version of the protocol or API
	Action     WsAction    `json:"action"`     // Action type the message corresponds to
	Timestamp  int64       `json:"timestamp"`  // Client's timestamp for the message. Format: Unix timestamp in milliseconds
	Parameters interface{} `json:"parameters"` // Parameters for the action
}

type WsError struct {
	Code    int    `json:"code"`    // Error code for identifying the issue
	Message string `json:"message"` // Human-readable error message
	Details string `json:"details"` // Optional additional details about the error
}

type WsResponse struct {
	Version   string      `json:"version"`   // Version of the protocol or API
	Status    string      `json:"status"`    // Response status (e.g., "success", "error")
	Action    WsAction    `json:"action"`    // Action type the response corresponds to
	Timestamp int64       `json:"timestamp"` // Server's timestamp for the response
	Data      interface{} `json:"data"`      // Data payload for the response
	Error     *WsError    `json:"error"`     // Error details if status is "error"
}

type WsHandler struct {
	userUsecase *usecase.UserUsecase
	clients     map[*websocket.Conn]bool // Keep track of active clients
	broadcast   chan []byte              // Broadcast channel for messages
	upgrader    websocket.Upgrader
}

// NewWsHandler creates a new WebSocket handler
func NewWsHandler(userUsecase *usecase.UserUsecase) *WsHandler {
	return &WsHandler{
		userUsecase: userUsecase,
		clients:     make(map[*websocket.Conn]bool),
		broadcast:   make(chan []byte),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow requests from any origin
				return true
			},
		},
	}
}

// ServeWs handles WebSocket connections
func (h *WsHandler) ServeWs(w http.ResponseWriter, r *http.Request) error {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Error("Error upgrading to WebSocket")
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return err
	}

	// Register new client
	h.clients[conn] = true
	logrus.Info("New WebSocket connection opened")

	go h.handleConnection(conn)
	go h.handleBroadcast()

	return nil
}

// handleConnection reads messages from the WebSocket connection
func (h *WsHandler) handleConnection(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		delete(h.clients, conn)
		logrus.Info("WebSocket connection closed")
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.WithError(err).Error("Error reading message from WebSocket")
			break
		}

		logrus.Infof("Received message: %s", msg)

		// Process and validate the message, then add it to the broadcast channel
		h.broadcast <- msg
	}
}

// handleBroadcast sends messages to all connected clients
func (h *WsHandler) handleBroadcast() {
	for {
		msg := <-h.broadcast

		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logrus.WithError(err).Error("Error writing message to WebSocket")
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}
