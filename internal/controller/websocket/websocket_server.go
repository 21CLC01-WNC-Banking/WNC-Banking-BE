package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
)

type UserConnection struct {
	Conn *websocket.Conn
}

var ()

type Server struct {
	userCons     map[int]*UserConnection
	upgrader     *websocket.Upgrader
	userConsLock sync.RWMutex
}

func NewServer() *Server {
	userCons := make(map[int]*UserConnection)
	upgrader := websocket.Upgrader{} // Use default options
	return &Server{
		userCons: userCons,
		upgrader: &upgrader,
	}
}

func (s *Server) Run() {
	http.HandleFunc("/ws", s.handleWebSocket)
	fmt.Println("WebSocket server started at :3636")
	err := http.ListenAndServe(":3636", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Handle WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WebSocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}

	// Example: Get user ID from query params (use authentication in production)
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Missing userId"))
		conn.Close()
		return
	}

	userID, err := strconv.Atoi(userIDStr)

	// Add the connection to the map
	s.userConsLock.Lock()
	s.userCons[userID] = &UserConnection{Conn: conn}
	s.userConsLock.Unlock()

	// Notify user that the connection is established
	conn.WriteMessage(websocket.TextMessage, []byte("established"))
}

// Read messages from a user connection
func readMessages(userID int, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message from user %d: %v\n", userID, err)
			return
		}
		fmt.Printf("Received from user %d: %s\n", userID, message)
	}
}

// Send a message to a specific user
func (s *Server) SendToUser(userID int, message string) {
	s.userConsLock.RLock()
	defer s.userConsLock.RUnlock()

	if userConn, exists := s.userCons[userID]; exists {
		err := userConn.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Printf("Error sending message to user %d: %v\n", userID, err)
		}
	} else {
		fmt.Printf("User %d not connected\n", userID)
	}
}

// Broadcast a message to all users
func (s *Server) broadcastMessage(message string) {
	s.userConsLock.RLock()
	defer s.userConsLock.RUnlock()

	for userID, userConn := range s.userCons {
		err := userConn.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Printf("Error broadcasting message to user %d: %v\n", userID, err)
		}
	}
}
