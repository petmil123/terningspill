package main

import (
	"dice/gameInitiator"
	"dice/gameManager"
	"dice/socketHandler"

	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// This runs the function handleRequest when called.
	connectRequestChan := make(chan gameInitiator.ConnectionRequest, 4)
	go gameInitiator.GameInitiator(connectRequestChan)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//receive a connect request
		fmt.Println("Received request from client")
		//Upgrade to websocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Failed to upgrade to websocket:", err)
			return
		}
		// Add player to the list
		connectionID := uuid.New().String()
		player := gameManager.Player{
			ConnectionID: connectionID,
			SendChan:     make(chan socketHandler.PlayerMessage, 4),
			ReceiveChan:  make(chan socketHandler.PlayerMessage, 4),
		}
		acceptCh := make(chan bool)
		connectRequestChan <- gameInitiator.ConnectionRequest{
			Player:   player,
			AcceptCh: acceptCh,
		}
		accepted := <-acceptCh
		if !accepted {
			fmt.Println("Connection rejected for player:", connectionID)
			conn.Close()
			return
		}
		fmt.Println("Connection accepted for player:", connectionID)
		go socketHandler.HandleIncomingMessages(conn, connectionID, player.ReceiveChan)
		go socketHandler.HandleOutgoingMessages(conn, player.SendChan)

	})
	fmt.Println("Listening to port 8080")
	http.ListenAndServe(":8080", nil)
}
