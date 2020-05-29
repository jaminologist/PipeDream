package pkg

import "github.com/bryjammin/pipedream/lobby/multiplayer/game/model"

//CreateTestBoard Allows you to create a board in a human readable fashion for easier testing
func CreateTestBoard(numberOfColumns int, numberOfRows int, rowsTopToBottom ...[]*model.Pipe) model.Board {
	testBoard := model.Board{
		Cells: make([][]*model.Pipe, numberOfColumns),
	}

	for i := 0; i < len(testBoard.Cells); i++ {
		testBoard.Cells[i] = make([]*model.Pipe, numberOfRows)
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
