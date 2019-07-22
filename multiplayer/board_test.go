package multiplayer

import (
	"fmt"
	"reflect"
	"testing"
	"time"
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

func TestPipe_pointsTo(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		p    *Pipe
		args args
		want []point
	}{
		{name: "END", p: &Pipe{Type: END, Direction: UP}, args: args{0, 0}, want: []point{point{0, 1}}},
		{name: "END", p: &Pipe{Type: END, Direction: RIGHT}, args: args{0, 0}, want: []point{point{1, 0}}},
		{name: "END", p: &Pipe{Type: END, Direction: DOWN}, args: args{0, 0}, want: []point{point{0, -1}}},
		{name: "END", p: &Pipe{Type: END, Direction: LEFT}, args: args{0, 0}, want: []point{point{-1, 0}}},

		{name: "LINE/UP", p: &Pipe{Type: LINE, Direction: UP}, args: args{0, 0}, want: []point{point{0, 1}, point{0, -1}}},
		{name: "LINE/RIGHT", p: &Pipe{Type: LINE, Direction: RIGHT}, args: args{0, 0}, want: []point{point{-1, 0}, point{1, 0}}},
		{name: "LINE/DOWN", p: &Pipe{Type: LINE, Direction: DOWN}, args: args{0, 0}, want: []point{point{0, 1}, point{0, -1}}},
		{name: "LINE/LEFT", p: &Pipe{Type: LINE, Direction: LEFT}, args: args{0, 0}, want: []point{point{-1, 0}, point{1, 0}}},

		{name: "LPIPE/UP", p: &Pipe{Type: LPIPE, Direction: UP}, args: args{0, 0}, want: []point{point{0, 1}, point{1, 0}}},
		{name: "LPIPE/RIGHT", p: &Pipe{Type: LPIPE, Direction: RIGHT}, args: args{0, 0}, want: []point{point{1, 0}, point{0, -1}}},
		{name: "LPIPE/DOWN", p: &Pipe{Type: LPIPE, Direction: DOWN}, args: args{0, 0}, want: []point{point{0, -1}, point{-1, 0}}},
		{name: "LPIPE/LEFT", p: &Pipe{Type: LPIPE, Direction: LEFT}, args: args{0, 0}, want: []point{point{-1, 0}, point{0, 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.pointsTo(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pipe.pointsTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

//Allows you to create a board in a human readable fashion for easier testing
func createTestBoard(numberOfColumns int, numberOfRows int, rowsTopToBottom ...[]*Pipe) Board {
	testBoard := Board{
		Cells: make([][]*Pipe, numberOfColumns),
	}

	for i := 0; i < len(testBoard.Cells); i++ {
		testBoard.Cells[i] = make([]*Pipe, numberOfRows)
	}

	height := numberOfRows - 1

	for i := 0; i < len(rowsTopToBottom); i++ {
		println(i)
		for index, pipe := range rowsTopToBottom[i] {
			testBoard.Cells[index][height-i] = pipe
		}
	}

	return testBoard
}

func TestBoard_findAllClosedPipeTrees(t *testing.T) {

	//Future Note: The First array is the x-axis the inner array is the y-axis so come up with a method to
	//Better board out a board
	testBoard := createTestBoard(3, 3,
		[]*Pipe{&Pipe{END, DOWN, 0}, &Pipe{END, DOWN, 0}, &Pipe{END, DOWN, 0}},
		[]*Pipe{&Pipe{LPIPE, UP, 0}, &Pipe{LINE, LEFT, 0}, &Pipe{LPIPE, DOWN, 0}},
		[]*Pipe{&Pipe{END, DOWN, 0}, &Pipe{LPIPE, UP, 0}, &Pipe{END, UP, 0}},
	)

	/*testBoard.Cells[2] = []*Pipe{&Pipe{END, DOWN, 0}, &Pipe{END, UP, 0}, &Pipe{LPIPE, UP, 0}}
	testBoard.Cells[1] = []*Pipe{&Pipe{END, UP, 0}, &Pipe{LINE, LEFT, 0}, &Pipe{LPIPE, DOWN, 0}}
	testBoard.Cells[0] = []*Pipe{&Pipe{END, DOWN, 0}, &Pipe{LPIPE, UP, 0}, &Pipe{END, DOWN, 0}}*/

	tests := []struct {
		name string
		b    *Board
		want int
	}{
		{"3x3 board test", &testBoard, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.b.findAllClosedPipeTrees(); len(got) != tt.want {
				t.Errorf("Board.findAllClosedPipeTrees() = expected length %v, got %v", tt.want, len(got))
			}
		})
	}
}

func TestBoard_addMissingPipesToBoard(t *testing.T) {

	testBoard := createTestBoard(3, 3,
		[]*Pipe{&Pipe{END, DOWN, 0}, &Pipe{END, DOWN, 0}, &Pipe{END, DOWN, 0}},
		[]*Pipe{nil, &Pipe{LINE, LEFT, 0}, &Pipe{LPIPE, DOWN, 0}},
		[]*Pipe{&Pipe{END, DOWN, 0}, &Pipe{END, DOWN, 0}, &Pipe{END, UP, 0}},
	)

	tests := []struct {
		name                       string
		b                          *Board
		wantPipeMovementAnimations []PipeMovementAnimation
		wantMaximumTime            time.Duration
	}{
		{"Starter Grid", &testBoard, []PipeMovementAnimation{}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPipeMovementAnimations, gotMaximumTime := tt.b.addMissingPipesToBoard()
			if !reflect.DeepEqual(gotPipeMovementAnimations, tt.wantPipeMovementAnimations) {
				fmt.Println(gotPipeMovementAnimations)
				fmt.Println(gotMaximumTime)
				t.Errorf("Board.addMissingPipesToBoard() gotPipeMovementAnimations = %v, want %v", gotPipeMovementAnimations, tt.wantPipeMovementAnimations)
			}
			if !reflect.DeepEqual(gotMaximumTime, tt.wantMaximumTime) {
				t.Errorf("Board.addMissingPipesToBoard() gotMaximumTime = %v, want %v", gotMaximumTime, tt.wantMaximumTime)
			}
		})
	}
}
