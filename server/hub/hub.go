package hub

import (
	"chat-server/server/actions"
	"chat-server/server/models"
	"chat-server/server/player"
	"sync"
	"time"
)

type Message struct {
	Content string
	Created time.Time
}

type Hub struct {
	players  []*player.Player
	messages []Message
	lock     sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		players:  make([]*player.Player, 0),
		messages: make([]Message, 0),
	}
}

func (h *Hub) AddPlayer(player *player.Player) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.players = append(h.players, player)
}

func (h *Hub) InitMessages(player *player.Player) {
	h.lock.Lock()
	defer h.lock.Unlock()

	listMessage := []models.ResponseMessageData{}

	for _, msg := range h.messages {
		listMessage = append(listMessage, models.ResponseMessageData{Content: msg.Content, Created: msg.Created})
	}

	response := models.ResponseData{Action: actions.ListMessages}
	response.SetData(listMessage)

	player.WriteMessage(response)
}

func (h *Hub) RemovePlayer(player *player.Player) {
	h.lock.Lock()
	defer h.lock.Unlock()
	for i, c := range h.players {
		if c == player {
			h.players = append(h.players[:i], h.players[i+1:]...)
			break
		}
	}
}

func (h *Hub) BroadcastMessage(message string, sender *player.Player) {
	h.lock.Lock()
	defer h.lock.Unlock()

	formattedMessage := sender.Name + ": " + message
	currentTime := time.Now()
	h.messages = append(h.messages, Message{Content: formattedMessage, Created: currentTime})
	if len(h.messages) > 100 {
		h.messages = h.messages[1:]
	}

	newMessage := models.ResponseData{Action: actions.NewMessage}
	data := models.ResponseMessageData{Content: formattedMessage, Created: currentTime}
	newMessage.SetData(data)

	h.sendDataToAuthenticatedPlayers(newMessage, sender)
}

func (h *Hub) sendDataToAuthenticatedPlayers(response models.ResponseData, sender *player.Player) {
	for _, p := range h.players {

		if p != sender && p.IsAuthenticated {
			p.WriteMessage(response)
		}
	}
}

func (h *Hub) BroadcastPlayersOnline(player *player.Player) {
	h.lock.Lock()
	defer h.lock.Unlock()

	var players []models.PlayerOnline
	for _, p := range h.players {
		if p.IsAuthenticated {
			players = append(players, models.PlayerOnline{Name: p.Name, Rank: p.Rank})
		}
	}

	response := models.ResponseData{Action: actions.PlayersOnline}
	response.SetData(players)

	player.WriteMessage(response)
}

func (h *Hub) BroadcastUserJoin(player *player.Player) {
	h.lock.Lock()
	defer h.lock.Unlock()

	data := models.PlayerJoined{ID: player.ID, Name: player.Name, Rank: player.Rank}
	response := models.ResponseData{Action: "user_joined"}
	response.SetData(data)

	h.sendDataToAuthenticatedPlayers(response, player)
}

func (h *Hub) BroadcastUserLeave(player *player.Player) {
	h.lock.Lock()
	defer h.lock.Unlock()

	data := models.PlayerLeft{ID: player.ID}
	response := models.ResponseData{Action: actions.PlayerLeft}
	response.SetData(data)

	h.sendDataToAuthenticatedPlayers(response, player)
}
