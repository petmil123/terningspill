package gameManager

import (
	"dice/socketHandler"
	"fmt"
)

type Player struct {
	ConnectionID string
	SendChan     chan socketHandler.PlayerMessage
	ReceiveChan  chan socketHandler.PlayerMessage
}

func RunGame(player1 *Player, player2 *Player) {
	fmt.Println("Game started between", player1.ConnectionID, "and", player2.ConnectionID)
}
