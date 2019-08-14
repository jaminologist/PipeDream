package multiplayer

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const SINGLEPLAYERBLITZGAMETIMELIMIT = 90

type SinglePlayerLobby struct {
	player *Player

	game *SinglePlayerBlitzGame

	register   chan *Player
	unregister chan *Player

	playerMessageChannel chan *PlayerMessage

	boardcastAll chan *Message
}

type BoardInput struct {
	X int
	Y int
}

type PlayerBoardInput struct {
	player *Player
	BoardInput
}

func NewSinglePlayerLobby() *SinglePlayerLobby {

	return &SinglePlayerLobby{
		player:               nil,
		register:             make(chan *Player),
		unregister:           make(chan *Player),
		playerMessageChannel: make(chan *PlayerMessage),
		boardcastAll:         make(chan *Message),
	}

}

func (lobby *SinglePlayerLobby) AddPlayer(p *Player) bool {

	if lobby.player == nil {
		lobby.player = p
		p.playerRegister = lobby
		p.PlayerMessageReceiver = lobby
		go p.run()
		return true
	}
	return false
}

func (lobby *SinglePlayerLobby) unregisterPlayer(player *Player) {
	lobby.unregister <- player
}

func (lobby *SinglePlayerLobby) SendMessage(message *PlayerMessage) {
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
	game    *VersusPlayerBlitzGame

	register   chan *Player
	unregister chan *Player

	messagesToPlayersChannel chan *PlayerMessage
	playerMessageChannel     chan *PlayerMessage

	boardcastAll chan *Message

	manager *VersusLobbyManager

	isFull bool
}

type LobbyBegin struct {
	IsStarted bool
}

func NewVersusLobby(vlm *VersusLobbyManager) *VersusLobby {

	return &VersusLobby{
		players:                  make(map[*Player](bool)),
		register:                 make(chan *Player),
		unregister:               make(chan *Player),
		messagesToPlayersChannel: make(chan *PlayerMessage),
		playerMessageChannel:     make(chan *PlayerMessage),
		boardcastAll:             make(chan *Message),
		manager:                  vlm,
	}

}

func (lobby *VersusLobby) AddPlayer(p *Player) bool {

	if len(lobby.players) < 2 {
		lobby.players[p] = true
		if len(lobby.players) >= 2 {
			lobby.isFull = true
		}
		return true
	}
	return false
}

func (lobby *VersusLobby) RemovePlayer(p *Player) bool {

	if _, ok := lobby.players[p]; ok {
		delete(lobby.players, p)
		return true
	}
	return false
}

func (lobby *VersusLobby) unregisterPlayer(player *Player) {
	lobby.unregister <- player
}

func (lobby *VersusLobby) SendMessage(message *PlayerMessage) {
	select {
	case lobby.playerMessageChannel <- message:
	}
}

func (lobby *VersusLobby) Run() {

	for player := range lobby.players {
		player.playerRegister = lobby
		player.PlayerMessageReceiver = lobby
	}

	lobby.game = NewVersusPlayerBlitzGame(lobby, SINGLEPLAYERBLITZGAMETIMELIMIT*time.Second)

	go func() {
		log.Println("Beginning Versus Game...")
		sendMessageToAll(&LobbyBegin{
			IsStarted: true,
		}, lobby.boardcastAll)
		go lobby.game.Run()
	}()

OuterLoop:
	for {
		select {

		case unRegisteringPlayer := <-lobby.unregister:

			log.Println("Unregistering Player...")

			lobby.RemovePlayer(unRegisteringPlayer)
			if len(lobby.players) <= 0 {
				break OuterLoop
			}

		case messageToPlayer := <-lobby.messagesToPlayersChannel:
			if _, ok := lobby.players[messageToPlayer.player]; ok {
				if err := messageToPlayer.player.conn.WriteMessage(messageToPlayer.messageType, messageToPlayer.message); err != nil {
					log.Println(err)
				}
			}
		case messageFromPlayer := <-lobby.playerMessageChannel:
			var input BoardInput
			err := json.Unmarshal(messageFromPlayer.message, &input)
			if err == nil {
				lobby.game.playerInputChannel <- &PlayerBoardInput{player: messageFromPlayer.player, BoardInput: input}
			}
		case messageToAll := <-lobby.boardcastAll:

			for player := range lobby.players {
				if err := player.conn.WriteMessage(messageToAll.messageType, messageToAll.message); err != nil {
					log.Println("Player Connection Error: ")
				}
			}
		}
	}

	log.Println("Versus Lobby Closed")

}
