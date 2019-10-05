package player

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"bryjamin.com/multiplayer/game/model"
)

//AIPlayer Used to mock a player and fill spaces for waiting players
type AIBlitzPlayer struct {
	PlayerMessageReceiver
}

//NewAIBlitzPlayer Returns a new AI Player that can play the blitz mode
func NewAIBlitzPlayer() *Player {

	messageFromServerCh := make(chan []byte, 0)

	player := &Player{}

	player.Conn = &AIBlitzPlayerConnection{
		Player:              player,
		messageFromServerCh: messageFromServerCh,
	}

	player.PlayerRunner = &AIBlitzPlayerRunner{
		messageFromServerCh: messageFromServerCh,
		player:              player,
	}

	return player
}

type AIBlitzPlayerRunner struct {
	player               *Player
	messageFromServerCh  chan []byte
	moves                []*model.Point
	recentBlitzGameState *model.BlitzGameState
}

func (runner *AIBlitzPlayerRunner) Run() {
	fmt.Println("Running ai blitz")
	for {

		select {
		case message := <-runner.messageFromServerCh:
			_ = message

			var gameState model.BlitzGameState
			err := json.Unmarshal(message, &gameState)

			if err != nil {
				log.Printf("Error: %+v", err)
				break
			}

			if gameState.TimeLimit != nil {
				log.Printf("Sent Time Limit")
				break
			}

			var enemyGameState model.VersusPlayerBlitzGamePlayerInformationSentToPlayers
			err = json.Unmarshal(message, &enemyGameState)

			if err != nil {
				log.Printf("%+v", err)
				break
			}

			if enemyGameState.EnemyInformation != nil {
				log.Printf("Sent Enemy Information")
				break
			}

			for _, report := range gameState.BoardReports {
				time.Sleep(report.MaximumAnimationTime)
			}

			if len(runner.moves) > 0 {
				time.Sleep(time.Duration(600) * time.Millisecond / 4)
				runner.sendNextMoveToLobby()
			} else {
				log.Println("Is it ever down here?")
				if gameState.IsOver {
					return
				}
				if len(gameState.BoardReports) > 0 {
					log.Printf("Sent Board Report Information")
					runner.recentBlitzGameState = &gameState
					moves, _ := runner.getNextMoves()
					runner.moves = moves
					runner.sendNextMoveToLobby()
				}
			}
		}

	}

}

func (runner *AIBlitzPlayerRunner) sendNextMoveToLobby() {
	var move *model.Point
	move, runner.moves = runner.moves[0], runner.moves[1:]

	bytes, _ := json.Marshal(move)

	runner.player.SendMessage(&PlayerMessage{
		MessageType: 100,
		Message:     bytes,
		Player:      runner.player,
	})

	time.Sleep(time.Millisecond * 20)
}

type AIBlitzPlayerConnection struct {
	*Player
	messageFromServerCh chan []byte
}

type TimeLimit struct {
	Time time.Duration
}

func (p *AIBlitzPlayerConnection) WriteMessage(messageType int, data []byte) error {

	go func() {
		p.messageFromServerCh <- data
	}()
	return nil
}

func (p *AIBlitzPlayerRunner) getNextMoves() ([]*model.Point, error) {

	blitzGameState := p.recentBlitzGameState
	board := blitzGameState.BoardReports[len(blitzGameState.BoardReports)-1].Board
	moves, err := model.BoardSolve(board)

	if err != nil {
		return []*model.Point{}, err
	}

	return moves, nil
}
