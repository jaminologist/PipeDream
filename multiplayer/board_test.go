package multiplayer_test

import (
	"testing"
	"time"

	"bryjamin.com/multiplayer"
)

//Allows you to create a board in a human readable fashion for easier testing
func createTestBoard(numberOfColumns int, numberOfRows int, rowsTopToBottom ...[]*multiplayer.Pipe) multiplayer.Board {
	testBoard := multiplayer.Board{
		Cells: make([][]*multiplayer.Pipe, numberOfColumns),
	}

	for i := 0; i < len(testBoard.Cells); i++ {
		testBoard.Cells[i] = make([]*multiplayer.Pipe, numberOfRows)
	}

	height := numberOfRows - 1

	for i := 0; i < len(rowsTopToBottom); i++ {
		for index, pipe := range rowsTopToBottom[i] {
			testBoard.Cells[index][height-i] = pipe
			pipe.X = index
			pipe.Y = height - i
		}
	}

	return testBoard
}

func TestBoard_UpdateBoardPipeConnections3x3(t *testing.T) {

	testBoard := createTestBoard(3, 3,
		[]*multiplayer.Pipe{&multiplayer.Pipe{Type: multiplayer.END, Direction: multiplayer.DOWN}, &multiplayer.Pipe{Type: multiplayer.END, Direction: multiplayer.UP}, &multiplayer.Pipe{Type: multiplayer.END, Direction: multiplayer.UP}},
		[]*multiplayer.Pipe{&multiplayer.Pipe{Type: multiplayer.LPIPE, Direction: multiplayer.UP}, &multiplayer.Pipe{Type: multiplayer.LINE, Direction: multiplayer.LEFT}, &multiplayer.Pipe{Type: multiplayer.LPIPE, Direction: multiplayer.DOWN}},
		[]*multiplayer.Pipe{&multiplayer.Pipe{Type: multiplayer.END, Direction: multiplayer.DOWN}, &multiplayer.Pipe{Type: multiplayer.LPIPE, Direction: multiplayer.UP}, &multiplayer.Pipe{Type: multiplayer.END, Direction: multiplayer.UP}},
	)

	/*

		IMPORTANT NOTE FOR FUTURE TEST WRTING:

		AN EXPLOSIVE PIPE IS ADDED SO THERE IS ONE LESS PIPEMOVEMENTANIMATION THAN EPECTED SINCE THAT SPACE HAS BEEN FILLED

		HOWEVER, AS THAT PLACEMENT IS RANDOM DUE TO NOT CHANGING THE SEED THE TEST RUNS THE SAME, BUT ARE QUITE EASILY BREAKABLE

		THIS WILL NEED TO BE LOOKED INTO IN FUTURE AS THERE NEEDS TO BE A WAY TO TEST WITHOUT RANDOMNESS RUINING TEST DATA.
	*/

	tests := []struct {
		name     string
		b        *multiplayer.Board
		expected []multiplayer.BoardReport
	}{
		{name: "3x3 Test", b: &testBoard, expected: []multiplayer.BoardReport{
			{
				DestroyedPipes: []multiplayer.DestroyedPipe{
					multiplayer.DestroyedPipe{Type: multiplayer.END, X: 0, Y: 2},
					multiplayer.DestroyedPipe{Type: multiplayer.LPIPE, X: 0, Y: 1},
					multiplayer.DestroyedPipe{Type: multiplayer.LINE, X: 1, Y: 1},
					multiplayer.DestroyedPipe{Type: multiplayer.LPIPE, X: 2, Y: 1},
					multiplayer.DestroyedPipe{Type: multiplayer.END, X: 2, Y: 0},
				},
				PipeMovementAnimations: []multiplayer.PipeMovementAnimation{
					multiplayer.PipeMovementAnimation{X: 0, StartY: 2, EndY: 1, TravelTime: time.Millisecond * 100},
					multiplayer.PipeMovementAnimation{X: 1, StartY: 2, EndY: 1, TravelTime: time.Millisecond * 100},
					multiplayer.PipeMovementAnimation{X: 2, StartY: 2, EndY: 0, TravelTime: time.Millisecond * 200},
					multiplayer.PipeMovementAnimation{X: 0, StartY: 3, EndY: 2, TravelTime: time.Millisecond * 100},
					multiplayer.PipeMovementAnimation{X: 1, StartY: 3, EndY: 2, TravelTime: time.Millisecond * 100},
					multiplayer.PipeMovementAnimation{X: 2, StartY: 3, EndY: 1, TravelTime: time.Millisecond * 200},
					multiplayer.PipeMovementAnimation{X: 2, StartY: 4, EndY: 2, TravelTime: time.Millisecond * 200},
				},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boardReports := tt.b.UpdateBoardPipeConnections()

			for index, _ := range tt.expected {
				for _, destroyedPipe := range tt.expected[index].DestroyedPipes {
					if !containsDestroyedPipe(boardReports[index].DestroyedPipes, destroyedPipe) {
						t.Errorf("Board.UpdateBoardPipeConnections() Expected Destoyed Pipe = %v, Not found inside of Destroyed Pipes = %v ", destroyedPipe, boardReports[index].DestroyedPipes)
					}
				}

				for _, pipeMovementAnimation := range tt.expected[index].PipeMovementAnimations {
					if !containsPipeMovementAnimation(boardReports[index].PipeMovementAnimations, pipeMovementAnimation) {
						t.Errorf("Board.UpdateBoardPipeConnections() Expected Destoyed Pipe = %v, Not found inside of Destroyed Pipes = %v ", pipeMovementAnimation, boardReports[index].PipeMovementAnimations)
					}
				}

				if len(boardReports[index].DestroyedPipes) != len(tt.expected[index].DestroyedPipes) {
					t.Errorf("Board.UpdateBoardPipeConnections() Incorrect length of Destroyed Pipes found. "+
						"Expected = %v, Found = %v ", tt.expected[index].DestroyedPipes, boardReports[index].DestroyedPipes)
				}

				if len(boardReports[index].PipeMovementAnimations) != len(tt.expected[index].PipeMovementAnimations) {
					t.Errorf("Board.UpdateBoardPipeConnections() Incorrect length of PipeMovementAnimations found. "+
						"Expected = %v, Found = %v ", tt.expected[index].PipeMovementAnimations, boardReports[index].PipeMovementAnimations)
				}
			}

			/*if !reflect.DeepEqual(boardReports, tt.expected) {
				t.Errorf("Board.UpdateBoardPipeConnections() = %v, want %v", boardReports, tt.expected)
			}*/
		})
	}
}

func containsDestroyedPipe(pipes []multiplayer.DestroyedPipe, pipe multiplayer.DestroyedPipe) bool {

	for i := 0; i < len(pipes); i++ {
		if pipes[i] == pipe {
			return true
		}
	}
	return false
}

func containsPipeMovementAnimation(pipes []multiplayer.PipeMovementAnimation, pipe multiplayer.PipeMovementAnimation) bool {

	for i := 0; i < len(pipes); i++ {
		if pipes[i] == pipe {
			return true
		}
	}
	return false
}
