package lobby

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"bryjamin.com/multiplayer/game"
	"bryjamin.com/multiplayer/message"
	"bryjamin.com/multiplayer/player"
)

const SINGLEPLAYERBLITZGAMETIMELIMIT = 60

type SinglePlayerLobby struct {
	player *player.Player

	game *game.SinglePlayerBlitzGame

	register   chan *player.Player
	unregister chan *player.Player

	playerMessageChannel chan *player.PlayerMessage

	boardcastAll chan *message.Message
}

func NewSinglePlayerLobby() *SinglePlayerLobby {

	return &SinglePlayerLobby{
		player:               nil,
		register:             make(chan *player.Player),
		unregister:           make(chan *player.Player),
		playerMessageChannel: make(chan *player.PlayerMessage),
		boardcastAll:         make(chan *message.Message),
	}

}

func (lobby *SinglePlayerLobby) AddPlayer(p *player.Player) bool {

	if lobby.player == nil {
		lobby.player = p
		p.PlayerRegister = lobby
		p.PlayerMessageReceiver = lobby
		go p.Run()
		return true
	}
	return false
}

func (lobby *SinglePlayerLobby) UnregisterPlayer(player *player.Player) {
	lobby.unregister <- player
}

func (lobby *SinglePlayerLobby) SendMessage(message *player.PlayerMessage) {
	select {
	case lobby.playerMessageChannel <- message:
	}
}

func (l *SinglePlayerLobby) Run() {

	go func() {
		l.game = game.NewSinglePlayerBlitzGame(l.boardcastAll, SINGLEPLAYERBLITZGAMETIMELIMIT*time.Second)
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
			var input message.BoardInput
			err := json.Unmarshal(messageFromPlayer.Message, &input)
			if err == nil {
				l.game.SendBoardInput(&input)
			}

		case message := <-l.boardcastAll:
			if err := l.player.WriteMessage(message.MessageType, message.Message); err != nil {
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
