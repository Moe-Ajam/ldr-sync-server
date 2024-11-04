package main

import (
	"fmt"
	"net/http"

	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
)

func (cfg *apiConfig) handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}
		fmt.Printf("message recieved: %s\n", msg)
		broadcast <- msg
	}
}
