package game

import (
	"github.com/bryjammin/pipedream/lobby/multiplayer/game/model"
	"github.com/bryjammin/pipedream/lobby/multiplayer/message"
	"github.com/bryjammin/pipedream/lobby/multiplayer/send"
)

//TutorialGame Displays a tutorial to a player
type TutorialGame struct {
	inboundBoardInputCh chan *message.BoardInput
	outboundMessageCh   chan *message.Message
	sendNewBoardCh      chan bool

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

	third := model.NewEmptyBoard(2, 5)
	third.InsertPipe(0, 0, model.LPIPE, model.DOWN)
	third.InsertPipe(0, 1, model.LINE, model.RIGHT)
	third.InsertPipe(0, 2, model.LINE, model.UP)
	third.InsertPipe(0, 3, model.END, model.UP)
	third.InsertPipe(1, 0, model.LPIPE, model.DOWN)
	third.InsertPipe(1, 1, model.LINE, model.RIGHT)
	third.InsertPipe(1, 2, model.LINE, model.UP)
	third.InsertPipe(1, 3, model.LINE, model.UP)
	third.InsertPipe(1, 4, model.END, model.UP)

	fourth := model.NewEmptyBoard(3, 3)
	fourth.InsertPipe(0, 0, model.LINE, model.DOWN)
	fourth.InsertPipe(0, 1, model.LINE, model.RIGHT)
	fourth.InsertPipe(0, 2, model.LINE, model.UP)
	fourth.InsertPipe(1, 0, model.END, model.DOWN)
	fourth.InsertPipe(1, 1, model.ENDEXPLOSION2, model.UP)
	fourth.InsertPipe(1, 2, model.LINE, model.RIGHT)
	fourth.InsertPipe(2, 0, model.LINE, model.UP)
	fourth.InsertPipe(2, 1, model.LINE, model.LEFT)
	fourth.InsertPipe(2, 2, model.LINE, model.UP)

	fifth := model.NewEmptyBoard(5, 5)
	fifth.InsertPipe(0, 0, model.LINE, model.DOWN)
	fifth.InsertPipe(0, 1, model.LINE, model.RIGHT)
	fifth.InsertPipe(0, 2, model.LPIPE, model.UP)
	fifth.InsertPipe(0, 3, model.LINE, model.UP)
	fifth.InsertPipe(0, 4, model.LPIPE, model.UP)

	fifth.InsertPipe(1, 0, model.END, model.LEFT)
	fifth.InsertPipe(1, 1, model.LINE, model.RIGHT)
	fifth.InsertPipe(1, 2, model.LINE, model.UP)
	fifth.InsertPipe(1, 3, model.LINE, model.LEFT)
	fifth.InsertPipe(1, 4, model.LINE, model.UP)

	fifth.InsertPipe(2, 0, model.LINE, model.DOWN)
	fifth.InsertPipe(2, 1, model.LINE, model.RIGHT)
	fifth.InsertPipe(2, 2, model.ENDEXPLOSION3, model.DOWN)
	fifth.InsertPipe(2, 3, model.LINE, model.UP)
	fifth.InsertPipe(2, 4, model.LINE, model.UP)

	fifth.InsertPipe(3, 0, model.LINE, model.LEFT)
	fifth.InsertPipe(3, 1, model.LINE, model.RIGHT)
	fifth.InsertPipe(3, 2, model.LINE, model.UP)
	fifth.InsertPipe(3, 3, model.LINE, model.LEFT)
	fifth.InsertPipe(3, 4, model.LINE, model.UP)

	fifth.InsertPipe(4, 0, model.LPIPE, model.LEFT)
	fifth.InsertPipe(4, 1, model.LINE, model.RIGHT)
	fifth.InsertPipe(4, 2, model.LINE, model.UP)
	fifth.InsertPipe(4, 3, model.LINE, model.UP)
	fifth.InsertPipe(4, 4, model.LPIPE, model.UP)

	return &TutorialGame{
		inboundBoardInputCh: make(chan *message.BoardInput),
		outboundMessageCh:   outboundMessageCh,
		board:               &firstBoard,
		tutorialBoards: []*model.Board{
			&firstBoard, &second, &third, &fourth, &fifth,
		},
	}
}

func (game *TutorialGame) SendBoardInput(input *message.BoardInput) {
	game.inboundBoardInputCh <- input
}

func (g *TutorialGame) Run() {

	g.board, g.tutorialBoards = g.tutorialBoards[0], g.tutorialBoards[1:]
	g.board.UpdateBoardPipeConnectionsButNoNewPipes()

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
		case boardInput := <-g.inboundBoardInputCh:

			g.board.RotatePipeClockwise(boardInput.X, boardInput.Y)
			boardReports := g.board.UpdateBoardPipeConnectionsButNoNewPipes()

			isOver := false

			if g.board.IsEmpty() {
				if len(g.tutorialBoards) == 0 {
					isOver = true
				} else {
					g.board, g.tutorialBoards = g.tutorialBoards[0], g.tutorialBoards[1:]
					newBoardReports := g.board.UpdateBoardPipeConnectionsButNoNewPipes()
					newBoardReports[0].IsNewBoard = true
					boardReports = append(boardReports, newBoardReports...)
				}
			}

			gameState := model.BlitzGameState{
				BoardReports: boardReports,
				IsOver:       isOver,
			}
			send.SendMessageToAll(&gameState, g.outboundMessageCh)

			if isOver {
				break OuterLoop
			}
		}
	}

}
