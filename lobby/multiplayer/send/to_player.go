package send

import (
	"encoding/json"
	"log"

	"github.com/bryjammin/pipedream/lobby/multiplayer/message"
	"github.com/bryjammin/pipedream/lobby/multiplayer/player"
	"github.com/gorilla/websocket"
)

func SendMessageToPlayer(v interface{}, p *player.Player, messageToPlayerChannel chan *player.PlayerMessage) {
	messageBytes, err := json.Marshal(v)

	if err != nil {
		log.Println(err)
	} else {
		messageToPlayerChannel <- &player.PlayerMessage{Player: p, MessageType: websocket.TextMessage, Message: messageBytes}
	}
}

func SendMessageToAll(v interface{}, messageToAll chan *message.Message) {
	messageBytes, err := json.Marshal(v)

	if err != nil {
		log.Println(err)
	} else {
		messageToAll <- &message.Message{MessageType: websocket.TextMessage, Message: messageBytes}
	}
}
