package player

import (
	"log"

	"bryjamin.com/multiplayer/message"
)

type Player struct {
	Conn

	PlayerRegister
	PlayerMessageReceiver
	PlayerRunner
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

type PlayerRunner interface {
	Run()
}

//NewPlayer Returns a new Player containing the given connection
func NewPlayer(conn Conn) *Player {

	player := &Player{
		Conn: conn,
	}

	player.PlayerRunner = &ManualPlayerRunner{
		Player: player,
	}

	return player
}

type ManualPlayerRunner struct {
	*Player
}

func (p *ManualPlayerRunner) Run() {
	for {
		err := p.run()
		if err != nil {
			return
		}
	}
}

func (p *ManualPlayerRunner) run() error {
	messageType, message, err := p.ReadMessage()
	if err != nil {
		log.Println("Error Reading Message From Player, Unregistering Player")
		if p.PlayerRegister != nil {
			p.UnregisterPlayer(p.Player)
		}
		return err
	}

	if p.PlayerMessageReceiver != nil {
		p.SendMessage(&PlayerMessage{
			MessageType: messageType,
			Message:     message,
			Player:      p.Player,
		})
	}
	return nil
}
