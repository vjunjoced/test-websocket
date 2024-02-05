package main

import (
	"chat-server/server/actions"
	"chat-server/server/models"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	for i := 0; i < 5; i++ {
		go simulateUser(i)
	}
	select {} // Evita que el programa principal termine
}

func simulateUser(userID int) {
	// Conectarse al servidor WebSocket
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:4500/chat", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	jsonString := `{
    "token": "token_secreto",
    "name": "User ` + fmt.Sprintf("%d", userID) + `"
}`
	authMessage := models.RequestData{
		Action: actions.Authenticate,
		Data:   json.RawMessage(jsonString),
	}
	authMsg, error := json.Marshal(authMessage)
	if error != nil {
		log.Println("Error al serializar mensaje de autenticación:", error)
		return
	}

	log.Println("Sending auth message:", string(authMsg))
	c.WriteMessage(websocket.TextMessage, authMsg)

	// Enviar mensajes de manera aleatoria
	go func() {
		index := 0
		for {
			time.Sleep(time.Duration(rand.Intn(7)+16) * time.Second) // Espera 3-6 segundos
			randomMessage := models.RequestData{
				Action: actions.NewMessage,
				Data:   json.RawMessage(`{"content":"Hello, world! ` + fmt.Sprintf("%d", index) + `"}`),
			}
			sendMessage(c, randomMessage)

			index++
		}
	}()

	// Desconectarse y reconectarse
	go func() {
		time.Sleep(time.Duration(rand.Intn(15)+30) * time.Second) // Espera 10-16 segundos
		c.Close()
		time.Sleep(time.Duration(rand.Intn(10)+50) * time.Second) // Espera 6-10 segundos para reconectar
		simulateUser(userID)                                      // Se reconecta
	}()
	select {}
}

func sendMessage(c *websocket.Conn, msg interface{}) {
	msgData, err := json.Marshal(msg)
	if err != nil {
		log.Println("error marshalling message:", err)
		return
	}
	if err := c.WriteMessage(websocket.TextMessage, msgData); err != nil {
		log.Println("write:", err)
		return
	}
}
