package player

import (
	"chat-server/server/models"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	conn            *websocket.Conn
	ID              string
	Name            string
	Rank            int
	IsAuthenticated bool
	m               sync.Mutex
}

func NewPlayer(conn *websocket.Conn) *Player {
	return &Player{conn: conn}
}

func (c *Player) ReadMessage() (models.RequestData, error) {
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return models.RequestData{}, err
	}

	var request models.RequestData
	err = json.Unmarshal(msg, &request)
	return request, err
}

func (c *Player) WriteMessage(response models.ResponseData) error {
	c.m.Lock()
	defer c.m.Unlock()

	responseJson, err := json.Marshal(response)

	if err != nil {
		log.Printf("Error serialize message: %v", err)
		return err
	}

	return c.conn.WriteMessage(websocket.TextMessage, responseJson)
}

func (c *Player) Close() error {
	return c.conn.Close()
}
