package server

import (
	"chat-server/server/actions"
	"chat-server/server/hub"
	"chat-server/server/models"
	"chat-server/server/player"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Ajustar según necesidad
}

var hubManager = hub.NewHub()

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar a WebSocket:", err)
		return
	}

	player := player.NewPlayer(conn)
	hubManager.AddPlayer(player)

	go handleMessages(player)
}

func handleMessages(player *player.Player) {
	defer func() {
		hubManager.RemovePlayer(player)
		hubManager.BroadcastUserLeave(player)
		player.Close()
	}()

	for {
		message, err := player.ReadMessage()
		if err != nil {
			log.Printf("Error al leer mensaje: %v", err)
			break
		}

		switch message.Action {
		case actions.Authenticate:
			err := authenticateClient(player, message)

			if err != nil {
				log.Printf("Error al deserializar datos de autenticación: %v", err)
				continue
			}
		case actions.NewMessage:
			err := newMessage(player, message)

			if err != nil {
				log.Printf("Error al deserializar datos del mensaje: %v", err)
				continue
			}
		}
	}
}

func checkToken(player *player.Player, data models.AuthData) bool {
	// Aquí se puede implementar la lógica de autenticación
	// Se puede usar una base de datos, un servicio de autenticación, etc.
	return data.Token == "token_secreto"
}

func authenticateClient(player *player.Player, request models.RequestData) error {
	var authData models.AuthData

	err := request.GetData(&authData)
	if err != nil {
		log.Printf("Error al deserializar datos de autenticación: %v", err)
		return err
	}

	if checkToken(player, authData) {
		player.IsAuthenticated = true
		player.Name = authData.Name

		player.WriteMessage(models.ResponseData{
			Action:  actions.Authenticate,
			Success: true,
		})

		hubManager.InitMessages(player)
		hubManager.BroadcastPlayersOnline(player)
		hubManager.BroadcastUserJoin(player)
	} else {
		player.WriteMessage(models.ResponseData{
			Action: actions.Authenticate,
			Error:  "authentication_failed",
		})
	}

	return nil
}

func newMessage(player *player.Player, request models.RequestData) error {
	var dataMessage models.NewMessageData

	err := request.GetData(&dataMessage)
	if err != nil {
		log.Printf("Error al deserializar datos del mensaje: %v", err)
		return err
	}

	if player.IsAuthenticated {
		hubManager.BroadcastMessage(dataMessage.Content, player)
		confirmCreatedMessage(player)
	}

	return nil
}

func confirmCreatedMessage(player *player.Player) {
	ackMessage := models.ResponseData{
		Action:  actions.NewMessage,
		Success: true,
	}
	player.WriteMessage(ackMessage)
}
