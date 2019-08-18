package multiplayer

import (
	"testing"
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

func TestVersusLobbyManager_handleNewPlayer(t *testing.T) {
	type args struct {
		//p *Player
		numberOfPlayersToAdd int
	}
	type expected struct {
		numberOfOpenLobbies   int
		numberOfClosedLobbies int
	}
	tests := []struct {
		name string
		//vlm      *VersusLobbyManager
		args     args
		expected expected
	}{
		{name: "Add 1 Player", args: args{numberOfPlayersToAdd: 1}, expected: expected{numberOfOpenLobbies: 1, numberOfClosedLobbies: 0}},
		{name: "Add 2 Players", args: args{numberOfPlayersToAdd: 2}, expected: expected{numberOfOpenLobbies: 0, numberOfClosedLobbies: 1}},
		{name: "Add 5 Players", args: args{numberOfPlayersToAdd: 5}, expected: expected{numberOfOpenLobbies: 1, numberOfClosedLobbies: 2}},
		{name: "Add 10 Players", args: args{numberOfPlayersToAdd: 10}, expected: expected{numberOfOpenLobbies: 0, numberOfClosedLobbies: 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			vlm := NewVersusLobbyManager()

			for i := 0; i < tt.args.numberOfPlayersToAdd; i++ {
				newPlayer := newPlayer(&MockConn{})
				vlm.handleNewPlayer(newPlayer)
			}

			if len(vlm.openVersusLobbies) != tt.expected.numberOfOpenLobbies || len(vlm.closedVersusLobbies) != tt.expected.numberOfClosedLobbies {
				t.Errorf("Incorrect Number Lobbies Present. Open Lobbies: %v, expected: %v. Closed Lobbies: %v, expected %v.",
					len(vlm.openVersusLobbies),
					tt.expected.numberOfOpenLobbies,
					len(vlm.closedVersusLobbies),
					tt.expected.numberOfClosedLobbies)
			}
		})
	}

}
