package multiplayer

import (
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

type Message struct {
	messageType int
	message     []byte
}

type BoardInput struct {
	X int
	Y int
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
