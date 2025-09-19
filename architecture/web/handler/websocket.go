package handler

import (
	"log"
	"net/http"
	"real-time-forum/architecture/models"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	userID   int64
	username string
}

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var (
	upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients   = make(map[*websocket.Conn]*Client)
	clientsMu sync.Mutex
	broadcast = make(chan WebSocketMessage)
)

func (m *MainHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Verify session
	session, err := m.service.Session.GetByUuid(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Get user info
	user, err := m.service.User.GetByID(session.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Create client
	client := &Client{
		conn:     conn,
		userID:   user.ID,
		username: user.Nickname,
	}

	// Add client to map
	clientsMu.Lock()
	clients[conn] = client
	clientsMu.Unlock()

	log.Printf("User %s connected via WebSocket", user.Nickname)

	// Handle messages
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		switch msg.Type {
		case "private_message":
			m.handlePrivateMessage(client, msg)
		case "new_post":
			// Broadcast new post to all clients
			broadcast <- msg
		}
	}

	// Remove client when disconnected
	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
	log.Printf("User %s disconnected from WebSocket", user.Nickname)
}

func (m *MainHandler) handlePrivateMessage(sender *Client, msg WebSocketMessage) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return
	}

	recipientIDFloat, ok := payload["recipient_id"].(float64)
	if !ok {
		return
	}
	recipientID := int64(recipientIDFloat)

	content, ok := payload["content"].(string)
	if !ok {
		return
	}

	// Create message in database
	message := &models.PrivateMessage{
		SenderID:    sender.userID,
		RecipientID: recipientID,
		Content:     content,
		CreatedAt:   time.Now(),
	}

	_, err := m.service.PrivateMessage.Create(message)
	if err != nil {
		log.Printf("Failed to save private message: %v", err)
		return
	}

	// Find recipient client and send message
	clientsMu.Lock()
	for conn, client := range clients {
		if client.userID == recipientID {
			response := WebSocketMessage{
				Type: "private_message",
				Payload: map[string]interface{}{
					"sender_id":   sender.userID,
					"sender_name": sender.username,
					"content":     content,
					"created_at":  time.Now().Format(time.RFC3339),
				},
			}
			conn.WriteJSON(response)
		}
	}
	clientsMu.Unlock()
}

func HandleMessages() {
	for {
		msg := <-broadcast
		clientsMu.Lock()
		for conn := range clients {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("WebSocket error: %v", err)
				conn.Close()
				delete(clients, conn)
			}
		}
		clientsMu.Unlock()
	}
}
