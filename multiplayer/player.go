package multiplayer

import (
	"log"
)

type Player struct {
	conn Conn

	playerRegister
	PlayerMessageReceiver
}

//AIPlayer Used to mock a player and fill spaces for waiting players
type AIPlayer struct {
}

type Conn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type playerRegister interface {
	unregisterPlayer(player *Player)
}

type PlayerMessageReceiver interface {
	SendMessage(message *PlayerMessage)
}

func newPlayer(conn Conn) *Player {
	return &Player{
		conn: conn,
	}
}

func (p *Player) run() {

	for {
		messageType, message, err := p.conn.ReadMessage()
		if err != nil {
			log.Println("Error Reading Message From Player, Unregistering Player")
			p.unregisterPlayer(p)
			return
		}

		if p.PlayerMessageReceiver != nil {
			p.SendMessage(&PlayerMessage{
				messageType: messageType,
				message:     message,
				player:      p,
			})
		}
	}

}
