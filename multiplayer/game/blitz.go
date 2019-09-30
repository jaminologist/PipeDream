package game

import (
	"time"

	"bryjamin.com/multiplayer/game/model"
	"bryjamin.com/multiplayer/message"
	"bryjamin.com/multiplayer/send"
)

type blitzGameMode interface {
	countdown()
	processBoardInput(x int, y int)
}

type timer struct {
	timeLimit  time.Duration
	messageCh  chan *message.Message
	finishedCh chan bool
}

func (cd *timer) countdown() {
	for {
		cd.timeLimit = cd.timeLimit - serverTick

		send.SendMessageToAll(&TimeLimit{
			Time: cd.timeLimit,
		}, cd.messageCh)

		if cd.timeLimit <= 0 {
			cd.finishedCh <- true
		}

		time.Sleep(serverTick)
	}
}

type boardInputProcessor struct {
	board     *model.Board
	score     int
	messageCh chan *message.Message
}

func (bip *boardInputProcessor) processBoardInput(x int, y int) model.BlitzGameState {
	bip.board.RotatePipeClockwise(x, y)
	boardReports := bip.board.UpdateBoardPipeConnections()

	bip.score += calculateScoreFromBoardReports(boardReports)

	gameState := model.BlitzGameState{
		BoardReports: boardReports,
		Score:        bip.score,
	}
	send.SendMessageToAll(&gameState, bip.messageCh)

	return gameState
}

func (bip *boardInputProcessor) processGameBegin() {
	send.SendMessageToAll(&model.BlitzGameState{
		Board: bip.board,
		Score: bip.score,
	}, bip.messageCh)
}

func (bip *boardInputProcessor) processGameOver() {
	gameState := model.BlitzGameState{
		Score:  bip.score,
		IsOver: true,
	}
	send.SendMessageToAll(&gameState, bip.messageCh)
}

type SinglePlayerBlitzGame struct {
	*timer
	*boardInputProcessor
	playerInputChannel  chan *message.BoardInput
	playerOutputChannel chan *message.Message
}

func NewSinglePlayerBlitzGame(playerOutputChannel chan *message.Message, timeLimit time.Duration) *SinglePlayerBlitzGame {

	board := model.NewBoard(7, 7)

	return &SinglePlayerBlitzGame{
		playerInputChannel:  make(chan *message.BoardInput),
		playerOutputChannel: playerOutputChannel,
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

func (game *SinglePlayerBlitzGame) SendBoardInput(input *message.BoardInput) {
	game.playerInputChannel <- input
}

func (g *SinglePlayerBlitzGame) Run() {

	g.board.UpdateBoardPipeConnections()

	go func() {
		g.boardInputProcessor.processGameBegin()
		g.timer.countdown()
	}()

OuterLoop:
	for {
		select {
		case isOver := <-g.finishedCh:
			if isOver {
				g.boardInputProcessor.processGameOver()
				break OuterLoop
			}
		case boardInput := <-g.playerInputChannel:
			g.boardInputProcessor.processBoardInput(boardInput.X, boardInput.Y)
		}
	}

}
