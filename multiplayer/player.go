package multiplayer

import (
	"fmt"
)

type Player struct {
	score int
	conn  Conn

	lobby *Lobby
}

type Conn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

func newPlayer(conn Conn) *Player {
	return &Player{
		score: 0,
		conn:  conn,
	}
}

func (p *Player) run() {

	for {
		fmt.Println("running")
		messageType, message, err := p.conn.ReadMessage()
		if err != nil {
			fmt.Println("err?")
			p.lobby.unregister <- p
			return
		}
		fmt.Println("running2")
		select {
		case p.lobby.boardcast <- &MessageFromPlayer{
			messageType: messageType,
			message:     message,
			player:      p,
		}:
		}
	}

}
