package models

type PlayerOnline struct {
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

type PlayerLeft struct {
	ID string `json:"id"`
}

type PlayerJoined struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"rank"`
}
