package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ChatServer represents a single chat server instance
type ChatServer struct {
	ID        string
	Port      string
	Clients   map[*websocket.Conn]string
	Mutex     sync.Mutex
	Upgrader  websocket.Upgrader
	Broadcast chan string
}

// Create a new ChatServer instance
func NewChatServer(id string, port string) *ChatServer {
	return &ChatServer{
		ID:        id,
		Port:      port,
		Clients:   make(map[*websocket.Conn]string),
		Mutex:     sync.Mutex{},
		Upgrader:  websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		Broadcast: make(chan string),
	}
}

// Handle WebSocket connections for this server
func (s *ChatServer) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the incoming HTTP connection to a WebSocket connection
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Generate a temporary name for the user
	tempName := fmt.Sprintf("User-%d", len(s.Clients)+1)

	// Store the client with a temporary username
	s.Mutex.Lock()
	s.Clients[conn] = tempName
	s.Mutex.Unlock()

	// Send a prompt for username
	promptMsg := fmt.Sprintf("SYSTEM_MSG:Welcome to chat server %s! Please type your username to begin.", s.ID)
	conn.WriteMessage(websocket.TextMessage, []byte(promptMsg))

	// Listen for incoming messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			s.Mutex.Lock()
			username := s.Clients[conn]
			delete(s.Clients, conn)
			s.Mutex.Unlock()

			// Broadcast that a user has left
			s.Broadcast <- fmt.Sprintf("SYSTEM_MSG:%s has left the chat", username)
			break
		}

		messageStr := string(message)

		// Check if this is the first message (username)
		s.Mutex.Lock()
		currentUsername := s.Clients[conn]
		if currentUsername == tempName {
			// This is the username being set
			newUsername := messageStr
			s.Clients[conn] = newUsername
			s.Mutex.Unlock()

			// Send confirmation back to this client only
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("SYSTEM_MSG:Username set to %s. Welcome to the chat!", newUsername)))

			// Broadcast that a new user has joined
			s.Broadcast <- fmt.Sprintf("SYSTEM_MSG:%s has joined the chat", newUsername)

			log.Printf("[Server %s] User connected: %s\n", s.ID, newUsername)
		} else {
			s.Mutex.Unlock()
			// Regular message - broadcast with username
			s.Broadcast <- fmt.Sprintf("%s: %s", currentUsername, messageStr)
		}
	}
}

// Broadcast messages to all connected clients in this server
func (s *ChatServer) handleMessages() {
	for {
		// Get the next message from the broadcast channel
		message := <-s.Broadcast

		// Send the message to all connected clients
		s.Mutex.Lock()
		for client := range s.Clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println(err)
				client.Close()
				delete(s.Clients, client)
			}
		}
		s.Mutex.Unlock()
	}
}

// Start a chat server on the specified port
func (s *ChatServer) Start() {
	// Set up routes
	mux := http.NewServeMux()

	// WebSocket endpoint
	mux.HandleFunc("/ws", s.handleConnections)

	// Serve static files
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	// Start the message handling goroutine
	go s.handleMessages()

	// Start the server
	serverAddr := ":" + s.Port
	log.Printf("Starting chat server %s on %s\n", s.ID, serverAddr)

	// Start HTTP server with error handling
	err := http.ListenAndServe(serverAddr, mux)
	if err != nil {
		log.Printf("Error starting server %s: %v\n", s.ID, err)
	}
}

func main() {
	// Parse command-line flags for default server
	defaultPort := flag.String("port", "8080", "Default server port")
	flag.Parse()

	// Create and start the default server
	defaultServer := NewChatServer("main", *defaultPort)
	go defaultServer.Start()

	// Define additional servers
	servers := []*ChatServer{
		NewChatServer("chat1", "8081"),
		NewChatServer("chat2", "8082"),
	}

	// Start each additional server in its own goroutine
	for _, server := range servers {
		go server.Start()
	}

	// Wait indefinitely (you might want to add proper shutdown handling in a production app)
	select {}
}
