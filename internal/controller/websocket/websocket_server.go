package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
)

type DeviceConnection struct {
	Conn *websocket.Conn
}

var ()

type Server struct {
	deviceCons     map[int]*DeviceConnection
	upgrader       *websocket.Upgrader
	deviceConsLock sync.RWMutex
}

func NewServer() *Server {
	deviceCons := make(map[int]*DeviceConnection)
	upgrader := websocket.Upgrader{} // Use default options
	return &Server{
		deviceCons: deviceCons,
		upgrader:   &upgrader,
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

	// Example: Get device ID from query params (use authentication in production)
	deviceIDStr := r.URL.Query().Get("deviceId")
	if deviceIDStr == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Missing deviceId"))
		conn.Close()
		return
	}

	deviceID, err := strconv.Atoi(deviceIDStr)

	// Add the connection to the map
	s.deviceConsLock.Lock()
	s.deviceCons[deviceID] = &DeviceConnection{Conn: conn}
	s.deviceConsLock.Unlock()

	// Notify device that the connection is established
	conn.WriteMessage(websocket.TextMessage, []byte("established"))
}

// Read messages from a device connection
func readMessages(deviceID int, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message from device %d: %v\n", deviceID, err)
			return
		}
		fmt.Printf("Received from device %d: %s\n", deviceID, message)
	}
}

// Send a message to a specific device
func (s *Server) SendToDevice(deviceID int, message string) {
	s.deviceConsLock.RLock()
	defer s.deviceConsLock.RUnlock()

	if deviceConn, exists := s.deviceCons[deviceID]; exists {
		err := deviceConn.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Printf("Error sending message to device %d: %v\n", deviceID, err)
		}
	} else {
		fmt.Printf("Device %d not connected\n", deviceID)
	}
}

// Broadcast a message to all devices
func (s *Server) broadcastMessage(message string) {
	s.deviceConsLock.RLock()
	defer s.deviceConsLock.RUnlock()

	for deviceID, deviceConn := range s.deviceCons {
		err := deviceConn.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Printf("Error broadcasting message to device %d: %v\n", deviceID, err)
		}
	}
}
