package game_test

import (
	"testing"

	"bryjamin.com/multiplayer/game"
	"bryjamin.com/pkg"
)

func TestBoardSolve(t *testing.T) {

	simpleTestBoard := pkg.CreateTestBoard(1, 2,
		[]*game.Pipe{&game.Pipe{Type: game.END, Direction: game.UP}},
		[]*game.Pipe{&game.Pipe{Type: game.END, Direction: game.DOWN}})

	infiniteLoopTestBoard := pkg.CreateTestBoard(4, 4,
		[]*game.Pipe{
			&game.Pipe{Type: game.LPIPE, Direction: game.LEFT},
			&game.Pipe{Type: game.LPIPE, Direction: game.UP},
			&game.Pipe{Type: game.LINE, Direction: game.UP},
			&game.Pipe{Type: game.LINE, Direction: game.DOWN}},
		[]*game.Pipe{
			&game.Pipe{Type: game.END, Direction: game.UP},
			&game.Pipe{Type: game.END, Direction: game.UP},
			&game.Pipe{Type: game.LPIPE, Direction: game.LEFT},
			&game.Pipe{Type: game.END, Direction: game.DOWN}},
		[]*game.Pipe{
			&game.Pipe{Type: game.LPIPE, Direction: game.RIGHT},
			&game.Pipe{Type: game.END, Direction: game.UP},
			&game.Pipe{Type: game.LPIPE, Direction: game.LEFT},
			&game.Pipe{Type: game.LPIPE, Direction: game.LEFT}},
		[]*game.Pipe{
			&game.Pipe{Type: game.LPIPE, Direction: game.UP},
			&game.Pipe{Type: game.LPIPE, Direction: game.LEFT},
			&game.Pipe{Type: game.END, Direction: game.DOWN},
			&game.Pipe{Type: game.END, Direction: game.DOWN}})

	infiniteLoopTestBoard2 := pkg.CreateTestBoard(4, 4,
		[]*game.Pipe{
			&game.Pipe{Type: game.LPIPE, Direction: game.RIGHT},
			&game.Pipe{Type: game.LPIPE, Direction: game.RIGHT},
			&game.Pipe{Type: game.END, Direction: game.LEFT},
			&game.Pipe{Type: game.LPIPE, Direction: game.RIGHT}},
		[]*game.Pipe{
			&game.Pipe{Type: game.LPIPE, Direction: game.RIGHT},
			&game.Pipe{Type: game.LINE, Direction: game.UP},
			&game.Pipe{Type: game.LPIPE, Direction: game.DOWN},
			&game.Pipe{Type: game.END, Direction: game.UP}},
		[]*game.Pipe{
			&game.Pipe{Type: game.LPIPE, Direction: game.RIGHT},
			&game.Pipe{Type: game.LPIPE, Direction: game.DOWN},
			&game.Pipe{Type: game.END, Direction: game.DOWN},
			&game.Pipe{Type: game.LPIPE, Direction: game.RIGHT}},
		[]*game.Pipe{
			&game.Pipe{Type: game.LPIPE, Direction: game.UP},
			&game.Pipe{Type: game.END, Direction: game.RIGHT},
			&game.Pipe{Type: game.END, Direction: game.DOWN},
			&game.Pipe{Type: game.LPIPE, Direction: game.DOWN}})

	type args struct {
		b *game.Board
	}
	tests := []struct {
		name    string
		args    args
		want    []*game.Point
		wantErr bool
	}{
		{
			name:    "Let's see how to goes",
			args:    args{b: &simpleTestBoard},
			want:    []*game.Point{&game.Point{0, 0}, &game.Point{0, 0}, &game.Point{0, 1}, &game.Point{0, 1}},
			wantErr: false,
		},
		{
			name: "L Pipe 3-part loop that points to end piece, but is a child pipe",
			args: args{b: &infiniteLoopTestBoard},
			want: []*game.Point{&game.Point{0, 1}, &game.Point{0, 1}, &game.Point{0, 1},
				&game.Point{1, 1}, &game.Point{1, 1}, &game.Point{1, 1}, &game.Point{0, 2}, &game.Point{0, 2}},
			wantErr: false,
		},
		{
			name: "L Pipe 3-part loop that points to end piece, but is a root pipe",
			args: args{b: &infiniteLoopTestBoard2},
			want: []*game.Point{&game.Point{0, 1}, &game.Point{0, 1}, &game.Point{0, 1},
				&game.Point{1, 0}, &game.Point{1, 0}, &game.Point{1, 0}, &game.Point{1, 2}, &game.Point{2, 1},
				&game.Point{2, 1}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := game.BoardSolve(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoardSolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for index, point := range got {
				if *point != *tt.want[index] {

					t.Logf("got:")
					for _, g := range got {
						t.Logf("%+v", g)
					}
					t.Logf("want:")
					for _, w := range tt.want {
						t.Logf("%+v", w)
					}
					t.Errorf("BoardSolve() not correct")
					break
				}
			}
		})
	}
}
