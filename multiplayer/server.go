package multiplayer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const serverTick time.Duration = time.Duration(600) * time.Millisecond

//Server manages the number of active Lobbies
type Server struct {
	lobbyMap     map[*Lobby]bool
	emptyLobbies []*Lobby

	// Register
	register chan *Player

	singlePlayerRegister chan *Player

	unregister chan *Player

	playersLookingForLobby []*Player
	playersInLobby         map[*Player]bool
}

func NewServer() *Server {

	return &Server{
		lobbyMap:               make(map[*Lobby]bool),
		emptyLobbies:           make([]*Lobby, 0),
		register:               make(chan *Player),
		singlePlayerRegister:   make(chan *Player),
		unregister:             make(chan *Player),
		playersLookingForLobby: make([]*Player, 0),
		playersInLobby:         make(map[*Player]bool),
	}

}

//HandleNewConnection Creates a new WebSocket Connection with the Multiplayer server and Registers the player with the server
func (s *Server) HandleNewConnection(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	newPlayer := newPlayer(conn)
	s.register <- newPlayer
}

//CreateSinglePlayerSession Creates a new WebSocket Connection with the Multiplayer server and Registers the player for a singleplayer session
func (s *Server) CreateSinglePlayerSession(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	newPlayer := newPlayer(conn)
	s.singlePlayerRegister <- newPlayer
	fmt.Println("Created Single player session")
}

//Run starts the Server. The server handles putting players into lobbies and starting their games
func (s *Server) Run() {

	for {

		select {
		case newPlayer := <-s.register:
			s.playersLookingForLobby = append(s.playersLookingForLobby, newPlayer)
		case newSinglePlayer := <-s.singlePlayerRegister:
			var singlePlayerLobby = NewSinglePlayerLobby()
			singlePlayerLobby.AddPlayer(newSinglePlayer)
			go singlePlayerLobby.Run()
		}

		if len(s.playersLookingForLobby) > 0 {

			for i := 0; i < len(s.playersLookingForLobby); i++ {
				if len(s.emptyLobbies) == 0 {
					s.emptyLobbies = append(s.emptyLobbies, NewLobby())
				}

				success := s.emptyLobbies[0].AddPlayer(s.playersLookingForLobby[i])

				if len(s.emptyLobbies[0].players) == 2 {
					fmt.Println("success is: ", success)
					var fullLobby *Lobby
					fullLobby, s.emptyLobbies = s.emptyLobbies[0], s.emptyLobbies[1:]
					s.lobbyMap[fullLobby] = true
					go fullLobby.Run()
				}
			}

			s.playersLookingForLobby = make([]*Player, 0)
		}

	}

}

type Game struct {
	lobby     *Lobby
	timeLimit time.Duration
}

type GameState struct {
	Time   time.Duration
	IsOver bool
}

type Message struct {
	messageType int
	message     []byte
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

type BoardInput struct {
	X int
	Y int
}

type SinglePlayerBlitzGame struct {
	singlePlayerLobby *SinglePlayerLobby
	board             *Board
	timeLimit         time.Duration

	score int

	playerInputChannel chan *BoardInput
}

func NewSinglePlayerBlitzGame(spl *SinglePlayerLobby, timeLimit time.Duration) *SinglePlayerBlitzGame {

	board := NewBoard(7, 8)

	return &SinglePlayerBlitzGame{
		singlePlayerLobby:  spl,
		timeLimit:          timeLimit,
		board:              &board,
		playerInputChannel: make(chan *BoardInput),
	}

}

type SinglePlayerBlitzGameState struct {
	Board          *Board
	BoardReports   []BoardReport
	Score          int
	IsOver         bool
	DestroyedPipes []DestroyedPipe
}

type TimeLimit struct {
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
			isOver := g.timeLimit <= 0
			g.send(&TimeLimit{
				Time: g.timeLimit,
			})
			if isOver {
				break
			}

			time.Sleep(serverTick)
		}
	}()

	for {
		select {
		case boardInput := <-g.playerInputChannel:
			g.board.Cells[boardInput.X][boardInput.Y].RotateClockWise()
			boardReports := g.board.UpdateBoardPipeConnections()

			pipesDestroyed := 0
			for i := 0; i < len(boardReports); i++ {
				pipesDestroyed += len(boardReports[i].DestroyedPipes)
			}

			g.score += 1250 * pipesDestroyed

			fmt.Println(g.score)
			gameState := SinglePlayerBlitzGameState{
				BoardReports: boardReports,
				Score:        g.score,
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

type MessageFromPlayer struct {
	messageType int
	message     []byte
	player      *Player
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
