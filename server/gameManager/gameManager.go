package gameManager

import (
	"dice/socketHandler"
	"fmt"
	"math/rand"
	"strconv"
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

	for {
		player1.SendChan <- socketHandler.PlayerMessage{
			ConnectionID: player1.ConnectionID,
			Message:      socketHandler.Message{Action: "gameState", Data: gameState.toString(1)},
		}
		player2.SendChan <- socketHandler.PlayerMessage{
			ConnectionID: player2.ConnectionID,
			Message:      socketHandler.Message{Action: "gameState", Data: gameState.toString(2)},
		}
		select {
		case msg := <-player1.ReceiveChan:
			fmt.Println("Received message from Player 1:", msg)
			if gameState.turn != 1 {
				player1.SendChan <- socketHandler.PlayerMessage{
					ConnectionID: player1.ConnectionID,
					Message:      socketHandler.Message{Action: "error", Data: "Not your turn."},
				}
			} else {
				switch gameState.turnPhase {
				case toThrow:
					if msg.Message.Action != "throwDie" {
						player1.SendChan <- socketHandler.PlayerMessage{
							ConnectionID: player1.ConnectionID,
							Message:      socketHandler.Message{Action: "error", Data: "Invalid action. You need to throw the die."},
						}
					}
					// Simulate die throw
					fmt.Println("Player 1 is throwing the die.")
					gameState.diceRoll = rand.Intn(6) + 1
					gameState.turnPhase = toCover
				case toCover:
					if msg.Message.Action != "coverField" {
						player1.SendChan <- socketHandler.PlayerMessage{
							ConnectionID: player1.ConnectionID,
							Message:      socketHandler.Message{Action: "error", Data: "Invalid action. You need to cover a field."},
						}
					}
					// Process covering field
					fieldToCover, err := strconv.Atoi(msg.Message.Data)
					if err != nil || fieldToCover < 1 || fieldToCover > 6 {
						player1.SendChan <- socketHandler.PlayerMessage{
							ConnectionID: player1.ConnectionID,
							Message:      socketHandler.Message{Action: "error", Data: "Invalid field number."},
						}
					} else {
						gameState.chosenField = fieldToCover
						if gameState.player1Covered[fieldToCover-1] < 2 {
							gameState.player1Covered[fieldToCover-1]++
						}
						gameState.turnPhase = toCall
						gameState.turn = 2 // Switch turn to player 2
					}
				case toCall:
					if msg.Message.Action != "guess" {
						player1.SendChan <- socketHandler.PlayerMessage{
							ConnectionID: player1.ConnectionID,
							Message:      socketHandler.Message{Action: "error", Data: "Invalid action. You need to make a guess."},
						}
					}
					if !(msg.Message.Data == "bluff" || msg.Message.Data == "truth") {
						player1.SendChan <- socketHandler.PlayerMessage{
							ConnectionID: player1.ConnectionID,
							Message:      socketHandler.Message{Action: "error", Data: "Invalid guess. Must be 'bluff' or 'truth'."},
						}
					}
					// Process guess
					if msg.Message.Data == "truth" {
						gameState.turnPhase = toThrow
						continue
					}
					isBluff := gameState.diceRoll != gameState.chosenField
					if isBluff {
						//Guess was correct, remove the chosen cover.
						gameState.player2Covered[gameState.chosenField-1]--
					} else {
						for {
							//Guess was wrong, remove a random cover from player 1.
							field := rand.Intn(6)
							if gameState.player1Covered[field] > 0 {
								gameState.player1Covered[field]--
								break
							}
						}
					}
					gameState.turnPhase = toThrow
				}
			}
		case msg := <-player2.ReceiveChan:
			fmt.Println("Received message from Player 2:", msg)
			if gameState.turn != 2 {
				player2.SendChan <- socketHandler.PlayerMessage{
					ConnectionID: player2.ConnectionID,
					Message:      socketHandler.Message{Action: "error", Data: "Not your turn."},
				}
			} else {
				switch gameState.turnPhase {
				case toThrow:
					if msg.Message.Action != "throwDie" {
						player2.SendChan <- socketHandler.PlayerMessage{
							ConnectionID: player2.ConnectionID,
							Message:      socketHandler.Message{Action: "error", Data: "Invalid action. You need to throw the die."},
						}
					}
					// Simulate die throw
					fmt.Println("Player 2 is throwing the die.")
					gameState.diceRoll = rand.Intn(6) + 1
					gameState.turnPhase = toCover
				case toCover:
					if msg.Message.Action != "coverField" {
						player2.SendChan <- socketHandler.PlayerMessage{
							ConnectionID: player2.ConnectionID,
							Message:      socketHandler.Message{Action: "error", Data: "Invalid action. You need to cover a field."},
						}
					}
					// Process covering field
					fieldToCover, err := strconv.Atoi(msg.Message.Data)
					if err != nil || fieldToCover < 1 || fieldToCover > 6 {
						player2.SendChan <- socketHandler.PlayerMessage{
							ConnectionID: player2.ConnectionID,
							Message:      socketHandler.Message{Action: "error", Data: "Invalid field number."},
						}
					} else {
						gameState.chosenField = fieldToCover
						gameState.turnPhase = toCall
						if gameState.player2Covered[fieldToCover-1] < 2 {
							gameState.player2Covered[fieldToCover-1]++
						}
						gameState.turn = 1 // Switch turn to player 1
					}
				case toCall:
					if msg.Message.Action != "guess" {
						player2.SendChan <- socketHandler.PlayerMessage{
							ConnectionID: player2.ConnectionID,
							Message:      socketHandler.Message{Action: "error", Data: "Invalid action. You need to make a guess."},
						}
					}
					if !(msg.Message.Data == "bluff" || msg.Message.Data == "truth") {
						player2.SendChan <- socketHandler.PlayerMessage{
							ConnectionID: player2.ConnectionID,
							Message:      socketHandler.Message{Action: "error", Data: "Invalid guess. Must be 'bluff' or 'truth'."},
						}
					}
					// Process guess
					if msg.Message.Data == "truth" {
						gameState.turnPhase = toThrow
						continue
					}
					isBluff := gameState.diceRoll != gameState.chosenField
					if isBluff {
						//Guess was correct, remove the chosen cover.
						gameState.player1Covered[gameState.chosenField-1]--
					} else {
						for {
							//Guess was wrong, remove a random cover from player 2.
							field := rand.Intn(6)
							if gameState.player2Covered[field] > 0 {
								gameState.player2Covered[field]--
								break
							}
						}
					}
					gameState.turnPhase = toThrow
				}
			}
		}
	}
}
