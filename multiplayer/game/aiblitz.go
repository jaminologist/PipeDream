package game

import (
	"log"
	"time"

	"bryjamin.com/multiplayer/message"
	"bryjamin.com/multiplayer/send"
)

type AIBlitzGame struct {
	board     *Board
	timeLimit time.Duration
	isOver    bool
	score     int

	moves []*Point

	aiInputChannel       chan bool
	playerOutputChannel  chan *message.Message
	gameOverInputChannel chan bool
}

type AIBlitzGameState struct {
	Board          *Board
	BoardReports   []BoardReport
	Score          int
	IsOver         bool
	DestroyedPipes []DestroyedPipe
}

func NewAIBlitzGame(playerOutputChannel chan *message.Message, timeLimit time.Duration) *AIBlitzGame {

	board := NewBoard(7, 7)

	log.Println("New AI GAME CREATED")

	return &AIBlitzGame{
		timeLimit:            timeLimit,
		board:                &board,
		playerOutputChannel:  playerOutputChannel,
		gameOverInputChannel: make(chan bool),
		aiInputChannel:       make(chan bool),
	}

}

func (g *AIBlitzGame) Run() {

	g.board.UpdateBoardPipeConnections()

	go func() {

		send.SendMessageToAll(&SinglePlayerBlitzGameState{
			Board: g.board,
			Score: g.score,
		}, g.playerOutputChannel)

		g.aiInputChannel <- true

		for !g.isOver {
			g.timeLimit = g.timeLimit - serverTick

			send.SendMessageToAll(&TimeLimit{
				Time: g.timeLimit,
			}, g.playerOutputChannel)

			g.isOver = g.timeLimit <= 0
			if g.isOver {
				g.gameOverInputChannel <- g.isOver
			}

			time.Sleep(serverTick)
		}
	}()

OuterLoop:
	for {
		select {
		case isOver := <-g.gameOverInputChannel:
			if isOver {
				gameState := SinglePlayerBlitzGameState{
					Score:  g.score,
					IsOver: g.isOver,
				}
				send.SendMessageToAll(&gameState, g.playerOutputChannel)
				break OuterLoop
			}
		case _ = <-g.aiInputChannel:

			if len(g.moves) <= 0 {
				g.moves, _ = BoardSolve(g.board)
			}
			var move *Point

			if len(g.moves) > 0 {
				move, g.moves = g.moves[0], g.moves[1:]
				g.board.RotatePipeClockwise(move.X, move.Y)
			} else {
				log.Println("No Moves Available.")
			}
			boardReports := g.board.UpdateBoardPipeConnections()

			g.score += calculateScoreFromBoardReports(boardReports)

			gameState := SinglePlayerBlitzGameState{
				BoardReports: boardReports,
				Score:        g.score,
				IsOver:       g.isOver,
			}

			var pauseTime time.Duration

			for _, boardReport := range boardReports {
				pauseTime += boardReport.MaximumAnimationTime
			}

			send.SendMessageToAll(&gameState, g.playerOutputChannel)

			go func() {
				//Adds a pause to account for the
				time.Sleep(pauseTime)

				//To avoid sending too much data to the client this adds a pause before the next AI move.
				time.Sleep(serverTick / 8)
				g.aiInputChannel <- true
			}()
		}
	}

}
