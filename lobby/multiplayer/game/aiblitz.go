package game

import (
	"errors"
	"time"

	"github.com/bryjammin/pipedream/lobby/multiplayer/game/model"
	"github.com/bryjammin/pipedream/lobby/multiplayer/message"
)

type AIBlitzGame struct {
	*timer
	*boardInputProcessor
	moves []*model.Point

	aiInputChannel      chan bool
	playerOutputChannel chan *message.Message
}

func NewAIBlitzGame(playerOutputChannel chan *message.Message, timeLimit time.Duration) *AIBlitzGame {

	board := model.NewBoard(7, 7)

	return &AIBlitzGame{
		playerOutputChannel: playerOutputChannel,
		aiInputChannel:      make(chan bool),
		boardInputProcessor: &boardInputProcessor{
			board:     &board,
			messageCh: playerOutputChannel,
		},
		timer: &timer{
			timeLimit:  timeLimit,
			messageCh:  playerOutputChannel,
			finishedCh: make(chan bool),
		},
	}

}

func (g *AIBlitzGame) getNextMove() (*model.Point, error) {
	if len(g.moves) <= 0 {
		g.moves, _ = model.BoardSolve(g.board)
	}
	var move *model.Point

	if len(g.moves) <= 0 {
		return nil, errors.New("No Moves Available")
	}

	move, g.moves = g.moves[0], g.moves[1:]

	return move, nil
}

func (g *AIBlitzGame) Run() {

	g.board.UpdateBoardPipeConnections()

	go func() {
		g.processGameBegin()
		g.aiInputChannel <- true
		g.countdown()
	}()

OuterLoop:
	for {
		select {
		case isOver := <-g.finishedCh:
			if isOver {
				g.boardInputProcessor.processGameOver()
				break OuterLoop
			}
		case _ = <-g.aiInputChannel:

			move, err := g.getNextMove()
			if err != nil {
				continue
			}

			blitzGameState := g.processBoardInput(move.X, move.Y)
			pauseForAnimationTime := time.Duration(0)

			for _, boardReport := range blitzGameState.BoardReports {
				pauseForAnimationTime += boardReport.MaximumAnimationTime
			}

			go func() {
				time.Sleep(pauseForAnimationTime)

				//To avoid sending too much data to the client this adds a pause before the next AI move.
				time.Sleep(serverTick / 8)
				g.aiInputChannel <- true
			}()
		}
	}

}
