
# WebSocket Chat Server and Client

## Overview

This project includes a WebSocket server and client for a chat application. The server handles connections, authentication, and broadcasting messages to connected and authenticated clients.

## Getting Started

### Prerequisites

- Go (Golang) installed on your system
- Postman for testing WebSocket connections (optional)

### Running the Server

1. Navigate to the server directory:

   ```sh
   cd path/to/server
   ```

2. Run the server:

   ```sh
   go run main.go
   ```

   The server will start and listen for incoming WebSocket connections.

### Running the Client

1. Navigate to the client directory:

   ```sh
   cd path/to/client
   ```

2. Run the client:

   ```sh
   go run client.go
   ```

   This will start a client instance that connects to the WebSocket server and starts sending messages.

## Connecting with Postman

You can also test the WebSocket server using Postman:

1. Open Postman and create a new request.
2. Set the request type to WebSocket.
3. Enter the WebSocket server URL (e.g., `ws://localhost:4500/chat`) and connect.
4. Once connected, you can authenticate by sending a message in the following format:

   ```json
   {
       "action": "authenticate",
       "data": {
           "token": "token_secreto"
       }
   }
   ```

   This will authenticate your session, and you will start receiving broadcasted messages.
