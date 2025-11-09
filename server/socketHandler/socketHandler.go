package socketHandler

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Message struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

type PlayerMessage struct {
	ConnectionID string
	Message      Message
}

func HandleIncomingMessages(conn *websocket.Conn, connectionID string, receiveChan chan PlayerMessage) {
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading JSON:", err)
			break
		}
		fmt.Printf("Received message from %s: %+v\n", connectionID, msg)
		receiveChan <- PlayerMessage{ConnectionID: connectionID, Message: msg}
	}
}

func HandleOutgoingMessages(conn *websocket.Conn, sendChan chan PlayerMessage) {
	for {
		msg := <-sendChan
		fmt.Println("Sending message to", msg.ConnectionID, ":", msg.Message)
		err := conn.WriteJSON(msg.Message)
		if err != nil {
			fmt.Println("Error writing JSON:", err)
			break
		}
	}
}
