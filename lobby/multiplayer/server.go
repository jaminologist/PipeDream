package multiplayer

import (
	"log"
	"net/http"

	"github.com/bryjammin/pipedream/lobby/multiplayer/lobby"
	"github.com/bryjammin/pipedream/lobby/multiplayer/player"
	"github.com/gorilla/websocket"
)

//Server manages the number of active Lobbies
type Server struct {
	versusLobbyManager *VersusLobbyManager

	singlePlayerRegister  chan *player.Player
	aiBlitzPlayerRegister chan *player.Player
	twoPlayerRegister     chan *player.Player
	versusAIRegister      chan *player.Player
	tutortialRegister     chan *player.Player

	unregister chan *player.Player

	playersLookingForLobby []*player.Player
	playersInLobby         map[*player.Player]bool
}

//NewServer Creates a new PipeDream server that handles the regsitering of all players to different game modes
func NewServer() *Server {
	versusLobbyManager := NewVersusLobbyManager()
	go versusLobbyManager.Run()

	return &Server{
		versusLobbyManager:     &versusLobbyManager,
		singlePlayerRegister:   make(chan *player.Player),
		aiBlitzPlayerRegister:  make(chan *player.Player),
		twoPlayerRegister:      make(chan *player.Player),
		versusAIRegister:       make(chan *player.Player),
		tutortialRegister:      make(chan *player.Player),
		unregister:             make(chan *player.Player),
		playersLookingForLobby: make([]*player.Player, 0),
		playersInLobby:         make(map[*player.Player]bool),
	}

}

func registerPlayer(w http.ResponseWriter, r *http.Request, registerCh chan *player.Player) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	newPlayer := player.NewPlayer(conn)
	registerCh <- newPlayer
}

//CreateSinglePlayerSession Creates a new WebSocket Connection with the Multiplayer server and Registers the player for a singleplayer session
func (s *Server) CreateSinglePlayerSession(w http.ResponseWriter, r *http.Request) {
	registerPlayer(w, r, s.singlePlayerRegister)
	log.Println("Created Single Player Session")
}

//FindTwoPlayerSession Creates a new WebSocket Connection with the Multiplayer server and Registers the player for finding a two player mutiplayer session
func (s *Server) FindTwoPlayerSession(w http.ResponseWriter, r *http.Request) {
	registerPlayer(w, r, s.twoPlayerRegister)
	log.Println("Created Two Player Versus Session")
}

//FindAISession registers player for AI Blitz Game
func (s *Server) FindAISession(w http.ResponseWriter, r *http.Request) {
	registerPlayer(w, r, s.aiBlitzPlayerRegister)
	log.Println("Created AI Player Session")
}

//FindVersusAISession registers player for Versus AI Game
func (s *Server) FindVersusAISession(w http.ResponseWriter, r *http.Request) {
	registerPlayer(w, r, s.versusAIRegister)
	log.Println("Created Versus AI Player Session")
}

//FindTutorialSession creates a tutorial session for a player
func (s *Server) FindTutorialSession(w http.ResponseWriter, r *http.Request) {
	registerPlayer(w, r, s.tutortialRegister)
	log.Println("Created Tutorial Player Session")
}

//Run starts the Server. The server handles putting players into lobbies and starting their games
func (s *Server) Run() {

	for {
		select {
		case newSinglePlayer := <-s.singlePlayerRegister:
			var singlePlayerLobby = lobby.NewSinglePlayerLobby()
			singlePlayerLobby.AddPlayer(newSinglePlayer)
			go singlePlayerLobby.Run()
		case newTutorialPlayer := <-s.tutortialRegister:
			var tutorialLobby = lobby.NewTutorialLobby()
			tutorialLobby.AddPlayer(newTutorialPlayer)
			go tutorialLobby.Run()
		case newAiBlitzPlayer := <-s.aiBlitzPlayerRegister:
			var aiPlayerLobby = lobby.NewAIBlitzLobby()
			aiPlayerLobby.AddPlayer(newAiBlitzPlayer)
			go aiPlayerLobby.Run()
		case newVersusPlayer := <-s.twoPlayerRegister:
			s.versusLobbyManager.registerPlayer(newVersusPlayer)
		case newVersusAiPlayer := <-s.versusAIRegister:
			aiVersusLobby := lobby.NewVersusLobby()
			aiVersusLobby.AddPlayer(newVersusAiPlayer)
			go newVersusAiPlayer.Run()

			aiPlayer := player.NewAIBlitzPlayer()
			go aiPlayer.Run()
			aiVersusLobby.AddPlayer(aiPlayer)

			go aiVersusLobby.Run()
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type VersusLobbyManager struct {
	players map[*player.Player](*lobby.VersusLobby)

	openVersusLobbies   []*lobby.VersusLobby
	closedVersusLobbies map[*lobby.VersusLobby](bool)

	registerPlayerCh   chan *player.Player
	unregisterPlayerCh chan *player.Player

	registerLobbyCh   chan *lobby.VersusLobby
	unregisterLobbyCh chan *lobby.VersusLobby
}

func NewVersusLobbyManager() VersusLobbyManager {
	return VersusLobbyManager{
		openVersusLobbies:   make([]*lobby.VersusLobby, 0),
		closedVersusLobbies: make(map[*lobby.VersusLobby](bool)),
		registerPlayerCh:    make(chan *player.Player),
		unregisterPlayerCh:  make(chan *player.Player),
		registerLobbyCh:     make(chan *lobby.VersusLobby),
		unregisterLobbyCh:   make(chan *lobby.VersusLobby),
	}
}

func (vlm *VersusLobbyManager) unregisterLobby(vl *lobby.VersusLobby) {
	vlm.unregisterLobbyCh <- vl
}

func (vlm *VersusLobbyManager) registerLobby(vl *lobby.VersusLobby) {
	vlm.registerLobbyCh <- vl
}

func (vlm *VersusLobbyManager) UnregisterPlayer(player *player.Player) {
	vlm.unregisterPlayerCh <- player
}

func (vlm *VersusLobbyManager) registerPlayer(p *player.Player) {
	log.Println("Registering New Player To VersusLobby Manager...")
	vlm.registerPlayerCh <- p
}

func (vlm *VersusLobbyManager) Run() {

	log.Println("Starting Versus Lobby Manager...")

	for {
		select {
		case newPlayer := <-vlm.registerPlayerCh:
			log.Println("Handling New Player")
			newPlayer.PlayerRegister = vlm
			go newPlayer.Run()
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

func (vlm *VersusLobbyManager) handleNewPlayer(p *player.Player) {

	if len(vlm.openVersusLobbies) > 0 {
		openLobby := vlm.openVersusLobbies[0]
		openLobby.AddPlayer(p)
		if openLobby.IsFull() { //
			vlm.closedVersusLobbies[openLobby] = true
			vlm.openVersusLobbies = append(vlm.openVersusLobbies[:0], vlm.openVersusLobbies[0+1:]...) //Delete from open Lobbies
			go openLobby.Run()                                                                        //Begin Lobby
		}
	} else {
		newOpenLobby := lobby.NewVersusLobby()
		vlm.openVersusLobbies = append(vlm.openVersusLobbies, newOpenLobby)
		newOpenLobby.AddPlayer(p)
	}

}
