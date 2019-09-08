package lobby

import (
	"fmt"
	"log"
	"time"

	"bryjamin.com/multiplayer/game"
	"bryjamin.com/multiplayer/message"
	"bryjamin.com/multiplayer/player"
)

type AIBlitzLobby struct {
	player *player.Player

	game *game.AIBlitzGame

	register   chan *player.Player
	unregister chan *player.Player

	boardcastAll chan *message.Message
}

func NewAIBlitzLobby() *AIBlitzLobby {

	return &AIBlitzLobby{
		player:       nil,
		register:     make(chan *player.Player),
		unregister:   make(chan *player.Player),
		boardcastAll: make(chan *message.Message),
	}

}

func (lobby *AIBlitzLobby) AddPlayer(p *player.Player) bool {

	if lobby.player == nil {
		lobby.player = p
		p.PlayerRegister = lobby
		go p.Run()
		return true
	}
	return false
}

func (lobby *AIBlitzLobby) UnregisterPlayer(player *player.Player) {
	lobby.unregister <- player
}

func (l *AIBlitzLobby) Run() {

	go func() {
		l.game = game.NewAIBlitzGame(l.boardcastAll, SINGLEPLAYERBLITZGAMETIMELIMIT*time.Second)
		go l.game.Run()

	}()

OuterLoop:
	for {
		select {

		case player := <-l.unregister:
			if player == l.player {
				break OuterLoop
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
