package multiplayer

import (
	"math/rand"
)

//Board Used to describe the state of a player's pipe board
type Board struct {
	Cells [][]Pipe
}

var allTypes = []PipeType{
	LINE,
	L_PIPE,
	END,
}

var cornerTypes = []PipeType{
	L_PIPE,
	END,
}

//NewBoard returns a new board with the given number of rows and columns and radomized set of pipes
func NewBoard(numberOfColumns int, numberOfRows int) Board {

	cells := make([][]Pipe, numberOfColumns)

	for x := 0; x < numberOfColumns; x++ {
		cells[x] = make([]Pipe, numberOfRows)

		for y := 0; y < numberOfRows; y++ {
			cells[x][y] = newRandomizedPipe(x, y, numberOfColumns)
		}
	}

	/*slcB, _ := json.Marshal(Board{
		Cells: cells,
	})
	fmt.Println(string(slcB))*/

	return Board{
		Cells: cells,
	}

}

func newRandomizedPipe(x int, y int, numberOfColumns int) Pipe {
	var pipeTypesToUse []PipeType

	if x == 0 || x == numberOfColumns-1 {
		pipeTypesToUse = cornerTypes
	} else {
		pipeTypesToUse = allTypes
	}

	return Pipe{
		Type:      pipeTypesToUse[rand.Intn(len(pipeTypesToUse))],
		Direction: PipeDirections[rand.Intn(len(PipeDirections))],
	}
}

//func newRandomizedPipe() Pipe {

//}

//func newPipe() Pipe

type Pipe struct {
	Type      PipeType
	Direction PipeDirection
}

type PipeType int

const (
	NONE            PipeType = -1
	LINE            PipeType = 0
	L_PIPE          PipeType = 2
	END             PipeType = 4
	END_EXPLOSION_2 PipeType = 8
	END_EXPLOSION_3 PipeType = 16
)

type PipeDirection int

var PipeDirections = []PipeDirection{UP, RIGHT, DOWN, LEFT}

const (
	UP    PipeDirection = 0
	RIGHT PipeDirection = 90
	DOWN  PipeDirection = 180
	LEFT  PipeDirection = 270
)
