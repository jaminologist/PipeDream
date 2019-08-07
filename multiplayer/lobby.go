package multiplayer

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

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

const SINGLEPLAYERBLITZGAMETIMELIMIT = 90

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
		p.PlayerRegister = lobby
		p.PlayerMessageReceiver = lobby
		go p.run()
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

		l.game = NewSinglePlayerBlitzGame(l.boardcastAll, SINGLEPLAYERBLITZGAMETIMELIMIT*time.Second)
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
			var input BoardInput
			err := json.Unmarshal(messageFromPlayer.message, &input)

			if err == nil {
				l.game.playerInputChannel <- &input
			}

		case message := <-l.boardcastAll:
			if err := l.player.conn.WriteMessage(message.messageType, message.message); err != nil {
				log.Println(err)
				return
			}
		}
	}

	fmt.Println("Lobby Closed")

}

type VersusLobby struct {
	players map[*Player](bool)
	games   map[*Player](*SinglePlayerBlitzGame)

	register   chan *Player
	unregister chan *Player

	playerMessageChannel chan *MessageFromPlayer

	boardcastAll chan *Message

	manager *VersusLobbyManager

	isFull bool
}

func NewVersusLobby(vlm *VersusLobbyManager) *VersusLobby {

	return &VersusLobby{
		players:              make(map[*Player](bool)),
		games:                make(map[*Player](*SinglePlayerBlitzGame)),
		register:             make(chan *Player),
		unregister:           make(chan *Player),
		playerMessageChannel: make(chan *MessageFromPlayer),
		boardcastAll:         make(chan *Message),
		manager:              vlm,
	}

}

func (lobby *VersusLobby) AddPlayer(p *Player) bool {

	if len(lobby.players) < 2 {

		p.PlayerRegister = lobby
		p.PlayerMessageReceiver = lobby

		lobby.players[p] = true
		if len(lobby.players) >= 2 {
			lobby.isFull = true
		}

		//go p.run()
		return true
	}
	return false
}

func (lobby *VersusLobby) Unregister(player *Player) {
	lobby.unregister <- player
}

func (lobby *VersusLobby) SendMessage(message *MessageFromPlayer) {
	select {
	case lobby.playerMessageChannel <- message:
	}
}

func (lobby *VersusLobby) Run() {

	for player := range lobby.players {
		go player.run()
		lobby.games[player] = NewSinglePlayerBlitzGame(lobby, SINGLEPLAYERBLITZGAMETIMELIMIT*time.Second)
	}

	go func() {
		lobby.boardcastAll <- &Message{
			messageType: websocket.TextMessage,
			message:     []byte("found_lobby"),
		}

		l.game = NewSinglePlayerBlitzGame(l, SINGLEPLAYERBLITZGAMETIMELIMIT*time.Second)
		go l.game.Run()
	}()

OuterLoop:
	for {
		select {

		case unRegisteringPlayer := <-lobby.unregister:

			if _, ok := lobby.players[unRegisteringPlayer]; ok {
				delete(lobby.players, unRegisteringPlayer)
			}

			if len(lobby.players) <= 0 {
				break OuterLoop
			}

		case messageFromPlayer := <-lobby.playerMessageChannel:
			var input BoardInput
			err := json.Unmarshal(messageFromPlayer.message, &input)

			if err == nil {
				lobby.games[messageFromPlayer.player].playerInputChannel <- &input
			}

		case message := <-lobby.boardcastAll:

			for player := range lobby.players {
				if err := player.conn.WriteMessage(message.messageType, message.message); err != nil {
					log.Println(err)
					//return
				}
			}
		}
	}

	fmt.Println("Lobby Closed")

}
