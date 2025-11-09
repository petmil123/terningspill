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

func GameManager(connectRequestChan chan ConnectionRequest) {
	players := make([]*gameManager.Player, 0, 2)
	for {
		select {
		case req := <-connectRequestChan:
			if len(players) >= 2 {
				// Reject connection
				fmt.Println("Rejected")
				req.Player.SendChan <- socketHandler.PlayerMessage{
					ConnectionID: req.Player.ConnectionID,
					Message:      socketHandler.Message{Action: "error", Data: "Maximum players reached"},
				}
				req.AcceptCh <- false
			} else {
				players = append(players, &req.Player)
				// Acknowledge connection
				fmt.Println("Accepted connection for player:", req.Player.ConnectionID)
				req.Player.SendChan <- socketHandler.PlayerMessage{
					ConnectionID: req.Player.ConnectionID,
					Message:      socketHandler.Message{Action: "connected", Data: "Welcome Player"},
				}
				req.AcceptCh <- true
				if len(players) == 2 {
					fmt.Println("Two players connected. Starting game...")
					go gameManager.RunGame(players[0], players[1])
					// Reset players for next game
					players = make([]*gameManager.Player, 0, 2)
				}
			}
		}
	}
}
