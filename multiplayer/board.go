package multiplayer

import (
	"math/rand"
	"time"
)

//Board Used to describe the state of a player's pipe board
type Board struct {
	Cells [][]*Pipe

	NumberOfColumns int
	NumberOfRows    int
}

//Pipe Represents a pipe within the game has a Type, Direction and 'Level' and an X and Y position
type Pipe struct {
	Type      PipeType
	Direction PipeDirection
	Level     PipeLevel
	X         int
	Y         int
}

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

//BoardReport Sends back information about board updates that can be used to calculate client animations
type BoardReport struct {
	DestroyedPipes         []DestroyedPipe
	PipeMovementAnimations []PipeMovementAnimation
	MaximumAnimationTime   time.Duration
	Board                  *Board
}

//PipeMovementAnimation Used to detail information for pipe 'Falling' animation. From StartY to EndY and expected travel time
type PipeMovementAnimation struct {
	X          int
	StartY     int
	EndY       int
	TravelTime time.Duration
}

//DestroyedPipe Used to detail information for Pipe 'Destroyed' animation using X and Y position and pipe Type
type DestroyedPipe struct {
	Type PipeType

	X int
	Y int
}

type point struct {
	x int
	y int
}

var allTypes = []PipeType{
	LINE,
	LPIPE,
	END,
	//ENDEXPLOSION2,
}

var cornerTypes = []PipeType{
	LPIPE,
	END,
}

var level1Size = 2
var level2Size = 4
var level3Size = 6

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
		NumberOfColumns: numberOfColumns,
		NumberOfRows:    numberOfRows,
	}
}

func NewEmptyBoard(numberOfColumns int, numberOfRows int) Board {
	cells := make([][]*Pipe, numberOfColumns)
	for x := 0; x < numberOfColumns; x++ {
		cells[x] = make([]*Pipe, numberOfRows)
	}

	return Board{
		Cells:           cells,
		NumberOfColumns: numberOfColumns,
		NumberOfRows:    numberOfRows,
	}
}

func newRandomizedPipe(x int, y int, numberOfColumns int) Pipe {
	var pipeTypesToUse []PipeType

	if x == 0 || x == numberOfColumns-1 {
		pipeTypesToUse = cornerTypes
	} else {
		pipeTypesToUse = allTypes
	}

	return newPipe(x, y,
		pipeTypesToUse[rand.Intn(len(pipeTypesToUse))],
		getRandomPipeDirection())
}

func newPipe(x int, y int, pipeType PipeType, pipeDirection PipeDirection) Pipe {
	return Pipe{
		X:         x,
		Y:         y,
		Type:      pipeType,
		Direction: pipeDirection,
	}
}

func CopyBoard(b *Board) Board {

	newBoard := NewEmptyBoard(b.NumberOfColumns, b.NumberOfRows)

	for x := 0; x < b.NumberOfColumns; x++ {
		for y := 0; y < b.NumberOfRows; y++ {
			pipe := b.Cells[x][y]
			newBoard.Cells[x][y] = &Pipe{
				Type:      pipe.Type,
				Direction: pipe.Direction,
				Level:     pipe.Level,
				X:         pipe.X,
				Y:         pipe.Y,
			}
		}
	}

	return newBoard
}

//RotatePipeClockwise Rotates the pipe at the given x and y clockwise if the board contains the given x and y
func (b *Board) RotatePipeClockwise(x int, y int) {
	if b.containsPoint(&point{x, y}) {
		b.Cells[x][y].RotateClockWise()
	}
}

//UpdateBoardPipeConnections loops through the board and checks to see which pipes are connected together
//Returns true if a connections is found
func (b *Board) UpdateBoardPipeConnections() []BoardReport {

	connectionFound := true

	boardReports := make([]BoardReport, 0, 0)

	for connectionFound {

		boardReport := BoardReport{}

		closedTrees := b.findAllClosedPipeTrees()

		connectionFound = len(closedTrees) > 0
		//func calculatenewpositionsforexplosivespipes (as well as the exploding pipes)
		//delete pipes
		boardReport.DestroyedPipes = b.deletePipeTreesFromBoard(closedTrees)

		//add in explosive pipes
		b.addSpecialPipesToBoardUsingClosedTrees(closedTrees)

		//add in new pipes into the empty slots
		boardReport.PipeMovementAnimations, boardReport.MaximumAnimationTime = b.addMissingPipesToBoard()
		//return if the number of connects was larger than zero

		copyBoard := CopyBoard(b)
		boardReport.Board = &copyBoard

		boardReports = append(boardReports, boardReport)
	}

	return boardReports
}

func (b *Board) findAllClosedPipeTrees() []*pipeTree {

	visitedPoints := make(map[point]bool)

	closedTrees := make([]*pipeTree, 0, 0)

	for x := 0; x < len(b.Cells); x++ {

		for y := 0; y < len(b.Cells[x]); y++ {

			if _, visited := visitedPoints[point{x, y}]; visited {
				continue
			}

			visitedPoints[point{x, y}] = true
			rootPipeTree := newPipeTree(b.Cells[x][y], x, y)

			isClosedTree := traversePipeTreeToCheckForClosedConnection(&rootPipeTree, visitedPoints, b)

			if isClosedTree {
				closedTrees = append(closedTrees, &rootPipeTree)
			}

			pipeTrees := rootPipeTree.rootAndChildren()

			size := len(pipeTrees)

			for _, pipeTree := range pipeTrees {

				switch {
				case size < level1Size:
					pipeTree.Pipe.Level = LEVEL_0
				case size < level2Size:
					pipeTree.Pipe.Level = LEVEL_1
				case size < level3Size:
					pipeTree.Pipe.Level = LEVEL_2
				case size >= level3Size:
					pipeTree.Pipe.Level = LEVEL_3
				}
			}
		}
	}

	return closedTrees

}

func (b *Board) deletePipeTreesFromBoard(pipeTrees []*pipeTree) []DestroyedPipe {

	destroyedPipes := make([]DestroyedPipe, 0, 0)

	for _, rootpipeTree := range pipeTrees {
		for _, pipeTree := range rootpipeTree.rootAndChildren() {
			pipe := b.Cells[pipeTree.x][pipeTree.y]
			if pipe != nil {
				if destroyedPipe, ok := b.deleteFromBoard(pipeTree.x, pipeTree.y); ok {
					destroyedPipes = append(destroyedPipes, destroyedPipe...)
				}
			}
		}
	}

	return destroyedPipes
}

func (b *Board) deleteFromBoard(x int, y int) ([]DestroyedPipe, bool) {

	destroyedPipes := make([]DestroyedPipe, 0, 0)

	if b.containsPoint(&point{x, y}) {

		if b.Cells[x][y] != nil {
			pipe := b.Cells[x][y]

			destroyedPipes = append(destroyedPipes, DestroyedPipe{Type: pipe.Type, X: pipe.X, Y: pipe.Y})

			b.Cells[x][y] = nil

			if pipe.Type == ENDEXPLOSION2 {
				for _, point := range positionInSquareRange(pipe.X, pipe.Y, 2) {
					if destroyedPipe, ok := b.deleteFromBoard(point.x, point.y); ok {
						destroyedPipes = append(destroyedPipes, destroyedPipe...)
					}
				}
			} else if pipe.Type == ENDEXPLOSION3 {
				for _, point := range positionInSquareRange(pipe.X, pipe.Y, 3) {
					if destroyedPipe, ok := b.deleteFromBoard(point.x, point.y); ok {
						destroyedPipes = append(destroyedPipes, destroyedPipe...)
					}
				}
			}

			return destroyedPipes, true
		}
	}

	return []DestroyedPipe{}, false
}

func positionInSquareRange(startX int, startY int, size int) []point {
	positions := make([]point, 0, 0)

	if size <= 0 {
		return positions
	}

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if x == 0 && y == 0 {
				continue
			}

			positions = append(positions,
				point{x: startX + x, y: startY + y},
				point{x: startX + x, y: startY - y},
				point{x: startX - x, y: startY + y},
				point{x: startX - x, y: startY - y},
			)
		}
	}

	return positions
}

func getRandomPipeDirection() PipeDirection {
	return pipeDirections[rand.Intn(len(pipeDirections))]
}

func (b *Board) addSpecialPipesToBoardUsingClosedTrees(rootPipeTrees []*pipeTree) {

	for _, rootpipeTree := range rootPipeTrees {

		allPipes := rootpipeTree.rootAndChildren()
		pipeTree := allPipes[rand.Intn(len(allPipes))]

		switch rootpipeTree.Level {
		case LEVEL_2:
			newPipe := newPipe(pipeTree.X, pipeTree.Y, ENDEXPLOSION2, getRandomPipeDirection())
			b.Cells[pipeTree.X][pipeTree.Y] = &newPipe
		case LEVEL_3:
			newPipe := newPipe(pipeTree.X, pipeTree.Y, ENDEXPLOSION3, getRandomPipeDirection())
			b.Cells[pipeTree.X][pipeTree.Y] = &newPipe
		}
	}
}

func getTravelTime(i int) time.Duration {
	return time.Duration(i) * (time.Millisecond * 100)
}

func (b *Board) addMissingPipesToBoard() (pipeMovementAnimations []PipeMovementAnimation, maximumTime time.Duration) {

	pipeMovementAnimations = make([]PipeMovementAnimation, 0, 0)
	//Lowers Pipes above a Gap
	for x := 0; x < len(b.Cells); x++ {

		resetPosition := 0

		for y := 0; y < len(b.Cells[x]); y++ {

			if b.Cells[x][y] != nil {

				if y > resetPosition {
					pipe := b.Cells[x][y]
					b.Cells[x][y] = nil
					b.Cells[x][resetPosition] = pipe
					pipe.X = x
					pipe.Y = resetPosition

					travelTime := getTravelTime(y - resetPosition)

					if travelTime > maximumTime {
						maximumTime = travelTime
					}

					pipeMovementAnimations = append(pipeMovementAnimations, PipeMovementAnimation{
						X:          x,
						StartY:     y,
						EndY:       resetPosition,
						TravelTime: travelTime,
					})

					y = resetPosition - 1

				} else {
					resetPosition++
				}
			}
		}
	}

	//Fills all Gaps, assumes once a gap is found the rest is also empty
	for x := 0; x < len(b.Cells); x++ {
		depth := 0
		height := len(b.Cells[x])
		for y := 0; y < height; y++ {
			if b.Cells[x][y] == nil {

				newPipe := newRandomizedPipe(x, y, b.NumberOfColumns)
				b.Cells[x][y] = &newPipe

				startY := (height) + (y - depth)
				endY := y
				travelTime := getTravelTime(startY - endY)

				if travelTime > maximumTime {
					maximumTime = travelTime
				}
				pipeMovementAnimations = append(pipeMovementAnimations, PipeMovementAnimation{
					X:          x,
					StartY:     startY,
					EndY:       endY,
					TravelTime: travelTime,
				})
			} else {
				depth++
			}
		}
	}

	return
	//newPipe := newRandomizedPipe(x, resetPosition, b.NumberOfColumns)
	//b.Cells[x][resetPosition] = &newPipe
}

func (b *Board) containsPoint(p *point) bool {
	if p.x < 0 || p.x > len(b.Cells)-1 {
		return false
	} else if p.y < 0 || p.y > len(b.Cells[p.x])-1 {
		return false
	}
	return true
}

func traversePipeTreeToCheckForClosedConnection(rootPipeTree *pipeTree, visitedPoints map[point]bool, board *Board) bool {

	isClosedTree := true
	pointsTo := rootPipeTree.pointsTo(rootPipeTree.x, rootPipeTree.y)

	for i := 0; i < len(pointsTo); i++ {

		pointToPoint := pointsTo[i]

		if board.containsPoint(&pointToPoint) {

			childTree := newPipeTree(board.Cells[pointToPoint.x][pointToPoint.y], pointToPoint.x, pointToPoint.y)
			childPointsTo := childTree.pointsTo(childTree.x, childTree.y)

			childPointsToParent := false

			for j := 0; j < len(childPointsTo); j++ {

				if childPointsTo[j].x == rootPipeTree.x && childPointsTo[j].y == rootPipeTree.y {
					childPointsToParent = true
					break
				}
			}

			if childPointsToParent {

				if _, visited := visitedPoints[pointToPoint]; !visited {
					visitedPoints[pointToPoint] = true
					rootPipeTree.addChild(&childTree)
					isChildClosedTree := traversePipeTreeToCheckForClosedConnection(&childTree, visitedPoints, board)

					if isClosedTree == true {
						isClosedTree = isChildClosedTree
					}
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
	Children []*pipeTree
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
	pt.Children = append(pt.Children, childTree)
}

func (pt *pipeTree) treeSize() int {

	count := 1

	for _, child := range pt.Children {
		count += child.treeSize()
	}

	return count
}

func (pt *pipeTree) rootAndChildren() []*pipeTree {

	pipeTreeSlice := make([]*pipeTree, 0, 0)

	pipeTreeSlice = append(pipeTreeSlice, pt)

	for _, child := range pt.Children {
		pipeTreeSlice = append(pipeTreeSlice, child.rootAndChildren()...)
	}

	return pipeTreeSlice

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
