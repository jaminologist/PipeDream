package game

import (
	"bryjamin.com/multiplayer/game/model"
	"bryjamin.com/multiplayer/message"
	"bryjamin.com/multiplayer/player"
	"bryjamin.com/multiplayer/send"
	"fmt"
)

//TutorialGame Displays a tutorial to a player
type TutorialGame struct {
	inboundBoardInputCh chan *message.BoardInput
	outboundMessageCh   chan *message.Message
	sendNewBoardCh      chan bool
	finishedCh          chan bool

	board          *model.Board
	tutorialBoards []*model.Board
}

//NewTutorialGame Creates a new tutorial game
func NewTutorialGame(outboundMessageCh chan *message.Message) *TutorialGame {

	firstBoard := model.NewEmptyBoard(3, 1)
	firstBoard.InsertPipe(0, 0, model.END, model.UP)
	firstBoard.InsertPipe(1, 0, model.LINE, model.UP)
	firstBoard.InsertPipe(2, 0, model.END, model.UP)

	second := model.NewEmptyBoard(2, 2)
	second.InsertPipe(0, 0, model.LPIPE, model.DOWN)
	second.InsertPipe(0, 1, model.LPIPE, model.RIGHT)
	second.InsertPipe(1, 1, model.LPIPE, model.UP)
	second.InsertPipe(1, 0, model.LPIPE, model.LEFT)

	third := model.NewEmptyBoard(2, 2)
	third.InsertPipe(0, 0, model.LPIPE, model.DOWN)
	third.InsertPipe(0, 1, model.LPIPE, model.RIGHT)
	third.InsertPipe(1, 1, model.LPIPE, model.UP)
	third.InsertPipe(1, 0, model.LPIPE, model.LEFT)

	return &TutorialGame{
		inboundBoardInputCh: make(chan *message.BoardInput),
		outboundMessageCh:   outboundMessageCh,
		finishedCh:          make(chan bool),
		board:               &firstBoard,
		tutorialBoards: []*model.Board{
			&firstBoard, &second, &third,
		},
	}
}

func (game *TutorialGame) SendBoardInput(input *message.BoardInput) {
	game.inboundBoardInputCh <- input
}

//SendMessage receives and handles a player message
func (game *TutorialGame) SendMessage(message *player.PlayerMessage) {
	//game.inboundBoardInputCh <- input
}

func (g *TutorialGame) Run() {

	g.board, g.tutorialBoards = g.tutorialBoards[0], g.tutorialBoards[1:]

	send.SendMessageToAll(&model.BlitzGameState{
		BoardReports: []model.BoardReport{
			model.BoardReport{
				Board: g.board,
			},
		},
	}, g.outboundMessageCh)

OuterLoop:
	for {
		select {
		case isOver := <-g.finishedCh:

			fmt.Println("Up here!")

			if isOver {
				gameState := model.BlitzGameState{
					Score:  0,
					IsOver: true,
				}
				send.SendMessageToAll(&gameState, g.outboundMessageCh)
				break OuterLoop
			}
		case boardInput := <-g.inboundBoardInputCh:
			//g.boardInputProcessor.processBoardInput(boardInput.X, boardInput.Y)

			//UpdateBoardPipeConnectionsButNoNewPipes

			g.board.RotatePipeClockwise(boardInput.X, boardInput.Y)
			boardReports := g.board.UpdateBoardPipeConnectionsButNoNewPipes()

			//bip.score += calculateScoreFromBoardReports(boardReports)

			isOver := false

			if g.board.IsEmpty() {

				if len(g.tutorialBoards) == 0 {
					isOver = true
					//break
				} else {
					fmt.Printf("%+v", g.board)
					g.board, g.tutorialBoards = g.tutorialBoards[0], g.tutorialBoards[1:]
					newBoardReports := g.board.UpdateBoardPipeConnectionsButNoNewPipes()
					newBoardReports[0].IsNewBoard = true
					boardReports = append(boardReports, newBoardReports...)

					//fmt.Printf("%+v", boardReports)
					fmt.Printf("%+v", g.board)
				}
			}

			gameState := model.BlitzGameState{
				BoardReports: boardReports,
				IsOver:       isOver,
			}
			send.SendMessageToAll(&gameState, g.outboundMessageCh)
		}
	}

}
