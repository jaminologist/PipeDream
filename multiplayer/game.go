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

type TimeLimit struct {
	Time time.Duration
}

type GameOver struct {
	Time time.Duration
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
	board     *Board
	timeLimit time.Duration
	isOver    bool
	score     int

	playerInputChannel   chan *BoardInput
	playerOutputChannel  chan *Message
	gameOverInputChannel chan bool
}

type SinglePlayerBlitzGameState struct {
	Board          *Board
	BoardReports   []BoardReport
	Score          int
	IsOver         bool
	DestroyedPipes []DestroyedPipe
}

func NewSinglePlayerBlitzGame(playerOutputChannel chan *Message, timeLimit time.Duration) *SinglePlayerBlitzGame {

	board := NewBoard(7, 8)

	return &SinglePlayerBlitzGame{
		timeLimit:            timeLimit,
		board:                &board,
		playerInputChannel:   make(chan *BoardInput),
		playerOutputChannel:  playerOutputChannel,
		gameOverInputChannel: make(chan bool),
	}

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
		g.playerOutputChannel <- &Message{messageType: websocket.TextMessage, message: messageBytes}
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

type VersusPlayerBlitzGame struct {
	versusLobby *VersusLobby
	boards      map[*Player](*Board)
	scores      map[*Player](int)
	timeLimit   time.Duration
	isOver      bool

	playerInputChannel   chan *PlayerBoardInput
	gameOverInputChannel chan bool
}

func NewVersusPlayerBlitzGame(vl *VersusLobby, timeLimit time.Duration) *VersusPlayerBlitzGame {

	playerBoards := make(map[*Player](*Board))

	for player := range vl.players {
		newBoard := NewBoard(7, 8)
		newBoard.UpdateBoardPipeConnections() //Note: Need to add a way to generate a board where there are no connections straight away.
		playerBoards[player] = &newBoard
	}

	return &VersusPlayerBlitzGame{
		versusLobby:          vl,
		boards:               playerBoards,
		timeLimit:            timeLimit,
		playerInputChannel:   make(chan *PlayerBoardInput),
		gameOverInputChannel: make(chan bool),
	}
}

func (vpbg *VersusPlayerBlitzGame) Run() {

	go func() {

		for player, board := range vpbg.boards {
			sendMessageToPlayer(&SinglePlayerBlitzGameState{
				Board: board,
				Score: vpbg.scores[player],
			}, player, vpbg.versusLobby.messagesToPlayersChannel)
		}

		for !vpbg.isOver {
			vpbg.timeLimit = vpbg.timeLimit - serverTick
			sendMessageToAll(&TimeLimit{
				Time: vpbg.timeLimit,
			}, vpbg.versusLobby.boardcastAll)
			vpbg.isOver = vpbg.timeLimit <= 0
			time.Sleep(serverTick)
		}

		vpbg.gameOverInputChannel <- vpbg.isOver
	}()

OuterLoop:
	for {
		select {
		case isOver := <-vpbg.gameOverInputChannel:
			if isOver {

				for player, board := range vpbg.boards {
					sendMessageToPlayer(&SinglePlayerBlitzGameState{
						Board:  board,
						IsOver: vpbg.isOver,
						Score:  vpbg.scores[player],
					}, player, vpbg.versusLobby.messagesToPlayersChannel)
				}
				break OuterLoop
			}
		case playerBoardInput := <-vpbg.playerInputChannel:

			player := playerBoardInput.player
			board := vpbg.boards[player]
			board.RotatePipeClockwise(playerBoardInput.X, playerBoardInput.Y)

			boardReports := board.UpdateBoardPipeConnections()

			vpbg.scores[player] += calculateScoreFromBoardReports(boardReports)

			gameState := SinglePlayerBlitzGameState{
				BoardReports: boardReports,
				Score:        vpbg.scores[player],
				IsOver:       vpbg.isOver,
			}

			sendMessageToPlayer(gameState, player, vpbg.versusLobby.messagesToPlayersChannel)
		}
	}

}

func sendMessageToPlayer(v interface{}, player *Player, messageToPlayerChannel chan *MessageFromPlayer) {
	messageBytes, err := json.Marshal(v)

	if err != nil {
		log.Println(err)
	} else {
		messageToPlayerChannel <- &MessageFromPlayer{player: player, messageType: websocket.TextMessage, message: messageBytes}
	}
}

func sendMessageToAll(v interface{}, messageToAll chan *Message) {
	messageBytes, err := json.Marshal(v)

	if err != nil {
		log.Println(err)
	} else {
		messageToAll <- &Message{messageType: websocket.TextMessage, message: messageBytes}
	}
}
