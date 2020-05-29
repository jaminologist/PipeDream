package model_test

import (
	"testing"

	"github.com/bryjammin/pipedream/lobby/multiplayer/game/model"
	"github.com/bryjammin/pipedream/lobby/pkg"
)

func TestBoardSolve(t *testing.T) {

	simpleTestBoard := pkg.CreateTestBoard(1, 2, []*model.Pipe{&model.Pipe{Type: model.END, Direction: model.UP}},
		[]*model.Pipe{&model.Pipe{Type: model.END, Direction: model.DOWN}})

	infiniteLoopTestBoard := pkg.CreateTestBoard(4, 4,
		[]*model.Pipe{
			&model.Pipe{Type: model.LPIPE, Direction: model.LEFT},
			&model.Pipe{Type: model.LPIPE, Direction: model.UP},
			&model.Pipe{Type: model.LINE, Direction: model.UP},
			&model.Pipe{Type: model.LINE, Direction: model.DOWN}},
		[]*model.Pipe{
			&model.Pipe{Type: model.END, Direction: model.UP},
			&model.Pipe{Type: model.END, Direction: model.UP},
			&model.Pipe{Type: model.LPIPE, Direction: model.LEFT},
			&model.Pipe{Type: model.END, Direction: model.DOWN}},
		[]*model.Pipe{
			&model.Pipe{Type: model.LPIPE, Direction: model.RIGHT},
			&model.Pipe{Type: model.END, Direction: model.UP},
			&model.Pipe{Type: model.LPIPE, Direction: model.LEFT},
			&model.Pipe{Type: model.LPIPE, Direction: model.LEFT}},
		[]*model.Pipe{
			&model.Pipe{Type: model.LPIPE, Direction: model.UP},
			&model.Pipe{Type: model.LPIPE, Direction: model.LEFT},
			&model.Pipe{Type: model.END, Direction: model.DOWN},
			&model.Pipe{Type: model.END, Direction: model.DOWN}})

	infiniteLoopTestBoard2 := pkg.CreateTestBoard(4, 4,
		[]*model.Pipe{
			&model.Pipe{Type: model.LPIPE, Direction: model.RIGHT},
			&model.Pipe{Type: model.LPIPE, Direction: model.RIGHT},
			&model.Pipe{Type: model.END, Direction: model.LEFT},
			&model.Pipe{Type: model.LPIPE, Direction: model.RIGHT}},
		[]*model.Pipe{
			&model.Pipe{Type: model.LPIPE, Direction: model.RIGHT},
			&model.Pipe{Type: model.LINE, Direction: model.UP},
			&model.Pipe{Type: model.LPIPE, Direction: model.DOWN},
			&model.Pipe{Type: model.END, Direction: model.UP}},
		[]*model.Pipe{
			&model.Pipe{Type: model.LPIPE, Direction: model.RIGHT},
			&model.Pipe{Type: model.LPIPE, Direction: model.DOWN},
			&model.Pipe{Type: model.END, Direction: model.DOWN},
			&model.Pipe{Type: model.LPIPE, Direction: model.RIGHT}},
		[]*model.Pipe{
			&model.Pipe{Type: model.LPIPE, Direction: model.UP},
			&model.Pipe{Type: model.END, Direction: model.RIGHT},
			&model.Pipe{Type: model.END, Direction: model.DOWN},
			&model.Pipe{Type: model.LPIPE, Direction: model.DOWN}})

	type args struct {
		b *model.Board
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.Point
		wantErr bool
	}{
		{
			name:    "Let's see how to goes",
			args:    args{b: &simpleTestBoard},
			want:    []*model.Point{&model.Point{0, 0}, &model.Point{0, 0}, &model.Point{0, 1}, &model.Point{0, 1}},
			wantErr: false,
		},
		{
			name: "L Pipe 3-part loop that points to end piece, but is a child pipe",
			args: args{b: &infiniteLoopTestBoard},
			want: []*model.Point{&model.Point{0, 1}, &model.Point{0, 1}, &model.Point{0, 1},
				&model.Point{1, 1}, &model.Point{1, 1}, &model.Point{1, 1}, &model.Point{0, 2}, &model.Point{0, 2}},
			wantErr: false,
		},
		{
			name: "L Pipe 3-part loop that points to end piece, but is a root pipe",
			args: args{b: &infiniteLoopTestBoard2},
			want: []*model.Point{&model.Point{0, 1}, &model.Point{0, 1}, &model.Point{0, 1},
				&model.Point{1, 0}, &model.Point{1, 0}, &model.Point{1, 0}, &model.Point{1, 2}, &model.Point{2, 1},
				&model.Point{2, 1}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := model.BoardSolve(tt.args.b)
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
