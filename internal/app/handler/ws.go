package handler

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/tranminhquanq/gomess/internal/app/domain"
	"github.com/tranminhquanq/gomess/internal/app/usecase"
	"github.com/tranminhquanq/gomess/internal/config"
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

type WsClient struct {
	ID   string
	Conn *websocket.Conn
	User domain.User
}

type WsHandler struct {
	userUsecase  *usecase.UserUsecase
	serverId     string
	localClients sync.Map // Stores active connections on this server
	upgrader     websocket.Upgrader
	// redisClient *redis.Client
}

// NewWsHandler creates a new WebSocket handler
func NewWsHandler(
	globalConfig *config.GlobalConfiguration,
	userUsecase *usecase.UserUsecase) *WsHandler {
	return &WsHandler{
		userUsecase:  userUsecase,
		serverId:     globalConfig.API.ID,
		localClients: sync.Map{},
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

	user, err := h.userUsecase.GetUserFromToken(r.Context())
	if err != nil {
		logrus.WithError(err).Error("Error getting user from token")
		http.Error(w, "Could not authenticate user", http.StatusUnauthorized)
		return err
	}

	client := &WsClient{ID: user.ID, Conn: conn, User: user}
	h.localClients.Store(client.ID, client)
	// h.registerClientInRedis(client.Id, h.serverID)

	go h.HandleIncomingMessages(client)

	return nil
}

func (h *WsHandler) handleCloseConnection(client *WsClient) {
	client.Conn.Close()
	h.localClients.Delete(client.ID)

	// Unregister client from Redis
	// h.unregisterClientFromRedis(client.ID)
}

func (h *WsHandler) HandleIncomingMessages(client *WsClient) {
	defer func() {
		h.handleCloseConnection(client)
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			logrus.WithError(err).Error("Error reading message from WebSocket")
			return
		}
		// TODO: You can parse the message and take appropriate actions based on the action type.
		// For now, i just broadcast the message to all local clients.
		go h.Broadcast2AllLocalClients(client.ID, msg)
	}
}

func (h *WsHandler) Broadcast2AllLocalClients(clientId string, message []byte) {
	h.localClients.Range(func(key, value interface{}) bool {
		client, ok := value.(*WsClient)
		if !ok || client == nil {
			logrus.Error(("Invalid client in localClients map"))
			return true // Continue to the next client
		}

		if client.ID != clientId { // Avoid sending to the sender
			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logrus.WithError(err).Error("Error writing message to WebSocket")
				return false // Stop iterating over clients
			}
		}
		return true
	})

	// Publish message to Redis
	// publishToRedis(clientId, message)
}

// func (h *WsHandler) broadcastMessage2SpecificClients(clientIds []string, message []byte) {
// 	for _, clientId := range clientIds {
// 		client, ok := h.localClients.Load(clientId)
// 		if !ok {
// 			logrus.Errorf("Client not found in localClients map: %s", clientId)
// 			continue
// 		}

// 		c, ok := client.(*WsClient)
// 		if !ok || c == nil {
// 			logrus.Error(("Invalid client in localClients map"))
// 			continue
// 		}

// 		err := c.Conn.WriteMessage(websocket.TextMessage, message)
// 		if err != nil {
// 			logrus.WithError(err).Error("Error writing message to WebSocket")
// 			continue
// 		}
// 	}
// }

// Register client in Redis
// func (h *WsHandler) registerClientInRedis(clientID, serverID string) {
// 	h.redisClient.HSet(ctx, "ws_clients", clientID, serverID)
// }

// Unregister client from Redis
// func (h *WsHandler) unregisterClientFromRedis(clientID string) {
// 	h.redisClient.HDel(ctx, "ws_clients", clientID)
// }

// Publish message to Redis
// func (h *WsHandler) publishToRedis(senderID string, message []byte) {
// 	err := redisClient.Publish(ctx, redisChannel, string(message)).Err()
// 	if err != nil {
// 		log.Println("Failed to publish message:", err)
// 	}
// }

// // Subscribe to Redis channel
// func (h *WsHandler) subscribeToRedis() {
// 	pubsub := redisClient.Subscribe(ctx, redisChannel)
// 	defer pubsub.Close()

// 	for {
// 		msg, err := pubsub.ReceiveMessage(ctx)
// 		if err != nil {
// 			log.Println("Failed to receive message:", err)
// 			continue
// 		}

// 		// Broadcast received message locally
// 		broadcastMessageLocally([]byte(msg.Payload))
// 	}
// }
