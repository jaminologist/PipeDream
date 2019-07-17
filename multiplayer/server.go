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
	fmt.Println("yo")
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
			singlePlayerLobby.Run()
		}

		if len(s.playersLookingForLobby) > 0 {

			fmt.Println("playerlooking for lobby: ", len(s.playersLookingForLobby))

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

//Lobby manages the number of players and the interacts between players in a game
type Lobby struct {
	players map[*Player]bool

	// Register
	register chan *Player

	// Unregister
	unregister chan *Player

	boardcast chan *MessageFromPlayer

	boardcastAll chan *Message

	game *Game
}

type SinglePlayerLobby struct {
	player *Player

	game *SinglePlayerBlitzGame

	register   chan *Player
	unregister chan *Player

	playerMessageChannel chan *MessageFromPlayer

	boardcastAll chan *Message
}

func NewSinglePlayerLobby() *SinglePlayerLobby {

	return &SinglePlayerLobby{
		player:               nil,
		register:             make(chan *Player),
		unregister:           make(chan *Player),
		playerMessageChannel: make(chan *MessageFromPlayer),
		boardcastAll:         make(chan *Message),
	}

}

func (lobby *SinglePlayerLobby) AddPlayer(p *Player) bool {

	if lobby.player == nil {
		lobby.player = p
		return true
	}
	return false
}

func (lobby *SinglePlayerLobby) Unregister(player *Player) {
	lobby.unregister <- player
}

func (lobby *SinglePlayerLobby) SendMessage(message *MessageFromPlayer) {
	select {
	case lobby.playerMessageChannel <- message:
	}
}

func (l *SinglePlayerLobby) Run() {

	go func() {
		l.boardcastAll <- &Message{
			messageType: websocket.TextMessage,
			message:     []byte("found_lobby"),
		}

		l.game = NewSinglePlayerBlitzGame(l, 90*time.Second)
		go l.game.Run()
	}()

OuterLoop:
	for {
		select {

		case player := <-l.unregister:
			if player == l.player {
				break OuterLoop
			}

		case messageFromPlayer := <-l.playerMessageChannel:
			//TODO
			_ = messageFromPlayer
		case message := <-l.boardcastAll:
			if err := l.player.conn.WriteMessage(message.messageType, message.message); err != nil {
				log.Println(err)
				return
			}
		}
	}

	fmt.Println("Lobby Closed")

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

type SinglePlayerBlitzGame struct {
	singlePlayerLobby *SinglePlayerLobby
	board             *Board
	timeLimit         time.Duration
}

func NewSinglePlayerBlitzGame(spl *SinglePlayerLobby, timeLimit time.Duration) *SinglePlayerBlitzGame {

	board := NewBoard(7, 8)

	return &SinglePlayerBlitzGame{
		singlePlayerLobby: spl,
		timeLimit:         timeLimit,
		board:             &board,
	}

}

type SinglePlayerBlitzGameState struct {
	Board  *Board
	Score int
	Time   time.Duration
	IsOver bool
}

func (g *SinglePlayerBlitzGame) Run() {

	for {
		g.timeLimit = g.timeLimit - serverTick
		isOver := g.timeLimit <= 0
		messageBytes, err := json.Marshal(&SinglePlayerBlitzGameState{
			Board: g.board,
			Score: 100,
			Time: g.timeLimit, 
			IsOver: isOver,
		})

		if err != nil {
			log.Println(err)
		} else {
			g.singlePlayerLobby.boardcastAll <- &Message{messageType: websocket.TextMessage, message: messageBytes}
		}

		if isOver {
			break
		}

		time.Sleep(serverTick)
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

func NewLobby() *Lobby {

	return &Lobby{
		players:      make(map[*Player]bool),
		register:     make(chan *Player),
		unregister:   make(chan *Player),
		boardcast:    make(chan *MessageFromPlayer),
		boardcastAll: make(chan *Message),
	}

}

func (lobby *Lobby) AddPlayer(p *Player) bool {

	fmt.Println("Addplayer len is:", len(lobby.players))

	if len(lobby.players) < 2 {
		lobby.players[p] = true
		p.PlayerRegister = lobby
		p.PlayerMessageReceiver = lobby
		go p.run()
		return true
	}

	return false

}

func (lobby *Lobby) Unregister(player *Player) {
	lobby.unregister <- player
}

func (lobby *Lobby) SendMessage(message *MessageFromPlayer) {
	select {
	case lobby.boardcast <- message:
	}
}

func (l *Lobby) Run() {

	go func() {
		l.boardcastAll <- &Message{
			messageType: websocket.TextMessage,
			message:     []byte("found_lobby"),
		}

		l.game = NewGame(l, 90*time.Second)
		go l.game.Run()
	}()

OuterLoop:
	for {
		select {
		//case player := <-l.register:
		/*if len(l.players) < 2 {
			fmt.Println("New player has joined the server ")
			l.players[player] = true
			player.lobby = l
			go player.run()

			if len(l.players) == 2 {
				go func() {
					l.boardcastAll <- &Message{
						messageType: websocket.TextMessage,
						message:     []byte("found_lobby"),
					}
					l.game = NewGame(l, 60*time.Second)
					go l.game.Run()
				}()
			}

		}*/
		case player := <-l.unregister:
			if _, ok := l.players[player]; ok {
				delete(l.players, player)
			}

			if len(l.players) == 0 {
				break OuterLoop
			}

		case messageFromPlayer := <-l.boardcast:
			for player := range l.players {
				if messageFromPlayer.player != player {
					fmt.Println("Message??")

					if err := player.conn.WriteMessage(messageFromPlayer.messageType, messageFromPlayer.message); err != nil {
						log.Println(err)
						return
					}
				}
			}

		case message := <-l.boardcastAll:
			fmt.Println("Broadcast All")
			for player := range l.players {
				if err := player.conn.WriteMessage(message.messageType, message.message); err != nil {
					log.Println(err)
					return
				}
			}
		}

	}

	fmt.Println("Lobby Closed")

}

func (l *Lobby) addNewPlayer() bool {

	return true
}

func (l *Lobby) isFull() bool {
	return true
}
