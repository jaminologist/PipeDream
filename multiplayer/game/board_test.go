package game_test

import (
	"testing"
	"time"

	"bryjamin.com/multiplayer/game"
	"bryjamin.com/pkg"
)

func TestBoard_UpdateBoardPipeConnections3x3(t *testing.T) {

	testBoard := pkg.CreateTestBoard(3, 3,
		[]*game.Pipe{&game.Pipe{Type: game.END, Direction: game.DOWN}, &game.Pipe{Type: game.END, Direction: game.UP}, &game.Pipe{Type: game.END, Direction: game.UP}},
		[]*game.Pipe{&game.Pipe{Type: game.LPIPE, Direction: game.UP}, &game.Pipe{Type: game.LINE, Direction: game.LEFT}, &game.Pipe{Type: game.LPIPE, Direction: game.DOWN}},
		[]*game.Pipe{&game.Pipe{Type: game.END, Direction: game.DOWN}, &game.Pipe{Type: game.LPIPE, Direction: game.UP}, &game.Pipe{Type: game.END, Direction: game.UP}},
	)

	/*

		IMPORTANT NOTE FOR FUTURE TEST WRTING:

		AN EXPLOSIVE PIPE IS ADDED SO THERE IS ONE LESS PIPEMOVEMENTANIMATION THAN EPECTED SINCE THAT SPACE HAS BEEN FILLED

		HOWEVER, AS THAT PLACEMENT IS RANDOM DUE TO NOT CHANGING THE SEED THE TEST RUNS THE SAME, BUT ARE QUITE EASILY BREAKABLE

		THIS WILL NEED TO BE LOOKED INTO IN FUTURE AS THERE NEEDS TO BE A WAY TO TEST WITHOUT RANDOMNESS RUINING TEST DATA.
	*/

	tests := []struct {
		name     string
		b        *game.Board
		expected []game.BoardReport
	}{
		{name: "3x3 Test", b: &testBoard, expected: []game.BoardReport{
			{
				DestroyedPipes: []game.DestroyedPipe{
					game.DestroyedPipe{Type: game.END, X: 0, Y: 2},
					game.DestroyedPipe{Type: game.LPIPE, X: 0, Y: 1},
					game.DestroyedPipe{Type: game.LINE, X: 1, Y: 1},
					game.DestroyedPipe{Type: game.LPIPE, X: 2, Y: 1},
					game.DestroyedPipe{Type: game.END, X: 2, Y: 0},
				},
				PipeMovementAnimations: []game.PipeMovementAnimation{
					game.PipeMovementAnimation{X: 0, StartY: 2, EndY: 1, TravelTime: time.Millisecond * 100},
					game.PipeMovementAnimation{X: 1, StartY: 2, EndY: 1, TravelTime: time.Millisecond * 100},
					game.PipeMovementAnimation{X: 2, StartY: 2, EndY: 0, TravelTime: time.Millisecond * 200},
					game.PipeMovementAnimation{X: 0, StartY: 3, EndY: 2, TravelTime: time.Millisecond * 100},
					game.PipeMovementAnimation{X: 1, StartY: 3, EndY: 2, TravelTime: time.Millisecond * 100},
					game.PipeMovementAnimation{X: 2, StartY: 3, EndY: 1, TravelTime: time.Millisecond * 200},
					game.PipeMovementAnimation{X: 2, StartY: 4, EndY: 2, TravelTime: time.Millisecond * 200},
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

func containsDestroyedPipe(pipes []game.DestroyedPipe, pipe game.DestroyedPipe) bool {

	for i := 0; i < len(pipes); i++ {
		if pipes[i] == pipe {
			return true
		}
	}
	return false
}

func containsPipeMovementAnimation(pipes []game.PipeMovementAnimation, pipe game.PipeMovementAnimation) bool {

	for i := 0; i < len(pipes); i++ {
		if pipes[i] == pipe {
			return true
		}
	}
	return false
}
