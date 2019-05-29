package multiplayer

import (
	"testing"
	"time"
)

func TestServer_Run(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Run()
		})
	}
}

type MockConn struct{}

func (m *MockConn) ReadMessage() (messageType int, p []byte, err error) {
	return 1, make([]byte, 0), nil
}

func (m *MockConn) WriteMessage(messageType int, data []byte) error {
	return nil
}

func TestAddingMultiplePlayerToServer(t *testing.T) {

	s := NewServer()
	go func() {
		s.Run()
	}()

	s.register <- &Player{conn: &MockConn{}}
	s.register <- &Player{conn: &MockConn{}}
	s.register <- &Player{conn: &MockConn{}}
	s.register <- &Player{conn: &MockConn{}}

	time.Sleep(1 * time.Millisecond)

	if len(s.lobbyMap) != 2 {
		t.Error("You messed up: ", len(s.lobbyMap))
	}

}

func TestGame_Run(t *testing.T) {
	tests := []struct {
		name string
		g    *Game
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.Run()
		})
	}
}

func TestLobby_Run(t *testing.T) {
	tests := []struct {
		name string
		l    *Lobby
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Run()
		})
	}
}
