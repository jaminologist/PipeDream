package game_test

import (
	"reflect"
	"testing"

	"bryjamin.com/multiplayer/game"
	"bryjamin.com/pkg"
)

func TestBoardSolve(t *testing.T) {

	_ = pkg.CreateTestBoard(1, 2,
		[]*game.Pipe{&game.Pipe{Type: game.END, Direction: game.UP}},
		[]*game.Pipe{&game.Pipe{Type: game.END, Direction: game.UP}})

	type args struct {
		b *game.Board
	}
	tests := []struct {
		name    string
		args    args
		want    []*game.Point
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := game.BoardSolve(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoardSolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoardSolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
