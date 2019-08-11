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

	board := NewBoard(7, 7)

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
	versusLobby           *VersusLobby
	playerGameInformation map[*Player](*VersusPlayerBlitzGamePlayerInformation)
	timeLimit             time.Duration
	isOver                bool

	playerInputChannel   chan *PlayerBoardInput
	gameOverInputChannel chan bool
}

type VersusPlayerBlitzGamePlayerInformation struct {
	ID    int
	Board *Board
	Score int
}

type VersusPlayerBlitzGamePlayerInformationSentToPlayers struct {
	PlayerID         int
	EnemyInformation *VersusPlayerBlitzGameState
	Time             time.Duration
}

type VersusPlayerBlitzGameState struct {
	ID int

	Board        *Board
	BoardReports []BoardReport
	Score        int
	IsOver       bool
}

func NewVersusPlayerBlitzGame(vl *VersusLobby, timeLimit time.Duration) *VersusPlayerBlitzGame {

	playerGameInformation := make(map[*Player](*VersusPlayerBlitzGamePlayerInformation))

	i := 0
	for player := range vl.players {
		newBoard := NewBoard(7, 7)
		newBoard.UpdateBoardPipeConnections() //Note: Need to add a way to generate a board where there are no connections straight away.
		playerGameInformation[player] = &VersusPlayerBlitzGamePlayerInformation{
			i,
			&newBoard,
			0,
		}
		i++
	}

	return &VersusPlayerBlitzGame{
		versusLobby:           vl,
		playerGameInformation: playerGameInformation,
		timeLimit:             timeLimit,
		playerInputChannel:    make(chan *PlayerBoardInput),
		gameOverInputChannel:  make(chan bool),
	}
}

func (vpbg *VersusPlayerBlitzGame) Run() {

	go func() {

		for player, info := range vpbg.playerGameInformation {
			sendMessageToPlayer(&SinglePlayerBlitzGameState{
				Board: info.Board,
				Score: info.Score,
			}, player, vpbg.versusLobby.messagesToPlayersChannel)

			opponent := vpbg.getOpponent(player)

			sendMessageToPlayer(&VersusPlayerBlitzGamePlayerInformationSentToPlayers{
				EnemyInformation: &VersusPlayerBlitzGameState{
					Board: info.Board,
					Score: info.Score,
				},
			}, opponent, vpbg.versusLobby.messagesToPlayersChannel)

		}

		for !vpbg.isOver {
			vpbg.timeLimit = vpbg.timeLimit - serverTick

			//playerInformationArray := make([]*VersusPlayerBlitzGamePlayerInformation, 0)

			/*for _, info := range vpbg.playerGameInformation {
				playerInformationArray = append(playerInformationArray, info)
			}*/

			for player, info := range vpbg.playerGameInformation {
				//enemyInformation := vpbg.playerGameInformation[vpbg.getOpponent(player)]
				go sendMessageToPlayer(&VersusPlayerBlitzGamePlayerInformationSentToPlayers{
					PlayerID: info.ID,
					//EnemyInformation: enemyInformation,
					Time: vpbg.timeLimit,
				}, player, vpbg.versusLobby.messagesToPlayersChannel)
			}

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

				for player, info := range vpbg.playerGameInformation {
					sendMessageToPlayer(&SinglePlayerBlitzGameState{
						Board:  info.Board,
						IsOver: vpbg.isOver,
						Score:  info.Score,
					}, player, vpbg.versusLobby.messagesToPlayersChannel)
				}
				break OuterLoop
			}
		case playerBoardInput := <-vpbg.playerInputChannel:

			player := playerBoardInput.player
			info := vpbg.playerGameInformation[player]
			info.Board.RotatePipeClockwise(playerBoardInput.X, playerBoardInput.Y)

			boardReports := info.Board.UpdateBoardPipeConnections()

			info.Score += calculateScoreFromBoardReports(boardReports)

			gameState := SinglePlayerBlitzGameState{
				BoardReports: boardReports,
				Score:        info.Score,
				IsOver:       vpbg.isOver,
			}

			sendMessageToPlayer(gameState, player, vpbg.versusLobby.messagesToPlayersChannel)

			opponent := vpbg.getOpponent(player)
			gameStateSentToOpponent := &VersusPlayerBlitzGamePlayerInformationSentToPlayers{
				EnemyInformation: &VersusPlayerBlitzGameState{
					BoardReports: boardReports,
					Score:        info.Score,
				},
			}

			sendMessageToPlayer(gameStateSentToOpponent, opponent, vpbg.versusLobby.messagesToPlayersChannel)
		}
	}

}

func (vpbg *VersusPlayerBlitzGame) getOpponent(p *Player) *Player {
	for opponent := range vpbg.playerGameInformation {
		if p != opponent {
			return opponent
		}
	}
	return nil
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
