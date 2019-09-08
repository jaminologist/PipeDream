package message

type Message struct {
	MessageType int
	Message     []byte
}

type BoardInput struct {
	X int
	Y int
}
