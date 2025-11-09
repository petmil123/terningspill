package gameManager

import (
	"dice/socketHandler"
	"fmt"
	"math/rand"
)

type Player struct {
	ConnectionID string
	SendChan     chan socketHandler.PlayerMessage
	ReceiveChan  chan socketHandler.PlayerMessage
}

type turnPhase string

const (
	toThrow turnPhase = "toThrow"
	toCover turnPhase = "toCover"
	toCall  turnPhase = "toCall"
)

type gameState struct {
	player1Covered [6]int
	player2Covered [6]int
	turn           int // 1 or 2
	turnPhase      turnPhase
	chosenField    int // 1 to 6
	diceRoll       int // 1 to 6
}

func (gs *gameState) toString(playerNum int) string {
	if playerNum != 1 && playerNum != 2 {
		panic("Invalid player number")
	}
	if playerNum == 1 {
		return fmt.Sprintf("%v,%v,%t,%s,%d,%d",
			gs.player1Covered, gs.player2Covered, gs.turn == 1, gs.turnPhase, gs.chosenField, gs.diceRoll)
	}
	if playerNum == 2 {
		return fmt.Sprintf("%v,%v,%t,%s,%d,%d",
			gs.player2Covered, gs.player1Covered, gs.turn == 2, gs.turnPhase, gs.chosenField, gs.diceRoll)
	}
	return ""
}

func initializeGameState() gameState {
	return gameState{
		player1Covered: [6]int{0, 0, 0, 0, 0, 0},
		player2Covered: [6]int{0, 0, 0, 0, 0, 0},
		turn:           rand.Intn(2) + 1, // Randomly choose starting player
		turnPhase:      toThrow,
	}
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

	gameState := initializeGameState()
	fmt.Println(gameState.toString(1))
	fmt.Println(gameState.toString(2))
	player1.SendChan <- socketHandler.PlayerMessage{
		ConnectionID: player1.ConnectionID,
		Message:      socketHandler.Message{Action: "gameState", Data: gameState.toString(1)},
	}
	player2.SendChan <- socketHandler.PlayerMessage{
		ConnectionID: player2.ConnectionID,
		Message:      socketHandler.Message{Action: "gameState", Data: gameState.toString(2)},
	}
}
