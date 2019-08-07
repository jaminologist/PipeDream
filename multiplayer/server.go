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
	twoPlayerRegister    chan *Player

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
		twoPlayerRegister:      make(chan *Player),
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

//FindTwoPlayerSession Creates a new WebSocket Connection with the Multiplayer server and Registers the player for finding a two player mutiplayer session
func (s *Server) FindTwoPlayerSession(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	newPlayer := newPlayer(conn)
	s.singlePlayerRegister <- newPlayer
	fmt.Println("Created Single player session")
}

func (s *Server) handleNewVersusConnection(player *Player) {
	if !s.findVersusLobby(player) {
		s.openVersusLobby(player)
	}
}

func (s *Server) findVersusLobby(player *Player) bool {
	return true
}

func (s *Server) openVersusLobby(player *Player) {

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
		case newVersusPlayer := <-s.twoPlayerRegister:
			s.handleNewVersusConnection(newVersusPlayer)
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

type PlayerBoardInput struct {
	player *Player
	BoardInput
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

type VersusLobbyManager struct {
	openVersusLobbies   []*VersusLobby
	closedVersusLobbies map[*VersusLobby](bool)

	playerHandler chan *Player

	registerLobby   chan *VersusLobby
	unregisterLobby chan *VersusLobby
}

func NewVersusLobbyManager() VersusLobbyManager {
	return VersusLobbyManager{
		openVersusLobbies:   make([]*VersusLobby, 0),
		closedVersusLobbies: make(map[*VersusLobby](bool)),

		registerLobby:   make(chan *VersusLobby),
		unregisterLobby: make(chan *VersusLobby),
	}
}

func (vlm *VersusLobbyManager) Unregister(vl *VersusLobby) {
	vlm.unregisterLobby <- vl
}

func (vlm *VersusLobbyManager) Register(vl *VersusLobby) {
	vlm.registerLobby <- vl
}

func (vlm *VersusLobbyManager) RegisterPlayer(p *Player) {
	vlm.playerHandler <- p
}

func (vlm *VersusLobbyManager) Run() {

	for {
		select {
		case newPlayer := <-vlm.playerHandler:
			vlm.handleNewPlayer(newPlayer)
		case registeringLobby := <-vlm.registerLobby:
			_ = registeringLobby
		case unregisteringLobby := <-vlm.unregisterLobby:

			_, ok := vlm.closedVersusLobbies[unregisteringLobby]
			if ok {
				delete(vlm.closedVersusLobbies, unregisteringLobby)
				log.Print("Removed Closed Lobby, Address: ", unregisteringLobby)
			}

			for i, openLobby := range vlm.openVersusLobbies {
				if openLobby == unregisteringLobby { //Delete From
					vlm.openVersusLobbies = append(vlm.openVersusLobbies[:i], vlm.openVersusLobbies[i+1:]...)
					log.Print("Removed Open Lobby, Address: ", unregisteringLobby)
				}
			}
		}
	}

}

func (vlm *VersusLobbyManager) handleNewPlayer(p *Player) {

	if len(vlm.openVersusLobbies) > 0 {
		openLobby := vlm.openVersusLobbies[0]
		openLobby.AddPlayer(p)
		if openLobby.isFull { //
			vlm.closedVersusLobbies[openLobby] = true
			vlm.openVersusLobbies = append(vlm.openVersusLobbies[:0], vlm.openVersusLobbies[0+1:]...) //Delete from open Lobbies
		}
	} else {
		newOpenLobby := NewVersusLobby(vlm)
		vlm.openVersusLobbies = append(vlm.openVersusLobbies, newOpenLobby)
		newOpenLobby.AddPlayer(p)
	}

}
