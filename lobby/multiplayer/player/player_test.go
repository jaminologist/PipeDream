package player

import (
	"errors"
	"reflect"
	"testing"
)

type MockConn struct {
	returnError bool
}

func (m *MockConn) ReadMessage() (messageType int, p []byte, err error) {

	var mockError error

	if m.returnError {
		mockError = errors.New("Mock Error")
	}

	return 1, make([]byte, 0), mockError
}

func (m *MockConn) WriteMessage(messageType int, data []byte) error {
	return nil
}

type MockPlayerRegister struct{}

func (m *MockPlayerRegister) UnregisterPlayer(player *Player) {}

type MockPlayerMessageReceiver struct{}

func (m *MockPlayerMessageReceiver) SendMessage(message *PlayerMessage) {}

func TestNewPlayer(t *testing.T) {

	type args struct {
		conn Conn
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Conn is added to Player", args: args{conn: &MockConn{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPlayer(tt.args.conn); !reflect.DeepEqual(got.Conn, tt.args.conn) {
				t.Errorf("Connection passed to NewPlayer() not in returned Player object")
			}
		})
	}
}

func TestPlayer_run(t *testing.T) {
	tests := []struct {
		name    string
		p       *ManualPlayerRunner
		wantErr bool
	}{
		{name: "When ReadMessage errors this should return an error",
			p: &ManualPlayerRunner{Player: &Player{Conn: &MockConn{returnError: true}}}, wantErr: true},
		{name: "When ReadMessage does not error this should not return an error",
			p: &ManualPlayerRunner{&Player{Conn: &MockConn{returnError: false}}}, wantErr: false},
		{name: "When PlayerRegister is not nil and ReadMessage returns an error this should return an error",
			p: &ManualPlayerRunner{&Player{Conn: &MockConn{returnError: true}, PlayerRegister: &MockPlayerRegister{}}}, wantErr: true},
		{name: "When PlayerMessageReceiver is not nil this should not return an error",
			p: &ManualPlayerRunner{&Player{Conn: &MockConn{returnError: false}, PlayerMessageReceiver: &MockPlayerMessageReceiver{}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.run(); (err != nil) != tt.wantErr {
				t.Errorf("Player.run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
