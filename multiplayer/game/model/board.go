package model

import (
	"fmt"
	"math/rand"
	"time"
)

//Board Used to describe the state of a player's pipe board
type Board struct {
	Cells [][]*Pipe

	NumberOfColumns int
	NumberOfRows    int
}

//BoardReport Sends back information about board updates that can be used to calculate client animations
type BoardReport struct {
	DestroyedPipes         []DestroyedPipe
	PipeMovementAnimations []PipeMovementAnimation
	MaximumAnimationTime   time.Duration
	Board                  *Board
	IsNewBoard             bool
}

type Point struct {
	X int
	Y int
}

func newPoint(x int, y int) *Point {
	return &Point{x, y}
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

//InsertPipe inserts a pipe of given type and direction into the board, replaces pipe if one does not exist
func (b *Board) InsertPipe(x int, y int, pipeType PipeType, pipeDirection PipeDirection) error {
	if b.containsPoint(&Point{x, y}) {
		pipe := newPipe(x, y, pipeType, pipeDirection)
		b.Cells[x][y] = &pipe
		return nil
	}
	return fmt.Errorf("Pipe can not be inserted into, (%d, %d) position does not exist on board", x, y)
}

//IsEmpty checks if a board has no pipes inside of it
func (b *Board) IsEmpty() bool {
	for x := 0; x < len(b.Cells); x++ {
		for y := 0; y < len(b.Cells[x]); y++ {
			if b.Cells[x][y] != nil {
				return false
			}
		}
	}
	return true
}

//CopyBoard returns a copy of the given board
func CopyBoard(b *Board) Board {

	newBoard := NewEmptyBoard(b.NumberOfColumns, b.NumberOfRows)

	for x := 0; x < b.NumberOfColumns; x++ {
		for y := 0; y < b.NumberOfRows; y++ {
			pipe := b.Cells[x][y]
			//TODO: DECIDE IF YOU WANT TO CREATE AN 'EMPTY' PIPE TYPE
			if b.Cells[x][y] != nil {
				newBoard.Cells[x][y] = &Pipe{
					Type:      pipe.Type,
					Direction: pipe.Direction,
					Level:     pipe.Level,
					X:         pipe.X,
					Y:         pipe.Y,
				}
			}
		}
	}

	return newBoard
}

//RotatePipeClockwise Rotates the pipe at the given x and y clockwise if the board contains the given x and y
func (b *Board) RotatePipeClockwise(x int, y int) {
	if b.containsPoint(&Point{x, y}) && b.Cells[x][y] != nil {
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

		//drop pipes that are above a gap
		boardReport.PipeMovementAnimations, boardReport.MaximumAnimationTime = b.dropFloatingPipes()

		//add in new pipes into the empty slots
		pipeMovementAnimations, maxAnimationTime := b.addMissingPipesToBoard()

		boardReport.PipeMovementAnimations = append(boardReport.PipeMovementAnimations, pipeMovementAnimations...)

		if boardReport.MaximumAnimationTime < maxAnimationTime {
			boardReport.MaximumAnimationTime = maxAnimationTime
		}

		copyBoard := CopyBoard(b)
		boardReport.Board = &copyBoard

		boardReports = append(boardReports, boardReport)
	}

	return boardReports
}

func (b *Board) UpdateBoardPipeConnectionsButNoNewPipes() []BoardReport {

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
		//b.addSpecialPipesToBoardUsingClosedTrees(closedTrees)

		//drop pipes that are above a gap
		boardReport.PipeMovementAnimations, boardReport.MaximumAnimationTime = b.dropFloatingPipes()

		//return if the number of connects was larger than zero

		copyBoard := CopyBoard(b)
		boardReport.Board = &copyBoard

		boardReports = append(boardReports, boardReport)
	}

	return boardReports
}

func (b *Board) findAllClosedPipeTrees() []*pipeTree {

	visitedPoints := make(map[Point]bool)

	closedTrees := make([]*pipeTree, 0, 0)

	for x := 0; x < len(b.Cells); x++ {

		for y := 0; y < len(b.Cells[x]); y++ {

			if _, visited := visitedPoints[Point{x, y}]; visited {
				continue
			}

			visitedPoints[Point{x, y}] = true

			//TODO: DECIDE IF YOU WANT AN EMPTY PIPE TYPE
			if b.Cells[x][y] == nil {
				continue
			}

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
					pipeTree.pipe.Level = level0
				case size < level2Size:
					pipeTree.pipe.Level = level1
				case size < level3Size:
					pipeTree.pipe.Level = level2
				case size >= level3Size:
					pipeTree.pipe.Level = level3
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
			pipe := b.Cells[pipeTree.X][pipeTree.Y]
			if pipe != nil {
				if destroyedPipe, ok := b.deleteFromBoard(pipeTree.X, pipeTree.Y); ok {
					destroyedPipes = append(destroyedPipes, destroyedPipe...)
				}
			}
		}
	}

	return destroyedPipes
}

func (b *Board) deleteFromBoard(x int, y int) ([]DestroyedPipe, bool) {

	destroyedPipes := make([]DestroyedPipe, 0, 0)

	if b.containsPoint(&Point{x, y}) {

		if b.Cells[x][y] != nil {
			pipe := b.Cells[x][y]

			destroyedPipes = append(destroyedPipes, DestroyedPipe{Type: pipe.Type, X: pipe.X, Y: pipe.Y})

			b.Cells[x][y] = nil

			if pipe.Type == ENDEXPLOSION2 {
				for _, point := range positionInSquareRange(pipe.X, pipe.Y, 2) {
					if destroyedPipe, ok := b.deleteFromBoard(point.X, point.Y); ok {
						destroyedPipes = append(destroyedPipes, destroyedPipe...)
					}
				}
			} else if pipe.Type == ENDEXPLOSION3 {
				for _, point := range positionInSquareRange(pipe.X, pipe.Y, 3) {
					if destroyedPipe, ok := b.deleteFromBoard(point.X, point.Y); ok {
						destroyedPipes = append(destroyedPipes, destroyedPipe...)
					}
				}
			}

			return destroyedPipes, true
		}
	}

	return []DestroyedPipe{}, false
}

func positionInSquareRange(startX int, startY int, size int) []Point {
	positions := make([]Point, 0, 0)

	if size <= 0 {
		return positions
	}

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if x == 0 && y == 0 {
				continue
			}

			positions = append(positions,
				Point{X: startX + x, Y: startY + y},
				Point{X: startX + x, Y: startY - y},
				Point{X: startX - x, Y: startY + y},
				Point{X: startX - x, Y: startY - y},
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

		switch rootpipeTree.pipe.Level {
		case level2:
			newPipe := newPipe(pipeTree.X, pipeTree.Y, ENDEXPLOSION2, getRandomPipeDirection())
			b.Cells[pipeTree.X][pipeTree.Y] = &newPipe
		case level3:
			newPipe := newPipe(pipeTree.X, pipeTree.Y, ENDEXPLOSION3, getRandomPipeDirection())
			b.Cells[pipeTree.X][pipeTree.Y] = &newPipe
		}
	}
}

func getTravelTime(i int) time.Duration {
	return time.Duration(i) * (time.Millisecond * 100)
}

//Used to move pipes that are floating above a gap towards the ground
func (b *Board) dropFloatingPipes() (pipeMovementAnimations []PipeMovementAnimation, maximumTime time.Duration) {
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
	return
}

func (b *Board) addMissingPipesToBoard() (pipeMovementAnimations []PipeMovementAnimation, maximumTime time.Duration) {

	pipeMovementAnimations = make([]PipeMovementAnimation, 0, 0)

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
}

func (b *Board) containsPoint(p *Point) bool {
	if p.X < 0 || p.X > len(b.Cells)-1 {
		return false
	} else if p.Y < 0 || p.Y > len(b.Cells[p.X])-1 {
		return false
	}
	return true
}

func traversePipeTreeToCheckForClosedConnection(rootPipeTree *pipeTree, visitedPoints map[Point]bool, board *Board) bool {

	isClosedTree := true
	pointsTo := rootPipeTree.pipe.pointsTo()

	for i := 0; i < len(pointsTo); i++ {

		pointToPoint := pointsTo[i]

		if board.containsPoint(&pointToPoint) {

			childTree := newPipeTree(board.Cells[pointToPoint.X][pointToPoint.Y], pointToPoint.X, pointToPoint.Y)

			childPointsToParent := isPipePointingToPipe(childTree.pipe, rootPipeTree.pipe)
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
	pipe     *Pipe

	Point
}

func newPipeTree(pipe *Pipe, x int, y int) pipeTree {
	return pipeTree{
		pipe:  pipe,
		Point: Point{x, y},
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
func (p *Pipe) pointsTo() []Point {

	if p == nil {
		return []Point{}
	}

	currentPoint := Point{p.X, p.Y}

	switch p.Type {
	case END, ENDEXPLOSION2, ENDEXPLOSION3:
		return []Point{pointFromDirection(currentPoint, p.Direction)}
	case LINE:
		switch p.Direction {
		case UP, DOWN:
			return []Point{pointFromDirection(currentPoint, UP), pointFromDirection(currentPoint, DOWN)}
		case LEFT, RIGHT:
			return []Point{pointFromDirection(currentPoint, LEFT), pointFromDirection(currentPoint, RIGHT)}
		}
	case LPIPE:
		switch p.Direction {
		case UP:
			return []Point{pointFromDirection(currentPoint, UP), pointFromDirection(currentPoint, RIGHT)}
		case RIGHT:
			return []Point{pointFromDirection(currentPoint, RIGHT), pointFromDirection(currentPoint, DOWN)}
		case DOWN:
			return []Point{pointFromDirection(currentPoint, DOWN), pointFromDirection(currentPoint, LEFT)}
		case LEFT:
			return []Point{pointFromDirection(currentPoint, LEFT), pointFromDirection(currentPoint, UP)}

		}
	}
	return []Point{}
}

func (p *Pipe) DoesPipePointTo(x int, y int) bool {
	points := p.pointsTo()

	for i := 0; i < len(points); i++ {
		if p.X == points[i].X && p.Y == points[i].Y {
			return true
		}
	}

	return false
}

func pointFromDirection(p Point, d PipeDirection) Point {
	switch d {
	case UP:
		return Point{p.X, p.Y + 1}
	case RIGHT:
		return Point{p.X + 1, p.Y}
	case DOWN:
		return Point{p.X, p.Y - 1}
	case LEFT:
		return Point{p.X - 1, p.Y}
	}
	return p
}
