package game_test

import (
	"testing"

	"bryjamin.com/multiplayer/game"
	"bryjamin.com/pkg"
)

func TestBoardSolve(t *testing.T) {

	testBoard := pkg.CreateTestBoard(1, 2,
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
		{name: "Let's see how to goes", args: args{b: &testBoard}, want: []*game.Point{&game.Point{0, 1}, &game.Point{0, 1}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := game.BoardSolve(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoardSolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for index, point := range tt.want {
				if *point != *got[index] {
					t.Errorf("BoardSolve() = %v, want %v", *got[index], *point)
				}
			}
		})
	}
}
