package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]string) // Connected clients with their usernames
var broadcast = make(chan string)              // Channel to broadcast messages

// Handle WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the incoming HTTP connection to a WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Ask for the user's username
	conn.WriteMessage(websocket.TextMessage, []byte("Please enter your username:"))
	_, username, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	// Store the client with their username
	clients[conn] = string(username)
	log.Printf("New client connected: %s\n", username)

	// Send a welcome message
	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Welcome %s!", username)))

	// Listen for incoming messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(clients, conn)
			break
		}
		// Broadcast the message with the username
		broadcast <- fmt.Sprintf("%s: %s", username, message)
	}
}

// Broadcast messages to all connected clients
func handleMessages() {
	for {
		// Get the next message from the broadcast channel
		message := <-broadcast
		// Send the message to all connected clients
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// Set up a route to handle WebSocket connections
	http.HandleFunc("/ws", handleConnections)

	// Start the message handling goroutine
	go handleMessages()

	// Serve the frontend HTML
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Start the HTTP server
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
