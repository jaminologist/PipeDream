package game

import (
	"errors"
	"time"

	"bryjamin.com/multiplayer/message"
)

type AIBlitzGame struct {
	*timer
	*boardInputProcessor
	moves []*Point

	aiInputChannel      chan bool
	playerOutputChannel chan *message.Message
}

func NewAIBlitzGame(playerOutputChannel chan *message.Message, timeLimit time.Duration) *AIBlitzGame {

	board := NewBoard(7, 7)

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

func (g *AIBlitzGame) getNextMove() (*Point, error) {
	if len(g.moves) <= 0 {
		g.moves, _ = BoardSolve(g.board)
	}
	var move *Point

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
