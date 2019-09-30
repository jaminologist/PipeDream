package player

import (
	"encoding/json"
	"fmt"
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

			/*var timelimit TimeLimit
			_ = json.Unmarshal(message, &timelimit)

			fmt.Println("Up top baby")

			if timelimit.Time != time.Duration(0) {
				fmt.Println("Returnin' nothing baby")
				break
			}*/

			if len(runner.moves) > 0 {
				var move *model.Point
				move, runner.moves = runner.moves[0], runner.moves[1:]

				bytes, _ := json.Marshal(move)

				runner.player.SendMessage(&PlayerMessage{
					MessageType: 100,
					Message:     bytes,
					Player:      runner.player,
				})
			} else {
				var state model.BlitzGameState
				err := json.Unmarshal(message, &state)

				if err != nil {
					break
				}

				if len(state.BoardReports) > 0 && len(runner.moves) <= 0 {
					runner.recentBlitzGameState = &state
					moves, _ := runner.getNextMoves()
					runner.moves = moves
				}

				if state.Board != nil && len(runner.moves) <= 0 {
					fmt.Println("Sword and board")
					moves, _ := model.BoardSolve(state.Board)
					runner.moves = moves
				}
			}
		}

	}

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

	/*var move *model.Point

	if len(g.moves) <= 0 {
		return nil, errors.New("No Moves Available")
	}

	pauseForAnimationTime := time.Duration(0)

	for _, boardReport := range blitzGameState.BoardReports {
		pauseForAnimationTime += boardReport.MaximumAnimationTime
	}

	move, g.moves = g.moves[0], g.moves[1:]

	return move, nil*/
}
