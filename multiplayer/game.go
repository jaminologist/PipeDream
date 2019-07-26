package multiplayer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Game struct {
	lobby     *Lobby
	timeLimit time.Duration
}

type GameState struct {
	Time   time.Duration
	IsOver bool
}

func NewGame(l *Lobby, timeLimit time.Duration) *Game {

	return &Game{
		lobby:     l,
		timeLimit: timeLimit,
	}

}

func (g *Game) Run() {

	for {
		g.timeLimit = g.timeLimit - serverTick
		isOver := g.timeLimit <= 0
		messageBytes, err := json.Marshal(&GameState{Time: g.timeLimit, IsOver: isOver})

		if err != nil {
			log.Println(err)
		} else {
			g.lobby.boardcastAll <- &Message{messageType: websocket.TextMessage, message: messageBytes}
		}

		if isOver {
			break
		}

		time.Sleep(serverTick)
	}

}

type SinglePlayerBlitzGame struct {
	singlePlayerLobby *SinglePlayerLobby
	board             *Board
	timeLimit         time.Duration
	isOver            bool
	score             int

	playerInputChannel   chan *BoardInput
	gameOverInputChannel chan bool
}

type SinglePlayerBlitzGameState struct {
	Board          *Board
	BoardReports   []BoardReport
	Score          int
	IsOver         bool
	DestroyedPipes []DestroyedPipe
}

func NewSinglePlayerBlitzGame(spl *SinglePlayerLobby, timeLimit time.Duration) *SinglePlayerBlitzGame {

	board := NewBoard(7, 8)

	return &SinglePlayerBlitzGame{
		singlePlayerLobby:    spl,
		timeLimit:            timeLimit,
		board:                &board,
		playerInputChannel:   make(chan *BoardInput),
		gameOverInputChannel: make(chan bool),
	}

}

type TimeLimit struct {
	Time time.Duration
}

type GameOver struct {
	Time time.Duration
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
			g.board.Cells[boardInput.X][boardInput.Y].RotateClockWise()
			boardReports := g.board.UpdateBoardPipeConnections()

			pipesDestroyed := 0
			for i := 0; i < len(boardReports); i++ {
				pipesDestroyed += len(boardReports[i].DestroyedPipes)
			}

			g.score += 1250 * pipesDestroyed

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
		g.singlePlayerLobby.boardcastAll <- &Message{messageType: websocket.TextMessage, message: messageBytes}
	}

}
