package multiplayer

import (
	"fmt"
)

type Player struct {
	score int
	conn  Conn

	PlayerRegister
	PlayerMessageReceiver
}

type Conn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type PlayerRegister interface {
	Unregister(player *Player)
}

type PlayerMessageReceiver interface {
	SendMessage(message *MessageFromPlayer)
}

func newPlayer(conn Conn) *Player {
	return &Player{
		score: 0,
		conn:  conn,
	}
}

func (p *Player) run() {

	for {
		messageType, message, err := p.conn.ReadMessage()
		if err != nil {
			fmt.Println("err?")

			p.Unregister(p)
			return
		}

		p.SendMessage(&MessageFromPlayer{
			messageType: messageType,
			message:     message,
			player:      p,
		})

		/*select {
		case p.lobby.boardcast <- &MessageFromPlayer{
			messageType: messageType,
			message:     message,
			player:      p,
		}:
		}*/
	}

}
