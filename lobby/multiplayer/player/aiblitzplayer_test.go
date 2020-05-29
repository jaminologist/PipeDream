package player

import (
	"encoding/json"
	"testing"

	"github.com/bryjammin/pipedream/lobby/multiplayer/game/model"
)

func TestAIBlitzPlayer_WriteMessage(t *testing.T) {
	type args struct {
		messageType int
		data        []byte
	}

	a, _ := json.Marshal(model.BlitzGameState{
		Score: 1000,
	})
	b, _ := json.Marshal(model.Pipe{
		X: 5,
		Y: 9,
	})

	tests := []struct {
		name    string
		p       *AIBlitzPlayer
		args    args
		wantErr bool
	}{
		{"Correct Gamestate input", NewAIBlitzPlayer(), args{100, a}, false},
		{"Incorrect Gamestate input", NewAIBlitzPlayer(), args{100, b}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.WriteMessage(tt.args.messageType, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("AIBlitzPlayer.WriteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
