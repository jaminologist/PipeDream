package multiplayer

import (
	"math/rand"
)

//Board Used to describe the state of a player's pipe board
type Board struct {
	Cells [][]*Pipe

	numberOfColumns int
	numberOfRows    int
}

var allTypes = []PipeType{
	LINE,
	LPIPE,
	END,
}

var cornerTypes = []PipeType{
	LPIPE,
	END,
}

//NewBoard returns a new board with the given number of rows and columns and radomized set of pipes
func NewBoard(numberOfColumns int, numberOfRows int) Board {

	cells := make([][]*Pipe, numberOfColumns)

	for x := 0; x < numberOfColumns; x++ {
		cells[x] = make([]*Pipe, numberOfRows)

		for y := 0; y < numberOfRows; y++ {
			newPipe := newRandomizedPipe(x, y, numberOfColumns)
			cells[x][y] = &newPipe
		}
	}

	return Board{
		Cells:           cells,
		numberOfColumns: numberOfColumns,
		numberOfRows:    numberOfRows,
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
		Direction: pipeDirections[rand.Intn(len(pipeDirections))],
	}
}

//UpdateBoardPipeConnections loops through the board and checks to see which pipes are connected together
func (b *Board) UpdateBoardPipeConnections() {

	vistedPointsMap := make(map[point]bool)

	for x := 0; x < len(b.Cells); x++ {
		for y := 0; y < len(b.Cells[x]); y++ {
			vistedPointsMap[point{x, y}] = false
		}
	}

	for x := 0; x < len(b.Cells); x++ {

		for y := 0; y < len(b.Cells[x]); y++ {

			if visited, _ := vistedPointsMap[point{x, y}]; visited {
				continue
			}

			rootPipeTree := newPipeTree(b.Cells[x][y], x, y)

			isClosedTree := traversePipeTreeToCheckForClosedConnection(&rootPipeTree, vistedPointsMap, b)

			_ = isClosedTree

		}
	}

}

func (b *Board) containsPoint(p *point) bool {

	if p.x < 0 || p.x > len(b.Cells) {
		return false
	} else if p.y < 0 || p.y > len(b.Cells[p.x])-1 {
		return false
	}
	return true
}

func traversePipeTreeToCheckForClosedConnection(rootPipeTree *pipeTree, vistedPointsMap map[point]bool, board *Board) bool {

	isClosedTree := true
	pointsTo := rootPipeTree.pointsTo(rootPipeTree.x, rootPipeTree.y)

	for i := 0; i < len(pointsTo); i++ {

		point := pointsTo[i]

		if board.containsPoint(&point) {

			childTree := newPipeTree(board.Cells[point.x][point.y], point.x, point.y)
			childPointsTo := childTree.pointsTo(childTree.x, childTree.y)

			childPointsToParent := false

			for j := 0; j < len(childPointsTo); j++ {
				if childPointsTo[j].x == rootPipeTree.x && childPointsTo[j].y == rootPipeTree.y {
					childPointsToParent = true
					break
				}
			}

			if visited, _ := vistedPointsMap[point]; !visited && childPointsToParent {
				rootPipeTree.addChild(&childTree)
				isChildClosedTree := traversePipeTreeToCheckForClosedConnection(&childTree, vistedPointsMap, board)

				if isClosedTree == true {
					isClosedTree = isChildClosedTree
				}
			} else {
				isClosedTree = false
			}

		} else {
			isClosedTree = false
		}

	}

	return isClosedTree
}

type pipeTree struct {
	parent   *pipeTree
	children []*pipeTree
	*Pipe

	point
}

func newPipeTree(pipe *Pipe, x int, y int) pipeTree {
	return pipeTree{
		Pipe:  pipe,
		point: point{x, y},
	}
}

func (pt *pipeTree) addChild(childTree *pipeTree) {
	childTree.parent = pt
	pt.children = append(pt.children, childTree)
}

func (pt *pipeTree) treeSize() int {

	count := 1

	for _, child := range pt.children {
		count += child.treeSize()
	}

	return count
}

func (pt *pipeTree) rootAndChildren() []*pipeTree{

	pipeTreeSlice := make([]*pipeTree, 0, 0)

	pipeTreeSlice = append(pipeTreeSlice, pt)

	for _, child := range pt.children {
		pipeTreeSlice = append(pipeTreeSlice, child.rootAndChildren()...)
	}

	return pipeTreeSlice

}

//func newRandomizedPipe() Pipe {

//}

//func newPipe() Pipe

type Pipe struct {
	Type      PipeType
	Direction PipeDirection
	Level     PipeLevel
}

//RotateClockWise Rotates the direction of the pipe clockwise
func (p *Pipe) RotateClockWise() {
	switch p.Direction {
	case UP:
		p.Direction = RIGHT
	case RIGHT:
		p.Direction = DOWN
	case DOWN:
		p.Direction = LEFT
	case LEFT:
		p.Direction = UP
	}
}

//PointsTo Returns which x and y this pipe points to from the give x and y
func (p *Pipe) pointsTo(x int, y int) []point {

	switch p.Type {
	case END, ENDEXPLOSION2, ENDEXPLOSION3:
		return []point{pointFromDirection(point{x, y}, p.Direction)}
	case LINE:
		switch p.Direction {
		case UP, DOWN:
			return []point{pointFromDirection(point{x, y}, UP), pointFromDirection(point{x, y}, DOWN)}
		case LEFT, RIGHT:
			return []point{pointFromDirection(point{x, y}, LEFT), pointFromDirection(point{x, y}, RIGHT)}
		}
	case LPIPE:
		switch p.Direction {
		case UP:
			return []point{pointFromDirection(point{x, y}, UP), pointFromDirection(point{x, y}, RIGHT)}
		case RIGHT:
			return []point{pointFromDirection(point{x, y}, RIGHT), pointFromDirection(point{x, y}, DOWN)}
		case DOWN:
			return []point{pointFromDirection(point{x, y}, DOWN), pointFromDirection(point{x, y}, LEFT)}
		case LEFT:
			return []point{pointFromDirection(point{x, y}, LEFT), pointFromDirection(point{x, y}, UP)}

		}
	}
	return []point{}
}

func pointFromDirection(p point, d PipeDirection) point {
	switch d {
	case UP:
		return point{p.x, p.y + 1}
	case RIGHT:
		return point{p.x + 1, p.y}
	case DOWN:
		return point{p.x, p.y - 1}
	case LEFT:
		return point{p.x - 1, p.y}
	}
	return p
}

/*

#Returns which column and row this pipe points to from the give column and row
func points_to(column: int, row: int) -> Array:

    match type:
        PipeType.END, PipeType.END_EXPLOSION_2, PipeType.END_EXPLOSION_3:
            match direction:
                Direction.UP:
                    return [Vector2(column, row - 1)]
                Direction.DOWN:
                    return [Vector2(column, row + 1)]
                Direction.LEFT:
                    return [Vector2(column - 1, row)]
                Direction.RIGHT:
                    return [Vector2(column + 1, row)]
        PipeType.LINE:
            match direction:
                Direction.UP, Direction.DOWN:
                    return [Vector2(column, row + 1), Vector2(column, row - 1)]
                Direction.RIGHT, Direction.LEFT:
                    return [Vector2(column + 1, row), Vector2(column - 1, row)]
        PipeType.L_PIPE:
            match direction:
                Direction.UP:
                    return [Vector2(column + 1, row), Vector2(column, row - 1)]
                Direction.DOWN:
                    return [Vector2(column - 1, row), Vector2(column, row + 1)]
                Direction.LEFT:
                    return [Vector2(column - 1, row), Vector2(column, row - 1)]
                Direction.RIGHT:
                    return [Vector2(column + 1, row), Vector2(column, row + 1)]
    return []


*/

//PipeType the types of pipe that exist within the game.
type PipeType int

//Collection of all pipe types in the game
const (
	NONE          PipeType = -1
	LINE          PipeType = 0
	LPIPE         PipeType = 2
	END           PipeType = 4
	ENDEXPLOSION2 PipeType = 8
	ENDEXPLOSION3 PipeType = 16
)

//PipeDirection The Direction the pipe is facing
type PipeDirection int

//Collection of pipe directions set using Dir
const (
	UP    PipeDirection = 0
	RIGHT PipeDirection = 90
	DOWN  PipeDirection = 180
	LEFT  PipeDirection = 270
)

var pipeDirections = []PipeDirection{UP, RIGHT, DOWN, LEFT}

//PipeLevel Used to display the the level of the pipes connected to this pipe.
type PipeLevel int

const (
	LEVEL_0 PipeLevel = 0
	LEVEL_1 PipeLevel = 1
	LEVEL_2 PipeLevel = 2
	LEVEL_3 PipeLevel = 3
)

type point struct {
	x int
	y int
}
