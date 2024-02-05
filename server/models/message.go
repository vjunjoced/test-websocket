package models

import "time"

type NewMessageData struct {
	Content string `json:"content"`
}

type ResponseMessageData struct {
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}
