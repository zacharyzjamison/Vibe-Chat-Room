# Multi-Server WebSocket Chat

A flexible WebSocket-based chat application that allows users to connect to multiple chat servers running on different ports.

## Overview

This application demonstrates a distributed chat system where multiple chat servers can run simultaneously on different ports. Users can connect to any available server through a responsive web interface and communicate with other users on the same server.

## Features

- Multiple chat servers running simultaneously on different ports
- Easy server switching through the UI
- Custom server creation via the UI
- Persistent username across messages within a server
- Real-time message broadcasting
- Connection status indicators
- System messages for user join/leave events
- Responsive web design

## Technical Stack

- **Backend**: Go (Golang)
- **Frontend**: HTML, CSS, JavaScript
- **WebSockets**: gorilla/websocket library

## Project Structure

```
.
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── public/             # Static files directory
│   └── index.html      # Web UI for the chat application
└── server.go           # Main server code
```

## How It Works

### Server Side

The application uses a `ChatServer` struct to manage each chat server instance. The main features include:

- WebSocket connection handling
- Client tracking with usernames
- Message broadcasting to all connected clients
- Support for system messages
- Dynamic server creation

### Client Side

The frontend provides a user-friendly interface with:

- Server selection buttons
- Custom server connection option
- Real-time message display
- Connection status indicators
- Username prompt on first connection

## Running the Application

### Prerequisites

- Go 1.23.5 or later
- GitHub.com/gorilla/websocket v1.5.3

### Starting the Default Servers

```bash
go run server.go
```

This starts:
- Main server on port 8080
- Chat server 1 on port 8081
- Chat server 2 on port 8082

### Starting a Custom Server

```bash
go run server.go custom <port>
```

Replace `<port>` with your desired port number (e.g., 8083).

## Usage

1. Open `http://localhost:8080` in your web browser
2. Select a server to connect to
3. Enter your username when prompted
4. Start chatting!

## Custom Server Creation

Users can create and connect to custom servers in two ways:

1. **Via UI**: Enter a port number and connect through the web interface
2. **Via Command Line**: Start a custom server from terminal as shown above

After successful connection to a custom server, a new button is added to the server selection panel for easy reconnection.

## Server Communication

- Each server operates independently
- Messages are only broadcast to users connected to the same server
- Connection status and user presence are tracked per server

## Future Enhancements

Potential improvements for the application:

- User authentication
- Message persistence
- Cross-server communication
- Private messaging
- File sharing capabilities
- Enhanced UI with themes and customization options