package multiplayer

import (
	"encoding/json"
	"log"
	"time"
)

type VersusLobby struct {
	players map[*Player](bool)
	game    *VersusPlayerBlitzGame

	register   chan *Player
	unregister chan *Player

	lobbyToPlayerMessageCh chan *PlayerMessage
	playerToLobbyMessageCh chan *PlayerMessage

	boardcastAll chan *Message

	manager *VersusLobbyManager

	isFull bool
}

func NewVersusLobby(vlm *VersusLobbyManager) *VersusLobby {

	return &VersusLobby{
		players:                make(map[*Player](bool)),
		register:               make(chan *Player),
		unregister:             make(chan *Player),
		lobbyToPlayerMessageCh: make(chan *PlayerMessage),
		playerToLobbyMessageCh: make(chan *PlayerMessage),
		boardcastAll:           make(chan *Message),
		manager:                vlm,
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
	case lobby.playerToLobbyMessageCh <- message:
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

		case messageToPlayer := <-lobby.lobbyToPlayerMessageCh:
			if _, ok := lobby.players[messageToPlayer.player]; ok {
				if err := messageToPlayer.player.conn.WriteMessage(messageToPlayer.messageType, messageToPlayer.message); err != nil {
					log.Println(err)
				}
			}
		case messageFromPlayer := <-lobby.playerToLobbyMessageCh:
			var input BoardInput
			err := json.Unmarshal(messageFromPlayer.message, &input)
			if err == nil {
				lobby.game.playerInputChannel <- &PlayerBoardInput{player: messageFromPlayer.player, BoardInput: input}
			}
		case message := <-lobby.boardcastAll:
			for player := range lobby.players {
				if err := player.conn.WriteMessage(message.messageType, message.message); err != nil {
					log.Println("Player Connection Error: ")
				}
			}
		}
	}

	log.Println("Versus Lobby Closed")

}
