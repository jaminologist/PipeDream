package player

import (
	"log"

	"bryjamin.com/multiplayer/message"
)

type Player struct {
	Conn

	PlayerRegister
	PlayerMessageReceiver
}

type PlayerMessage struct {
	MessageType int
	Message     []byte
	Player      *Player
}

type PlayerBoardInput struct {
	Player *Player
	message.BoardInput
}

//AIPlayer Used to mock a player and fill spaces for waiting players
type AIPlayer struct {
}

type Conn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type PlayerRegister interface {
	UnregisterPlayer(player *Player)
}

type PlayerMessageReceiver interface {
	SendMessage(message *PlayerMessage)
}

//NewPlayer Returns a new Player containing the given connection
func NewPlayer(conn Conn) *Player {
	return &Player{
		Conn: conn,
	}
}

func (p *Player) run() {

	for {
		messageType, message, err := p.ReadMessage()
		if err != nil {
			log.Println("Error Reading Message From Player, Unregistering Player")
			p.UnregisterPlayer(p)
			return
		}

		if p.PlayerMessageReceiver != nil {
			p.SendMessage(&PlayerMessage{
				MessageType: messageType,
				Message:     message,
				Player:      p,
			})
		}
	}

}

func (p *Player) Run() {
	p.run()
}
