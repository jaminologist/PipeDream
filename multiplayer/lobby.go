package multiplayer

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
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

type LobbyBegin struct {
	IsStarted bool
}
