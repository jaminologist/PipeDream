package multiplayer

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const serverTick time.Duration = time.Duration(600) * time.Millisecond

//Server manages the number of active Lobbies
type Server struct {
	versusLobbyManager *VersusLobbyManager

	singlePlayerRegister chan *Player
	twoPlayerRegister    chan *Player

	unregister chan *Player

	playersLookingForLobby []*Player
	playersInLobby         map[*Player]bool
}

func NewServer() *Server {
	versusLobbyManager := NewVersusLobbyManager()
	go versusLobbyManager.Run()

	return &Server{
		versusLobbyManager:     &versusLobbyManager,
		singlePlayerRegister:   make(chan *Player),
		twoPlayerRegister:      make(chan *Player),
		unregister:             make(chan *Player),
		playersLookingForLobby: make([]*Player, 0),
		playersInLobby:         make(map[*Player]bool),
	}

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
	log.Println("Created Single Player Session")
}

//FindTwoPlayerSession Creates a new WebSocket Connection with the Multiplayer server and Registers the player for finding a two player mutiplayer session
func (s *Server) FindTwoPlayerSession(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	newPlayer := newPlayer(conn)
	s.twoPlayerRegister <- newPlayer
	log.Println("Created Two Player Versus Session")
}

//Run starts the Server. The server handles putting players into lobbies and starting their games
func (s *Server) Run() {

	for {
		select {
		case newSinglePlayer := <-s.singlePlayerRegister:
			var singlePlayerLobby = NewSinglePlayerLobby()
			singlePlayerLobby.AddPlayer(newSinglePlayer)
			go singlePlayerLobby.Run()
		case newVersusPlayer := <-s.twoPlayerRegister:
			s.versusLobbyManager.registerPlayer(newVersusPlayer)
		}
	}
}

type Message struct {
	messageType int
	message     []byte
}

type PlayerMessage struct {
	messageType int
	message     []byte
	player      *Player
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type VersusLobbyManager struct {
	players map[*Player](*VersusLobby)

	openVersusLobbies   []*VersusLobby
	closedVersusLobbies map[*VersusLobby](bool)

	registerPlayerCh   chan *Player
	unregisterPlayerCh chan *Player

	registerLobbyCh   chan *VersusLobby
	unregisterLobbyCh chan *VersusLobby
}

func NewVersusLobbyManager() VersusLobbyManager {
	return VersusLobbyManager{
		openVersusLobbies:   make([]*VersusLobby, 0),
		closedVersusLobbies: make(map[*VersusLobby](bool)),
		registerPlayerCh:    make(chan *Player),
		unregisterPlayerCh:  make(chan *Player),
		registerLobbyCh:     make(chan *VersusLobby),
		unregisterLobbyCh:   make(chan *VersusLobby),
	}
}

func (vlm *VersusLobbyManager) unregisterLobby(vl *VersusLobby) {
	vlm.unregisterLobbyCh <- vl
}

func (vlm *VersusLobbyManager) registerLobby(vl *VersusLobby) {
	vlm.registerLobbyCh <- vl
}

func (vlm *VersusLobbyManager) unregisterPlayer(player *Player) {
	vlm.unregisterPlayerCh <- player
}

func (vlm *VersusLobbyManager) registerPlayer(p *Player) {
	log.Println("Registering New Player To VersusLobby Manager...")
	vlm.registerPlayerCh <- p
}

func (vlm *VersusLobbyManager) Run() {

	log.Println("Starting Versus Lobby Manager...")

	for {
		select {
		case newPlayer := <-vlm.registerPlayerCh:
			log.Println("Handling New Player")
			newPlayer.playerRegister = vlm
			go newPlayer.run()
			vlm.handleNewPlayer(newPlayer)
		case unregisteringPlayer := <-vlm.unregisterPlayerCh:
			log.Println("Removing Player From Open Lobby...")
			for _, lobby := range vlm.openVersusLobbies { //S
				lobby.RemovePlayer(unregisteringPlayer)
			}
		case registeringLobby := <-vlm.registerLobbyCh:
			_ = registeringLobby
		case unregisteringLobby := <-vlm.unregisterLobbyCh:
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
		default:
		}
	}

	log.Println("Stopping Versus Lobby Manager...")

}

func (vlm *VersusLobbyManager) handleNewPlayer(p *Player) {

	if len(vlm.openVersusLobbies) > 0 {
		openLobby := vlm.openVersusLobbies[0]
		openLobby.AddPlayer(p)
		if openLobby.isFull { //
			vlm.closedVersusLobbies[openLobby] = true
			vlm.openVersusLobbies = append(vlm.openVersusLobbies[:0], vlm.openVersusLobbies[0+1:]...) //Delete from open Lobbies
			go openLobby.Run()                                                                        //Begin Lobby
		}
	} else {
		newOpenLobby := NewVersusLobby(vlm)
		vlm.openVersusLobbies = append(vlm.openVersusLobbies, newOpenLobby)
		newOpenLobby.AddPlayer(p)
	}

}
