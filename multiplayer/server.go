package multiplayer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//Server manages the number of active Lobbies
type Server struct {
	lobbyMap     map[*Lobby]bool
	emptyLobbies []*Lobby

	// Register
	register chan *Player

	playersLookingForLobby []*Player
	playersInLobby         map[*Player]bool
}

func NewServer() *Server {

	return &Server{
		lobbyMap:               make(map[*Lobby]bool),
		emptyLobbies:           make([]*Lobby, 0),
		register:               make(chan *Player),
		playersLookingForLobby: make([]*Player, 0),
		playersInLobby:         make(map[*Player]bool),
	}

}

func (s *Server) Run() {

	for {

		select {
		case newPlayer := <-s.register:
			s.playersLookingForLobby = append(s.playersLookingForLobby, newPlayer)
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

const tick time.Duration = time.Duration(600) * time.Millisecond

func (g *Game) Run() {

	for {
		g.timeLimit = g.timeLimit - tick
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

		time.Sleep(tick)
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

func (s *Server) HandleNewConnection(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	newPlayer := newPlayer(conn)
	s.register <- newPlayer
}

func (lobby *Lobby) JoinLobby(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	newPlayer := newPlayer(conn)
	lobby.register <- newPlayer
}

func (l *Lobby) AddPlayer(p *Player) bool {

	fmt.Println("Addplayer len is:", len(l.players))

	if len(l.players) < 2 {
		l.players[p] = true
		p.lobby = l
		go p.run()
		return true
	}

	return false

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
