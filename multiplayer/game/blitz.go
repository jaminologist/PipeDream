package game

import (
	"encoding/json"
	"log"
	"time"

	"bryjamin.com/multiplayer/message"
	"github.com/gorilla/websocket"
)

type TimeLimit struct {
	Time time.Duration
}

type GameOver struct {
	Time time.Duration
}

type SinglePlayerBlitzGame struct {
	board     *Board
	timeLimit time.Duration
	isOver    bool
	score     int

	playerInputChannel   chan *message.BoardInput
	playerOutputChannel  chan *message.Message
	gameOverInputChannel chan bool
}

type SinglePlayerBlitzGameState struct {
	Board          *Board
	BoardReports   []BoardReport
	Score          int
	IsOver         bool
	DestroyedPipes []DestroyedPipe
}

func NewSinglePlayerBlitzGame(playerOutputChannel chan *message.Message, timeLimit time.Duration) *SinglePlayerBlitzGame {

	board := NewBoard(7, 7)

	return &SinglePlayerBlitzGame{
		timeLimit:            timeLimit,
		board:                &board,
		playerInputChannel:   make(chan *message.BoardInput),
		playerOutputChannel:  playerOutputChannel,
		gameOverInputChannel: make(chan bool),
	}

}

func (game *SinglePlayerBlitzGame) SendBoardInput(input *message.BoardInput) {
	game.playerInputChannel <- input
}

func (g *SinglePlayerBlitzGame) Run() {

	g.board.UpdateBoardPipeConnections()

	go func() {

		g.send(&SinglePlayerBlitzGameState{
			Board: g.board,
			Score: g.score,
		})

		for {
			g.timeLimit = g.timeLimit - serverTick
			g.send(&TimeLimit{
				Time: g.timeLimit,
			})
			g.isOver = g.timeLimit <= 0
			if g.isOver {
				g.gameOverInputChannel <- g.isOver
				break
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
				g.send(&gameState)
				break OuterLoop
			}
		case boardInput := <-g.playerInputChannel:
			g.board.RotatePipeClockwise(boardInput.X, boardInput.Y)
			boardReports := g.board.UpdateBoardPipeConnections()

			g.score += calculateScoreFromBoardReports(boardReports)

			gameState := SinglePlayerBlitzGameState{
				BoardReports: boardReports,
				Score:        g.score,
				IsOver:       g.isOver,
			}

			g.send(&gameState)
		}
	}

}

func (g *SinglePlayerBlitzGame) send(v interface{}) {
	messageBytes, err := json.Marshal(v)

	if err != nil {
		log.Println(err)
	} else {
		g.playerOutputChannel <- &message.Message{MessageType: websocket.TextMessage, Message: messageBytes}
	}
}

func calculateScoreFromBoardReports(boardReports []BoardReport) int {

	pipesDestroyed := 0
	for i := 0; i < len(boardReports); i++ {
		pipesDestroyed += len(boardReports[i].DestroyedPipes)
	}

	score := 1250 * pipesDestroyed

	return score
}
