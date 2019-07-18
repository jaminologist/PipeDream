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

func TestPipe_RotateClockWise(t *testing.T) {
	tests := []struct {
		name     string
		p        *Pipe
		expected PipeDirection
	}{
		{name: "Up", p: &Pipe{Direction: UP}, expected: RIGHT},
		{name: "Right", p: &Pipe{Direction: RIGHT}, expected: DOWN},
		{name: "Down", p: &Pipe{Direction: DOWN}, expected: LEFT},
		{name: "Left", p: &Pipe{Direction: LEFT}, expected: UP},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.RotateClockWise()

			if tt.p.Direction != tt.expected {
				t.Errorf("Incorrect Direction after rotation. expected:%v, got %v", tt.expected, tt.p.Direction)
			}
		})
	}
}
