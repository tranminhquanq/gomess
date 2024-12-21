package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/tranminhquanq/gomess/internal/app/usecase"
)

type WSHandler struct {
	userUsecase *usecase.UserUsecase
	clients     map[*websocket.Conn]bool // Keep track of active clients
	broadcast   chan []byte              // Broadcast channel for messages
	upgrader    websocket.Upgrader
}

func NewWSHandler(userUsecase *usecase.UserUsecase) *WSHandler {
	return &WSHandler{
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

func (h *WSHandler) HandleWS(w http.ResponseWriter, r *http.Request) error {
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

func (h *WSHandler) handleConnection(conn *websocket.Conn) {
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

func (h *WSHandler) handleBroadcast() {
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
