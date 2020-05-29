package lobby

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bryjammin/pipedream/lobby/multiplayer/game"
	"github.com/bryjammin/pipedream/lobby/multiplayer/message"
	"github.com/bryjammin/pipedream/lobby/multiplayer/player"
)

type TutorialLobby struct {
	player *player.Player

	game *game.TutorialGame

	register   chan *player.Player
	unregister chan *player.Player

	inboundPlayerMessageCh chan *player.PlayerMessage

	outboundMessageCh chan *message.Message
}

func NewTutorialLobby() *TutorialLobby {

	return &TutorialLobby{
		player:                 nil,
		register:               make(chan *player.Player),
		unregister:             make(chan *player.Player),
		inboundPlayerMessageCh: make(chan *player.PlayerMessage),
		outboundMessageCh:      make(chan *message.Message),
	}

}

func (lobby *TutorialLobby) AddPlayer(p *player.Player) bool {

	if lobby.player == nil {
		lobby.player = p
		p.PlayerRegister = lobby
		p.PlayerMessageReceiver = lobby
		go p.Run()
		return true
	}
	return false
}

func (lobby *TutorialLobby) UnregisterPlayer(player *player.Player) {
	lobby.unregister <- player
}

func (lobby *TutorialLobby) SendMessage(message *player.PlayerMessage) {
	select {
	case lobby.inboundPlayerMessageCh <- message:
	}
}

func (l *TutorialLobby) Run() {

	go func() {
		l.game = game.NewTutorialGame(l.outboundMessageCh)
		go l.game.Run()
	}()

OuterLoop:
	for {
		select {

		case player := <-l.unregister:
			if player == l.player {
				break OuterLoop
			}
		case inboundPlayerMessage := <-l.inboundPlayerMessageCh:
			var input message.BoardInput
			err := json.Unmarshal(inboundPlayerMessage.Message, &input)
			if err == nil {
				l.game.SendBoardInput(&input)
			} else {
				log.Printf("%v", err)
			}
		case message := <-l.outboundMessageCh:
			if err := l.player.WriteMessage(message.MessageType, message.Message); err != nil {
				log.Println(err)
				return
			}
		}
	}
	fmt.Println("Lobby Closed")
}
