// Pairs players and initiates game between two connected players.

package gameInitiator

import (
	"dice/gameManager"
	"dice/socketHandler"
	"fmt"
)

type ConnectionRequest struct {
	Player   gameManager.Player
	AcceptCh chan bool
}

func GameInitiator(connectRequestChan chan ConnectionRequest) {
	players := make([]*gameManager.Player, 0, 2)
	for req := range connectRequestChan {
		fmt.Println("New connection request for player:", req.Player.ConnectionID)
		if len(players) >= 2 {
			// Reject connection
			req.Player.SendChan <- socketHandler.PlayerMessage{
				ConnectionID: req.Player.ConnectionID,
				Message:      socketHandler.Message{Action: "error", Data: "Maximum players reached"},
			}
			req.AcceptCh <- false
		} else {
			players = append(players, &req.Player)
			// Acknowledge connection
			req.AcceptCh <- true
			req.Player.SendChan <- socketHandler.PlayerMessage{
				ConnectionID: req.Player.ConnectionID,
				Message:      socketHandler.Message{Action: "connected", Data: "Welcome. You will be paired soon."},
			}
			if len(players) == 2 {
				fmt.Println("Two players connected. Starting game...")
				go gameManager.RunGame(players[0], players[1])
				// Reset players for next game
				players = make([]*gameManager.Player, 0, 2)
			}
		}
	}
}
