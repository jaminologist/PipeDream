package multiplayer

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	type args struct {
		numberOfColumns int
		numberOfRows    int
	}
	tests := []struct {
		name string
		args args
		//want Board
	}{
		{name: "Test New Board Contains Correct Pipe Types", args: args{numberOfColumns: 5, numberOfRows: 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			board := NewBoard(tt.args.numberOfColumns, tt.args.numberOfRows)

			for x := 0; x < len(board.Cells); x++ {
				for y := 0; y < len(board.Cells[x]); y++ {

					if x == 0 || x == len(board.Cells) {
						if board.Cells[x][y].Type == LINE {
							t.Errorf("Incorrect Pipe Type 'LINE' found in Corner part of board")
						}
					}
				}
			}
		})
	}
}
