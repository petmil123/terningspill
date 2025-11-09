package gameManager

import (
	"dice/socketHandler"
)

type Player struct {
	ConnectionID string
	SendChan     chan socketHandler.PlayerMessage
	ReceiveChan  chan socketHandler.PlayerMessage
}

func RunGame(player1 *Player, player2 *Player) {
	player1.SendChan <- socketHandler.PlayerMessage{
		ConnectionID: player1.ConnectionID,
		Message:      socketHandler.Message{Action: "start", Data: "Game started! You are Player 1."},
	}
	player2.SendChan <- socketHandler.PlayerMessage{
		ConnectionID: player2.ConnectionID,
		Message:      socketHandler.Message{Action: "start", Data: "Game started! You are Player 2."},
	}
}
