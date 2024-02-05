package main

import (
	"chat-server/server"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/chat", server.HandleConnections)
	log.Println("Servidor iniciado en :4500")
	err := http.ListenAndServe(":4500", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
